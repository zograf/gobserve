package gateway

import "github.com/labstack/echo/v4"

type Gateway struct {
	Component *Component
}
type Component struct {
	Info   *ServiceInfo
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
