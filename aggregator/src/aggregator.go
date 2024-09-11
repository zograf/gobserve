package aggregator

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
	"github.com/labstack/echo/v4/middleware"
)

func New() *Aggregator {
	p := os.Getenv("PORT")
	ip := os.Getenv("IP")
	srIp := os.Getenv("SERVICE_REGISTRY_IP")
	srPort := os.Getenv("SERVICE_REGISTRY_PORT")
	name := os.Getenv("NAME")
	gateway := &Aggregator{
		Component: &Component{
			Info: &ServiceInfo{
				Ip:   ip,
				Port: p,
				Name: name,
			},
			SRIP:   srIp,
			SRPort: srPort,
		},
	}
	return gateway
}

func (agg *Aggregator) Run() {
	e := echo.New()

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				Sr:      *agg,
			}
			return next(cc)
		}
	})

	e.GET("/health", healthCheck)
	e.GET("/log", getLogs)

	err := register(*agg)

	if err != nil {
		return
	}

	url := fmt.Sprintf("%s%s", agg.Component.Info.Ip, agg.Component.Info.Port)
	e.Logger.Fatal(e.Start(url))
}

func (gateway *Aggregator) GetInfos() (map[string]*ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	url := fmt.Sprintf("http://%s%s%s", gateway.Component.SRIP, gateway.Component.SRPort, "/serviceInfo")
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

func register(agg Aggregator) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := agg.Component.Info

	jsonPayload, err := json.Marshal(*info)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	url := fmt.Sprintf("http://%s%s/serviceInfo", agg.Component.SRIP, agg.Component.SRPort)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("HTTP request failed: %w", err)
		}
		fmt.Println("[*] Register successful!")
	}

	return nil
}
