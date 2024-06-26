package serviceregistry

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/zograf/gobserve/core"
)

type ServiceRegistry struct {
	Port string
	// TODO: Maybe don't use in memory storage?
	Infos map[string]*core.ServiceInfo
}

func (sr *ServiceRegistry) GetInfos() (map[string]*core.ServiceInfo, error) {
	return sr.Infos, nil
}

func (sr *ServiceRegistry) AddServiceInfo(si *core.ServiceInfo) error {
	infos, err := sr.GetInfos()
	if err != nil {
		return err
	}

	_, exists := infos[si.Name]

	if exists {
		return fmt.Errorf("ServiceInfo with name %s already exists", si.Name)
	}

	infos[si.Name] = si
	return nil
}

func New() *ServiceRegistry {
	p := os.Getenv("SERVICE_REGISTRY_PORT")
	sr := &ServiceRegistry{Port: p}
	sr.Infos = make(map[string]*core.ServiceInfo)
	return sr
}

func (sr *ServiceRegistry) Run() {
	e := echo.New()

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &core.CustomContext{
				Context: c,
				Sr:      sr,
			}
			return next(cc)
		}
	})

	// Routes
	e.GET("/serviceInfo", getAll)
	e.GET("/serviceInfo/:name", getByName)
	e.POST("/serviceInfo", register)

	e.Logger.Fatal(e.Start(sr.Port))
}
