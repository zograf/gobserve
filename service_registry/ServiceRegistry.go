package serviceregistry

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

type ServiceRegistry struct {
	Port  string
	Infos map[string]*ServiceInfo
}

type CustomContext struct {
	echo.Context
	Sr *ServiceRegistry
}

func (sr *ServiceRegistry) AddServiceInfo(si *ServiceInfo) error {
	_, exists := sr.Infos[si.Name]

	if exists {
		return fmt.Errorf("ServiceInfo with name %s already exists", si.Name)
	}

	sr.Infos[si.Name] = si
	return nil
}

func New() *ServiceRegistry {
	p := os.Getenv("SERVICE_REGISTRY_PORT")
	sr := &ServiceRegistry{Port: p}
	sr.Infos = make(map[string]*ServiceInfo)
	return sr
}

func (sr *ServiceRegistry) Run() {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, sr}
			return next(cc)
		}
	})

	e.GET("/serviceInfo", getAll)
	e.GET("/serviceInfo/:name", getByName)
	e.POST("/serviceInfo", register)

	e.Logger.Fatal(e.Start(sr.Port))
}
