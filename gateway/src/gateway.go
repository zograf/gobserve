package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func New() *Gateway {
	p := os.Getenv("PORT")
	ip := os.Getenv("IP")
	srIp := os.Getenv("SERVICE_REGISTRY_IP")
	srPort := os.Getenv("SERVICE_REGISTRY_PORT")
	gateway := &Gateway{
		Ip:     ip,
		Port:   p,
		SRIP:   srIp,
		SRPort: srPort,
	}
	return gateway
}

func (gateway *Gateway) Run() {
	e := echo.New()

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				Sr:      *gateway,
			}
			return next(cc)
		}
	})

	// Routes
	e.GET("/health", healthCheck)
	e.Any("/*", forwardRequest)

	url := fmt.Sprintf("%s%s", gateway.Ip, gateway.Port)
	e.Logger.Fatal(e.Start(url))
}

func (gateway *Gateway) GetInfos() (map[string]*ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	url := fmt.Sprintf("http://%s%s%s", gateway.SRIP, gateway.SRPort, "/serviceInfo")
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

	fmt.Print("[*] SR lookup: ")
	for k := range ret {
		fmt.Printf("%s, ", k)
	}
	fmt.Println()

	return ret, nil
}
