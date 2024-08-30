package microservice

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
	Ms Microservice
}

type Microservice struct {
	Component *Component
}

type Component struct {
	Info   *ServiceInfo
	SRIP   string
	SRPort string
}

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
