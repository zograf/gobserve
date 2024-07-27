package serviceregistry

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

type ServiceRegistry struct {
	Port string
	Ip   string
	// TODO: Maybe don't use in memory storage?
	Infos map[string]*ServiceInfo
}

func (sr *ServiceRegistry) GetInfos() (map[string]*ServiceInfo, error) {
	return sr.Infos, nil
}

func (sr *ServiceRegistry) AddServiceInfo(si *ServiceInfo) error {
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
	p := os.Getenv("PORT")
	ip := os.Getenv("IP")
	sr := &ServiceRegistry{
		Port: p,
		Ip:   ip,
	}
	sr.Infos = make(map[string]*ServiceInfo)
	return sr
}

func (sr *ServiceRegistry) Run() {
	e := echo.New()

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				Sr:      sr,
			}
			return next(cc)
		}
	})

	// Routes
	e.GET("/health", healthCheck)
	e.GET("/serviceInfo", getAll)
	e.GET("/serviceInfo/:name", getByName)
	e.POST("/serviceInfo", register)

	url := fmt.Sprintf("%s%s", sr.Ip, sr.Port)
	e.Logger.Fatal(e.Start(url))
}
