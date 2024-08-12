package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Proxy struct {
	Port           string
	Ip             string
	Infos          map[string]*ServiceInfo
	ProxiedService *ServiceInfo
	SRIP           string
	SRPort         string
}

func (proxy *Proxy) getRealInfos() (map[string]*ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	url := fmt.Sprintf("http://%s%s%s", proxy.SRIP, proxy.SRPort, "/serviceInfo")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req = req.WithContext(ctx)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var ret map[string]*ServiceInfo
	err = json.Unmarshal(body, &ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (proxy *Proxy) GetInfos() (map[string]*ServiceInfo, error) {
	if proxy.ProxiedService == nil {
		return nil, fmt.Errorf("microservice not yet registered")
	}
	realInfos, err := proxy.getRealInfos()

	if err != nil {
		return nil, fmt.Errorf("failed to get service infos.\nerror message: %s", err.Error())
	}

	if _, exists := realInfos[proxy.ProxiedService.Name]; exists {
		delete(realInfos, proxy.ProxiedService.Name)
	}

	for _, val := range realInfos {
		val.Ip = proxy.Ip
		val.Port = proxy.Port
	}

	proxy.Infos = realInfos
	return realInfos, nil
}

func (proxy *Proxy) addToRegistry(si *ServiceInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	jsonData, err := json.Marshal(si)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s%s%s", proxy.SRIP, proxy.SRPort, "/serviceInfo")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(jsonData)))

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	res.Body.Close()
	return nil
}

func (proxy *Proxy) AddServiceInfo(si *ServiceInfo) error {
	if proxy.ProxiedService != nil {
		return fmt.Errorf("ServiceInfo with name %s already exists", si.Name)
	}

	proxy.ProxiedService = si

	proxiedInfo := &ServiceInfo{
		Name: si.Name,
		Port: proxy.Port,
		Ip:   proxy.Ip,
	}

	err := proxy.addToRegistry(proxiedInfo)
	fmt.Println(proxiedInfo)
	return err
}

func New() *Proxy {
	p := os.Getenv("PORT")
	ip := os.Getenv("IP")
	srIp := os.Getenv("SERVICE_REGISTRY_IP")
	srPort := os.Getenv("SERVICE_REGISTRY_PORT")
	sr := &Proxy{
		Ip:     ip,
		Port:   p,
		SRIP:   srIp,
		SRPort: srPort,
	}
	sr.Infos = make(map[string]*ServiceInfo)
	return sr
}

func (proxy *Proxy) Run() {
	e := echo.New()

	// Middleware
	e.Use(logRequest)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				Sr:      proxy,
			}
			return next(cc)
		}
	})

	// Routes
	e.GET("/health", healthCheck)
	e.GET("/serviceInfo", getAll)
	e.GET("/serviceInfo/:name", getByName)
	e.POST("/serviceInfo", register)
	e.Any("/*", proxyPass)

	url := fmt.Sprintf("%s%s", proxy.Ip, proxy.Port)
	e.Logger.Fatal(e.Start(url))
}

func formatLogEntry(clientIP, method, url, protocol string, requestHeaders map[string]string, requestBody string, responseTimestamp string, statusCode int, responseHeaders map[string]string, responseBody string, duration int64) string {
	return fmt.Sprintf(
		"request_timestamp=%s|client_ip=%s|method=%s|url=%s|protocol=%s|request_headers=%s|request_body=%s|response_timestamp=%s|status_code=%d|response_headers=%s|response_body=%s|duration_ms=%d",
		time.Now().UTC().Format(time.RFC3339),
		clientIP,
		method,
		url,
		protocol,

		formatHeaders(requestHeaders),
		requestBody,

		responseTimestamp,
		statusCode,
		formatHeaders(responseHeaders),
		strings.ReplaceAll(responseBody, "\n", " "),
		duration,
	)
}

func formatHeaders(headers map[string]string) string {
	var headerParts []string
	for key, value := range headers {
		headerParts = append(headerParts, fmt.Sprintf("%s:%s", key, value))
	}
	return strings.Join(headerParts, ",")
}

type CustomResponseWriter struct {
	echo.Response
	Body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.Response.Writer.Write(b)
}

func logRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		clientIP := c.RealIP()

		var requestBody bytes.Buffer
		_, err := io.Copy(&requestBody, c.Request().Body)
		if err != nil {
			return fmt.Errorf("error reading request body: %w", err)
		}
		c.Request().Body = io.NopCloser(&requestBody)

		requestHeaders := make(map[string]string)
		for name, values := range c.Request().Header {
			requestHeaders[name] = values[0]
		}

		res := c.Response()
		rec := &CustomResponseWriter{
			Response: *res,
			Body:     new(bytes.Buffer),
		}
		c.Response().Writer = rec

		if err := next(c); err != nil {
			c.Error(err)
		}

		responseBody := rec.Body.String()

		responseHeaders := make(map[string]string)
		for name, values := range c.Response().Header() {
			responseHeaders[name] = values[0]
		}

		duration := time.Since(start).Milliseconds()

		logEntry := formatLogEntry(
			clientIP,
			c.Request().Method,
			c.Request().URL.String(),
			c.Request().Proto,
			requestHeaders,
			requestBody.String(),
			time.Now().UTC().Format(time.RFC3339),
			c.Response().Status,
			responseHeaders,
			responseBody,
			duration,
		)

		file, err := os.OpenFile("proxy.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error opening log file: %w", err)
		}
		defer file.Close()

		_, err = file.WriteString(logEntry + "\n")
		if err != nil {
			return fmt.Errorf("error writing to log file: %w", err)
		}

		return nil
	}
}
