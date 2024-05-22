package core

import "github.com/labstack/echo/v4"

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type CustomContext struct {
	echo.Context
	Sr Registry
}

type Registry interface {
	GetInfos() map[string]*ServiceInfo
	AddServiceInfo(si *ServiceInfo) error
}
