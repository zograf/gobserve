package gateway

import "github.com/labstack/echo/v4"

type Gateway struct {
	Port   string
	Ip     string
	SRIP   string
	SRPort string
}

type CustomContext struct {
	echo.Context
	Sr Gateway
}

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
