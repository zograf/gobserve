package microservice

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
	Ms Microservice
}

type Microservice struct {
	Port                string
	Ip                  string
	ServiceRegistryIp   string
	ServiceRegistryPort string
	Info                ServiceInfo
}

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
