package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zograf/gobserve/core"
)

type Proxy struct {
	Port           string
	Ip             string
	Infos          map[string]*core.ServiceInfo
	ProxiedService *core.ServiceInfo
	SRIP           string
	SRPort         string
}

func (proxy *Proxy) getRealInfos() (map[string]*core.ServiceInfo, error) {
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

	var ret map[string]*core.ServiceInfo
	err = json.Unmarshal(body, &ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (proxy *Proxy) GetInfos() (map[string]*core.ServiceInfo, error) {
	realInfos, err := proxy.getRealInfos()

	if err != nil {
		return nil, fmt.Errorf("failed to get service infos.\nerror message: %s", err.Error())
	}

	for _, val := range realInfos {
		val.Ip = proxy.Ip
		val.Port = proxy.Port
	}

	proxy.Infos = realInfos
	return realInfos, nil
}

func (proxy *Proxy) addToRegistry(si *core.ServiceInfo) error {
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

func (proxy *Proxy) AddServiceInfo(si *core.ServiceInfo) error {
	if proxy.ProxiedService != nil {
		return fmt.Errorf("ServiceInfo with name %s already exists", si.Name)
	}

	proxy.ProxiedService = si

	proxiedInfo := &core.ServiceInfo{
		Name: si.Name,
		Port: proxy.Port,
		Ip:   proxy.Ip,
	}

	err := proxy.addToRegistry(proxiedInfo)
	fmt.Println(proxiedInfo)
	return err
}

func New() *Proxy {
	p := os.Getenv("PROXY_PORT")
	ip := os.Getenv("PROXY_IP")
	srIp := os.Getenv("SERVICE_REGISTRY_IP")
	srPort := os.Getenv("SERVICE_REGISTRY_PORT")
	sr := &Proxy{
		Ip:     ip,
		Port:   p,
		SRIP:   srIp,
		SRPort: srPort,
	}
	sr.Infos = make(map[string]*core.ServiceInfo)
	return sr
}

func (proxy *Proxy) Run() {
	e := echo.New()

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &core.CustomContext{
				Context: c,
				Sr:      proxy,
			}
			return next(cc)
		}
	})

	// Routes
	e.GET("/serviceInfo", getAll)
	e.GET("/serviceInfo/:name", getByName)
	e.POST("/serviceInfo", register)
	e.Any("/*", proxyPass)

	e.Logger.Fatal(e.Start(proxy.Port))
}
