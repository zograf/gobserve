package microservice

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
)

func New() *Microservice {
	ip := os.Getenv("IP")
	p := os.Getenv("PORT")
	srIp := os.Getenv("SERVICE_REGISTRY_IP")
	srPort := os.Getenv("SERVICE_REGISTRY_PORT")

	ms := &Microservice{
		Ip:                  ip,
		Port:                p,
		ServiceRegistryIp:   srIp,
		ServiceRegistryPort: srPort,
		Info: ServiceInfo{
			Ip:   ip,
			Port: p,
			Name: "A",
		},
	}
	return ms
}

func (ms *Microservice) Run() {
	e := echo.New()

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				Ms:      *ms,
			}
			return next(cc)
		}
	})

	// Routes
	e.GET("/health", healthCheck)
	e.GET("/findServiceInfo", getServiceInfo)
	e.POST("/register", register)
	e.GET("/", getTest)

	url := fmt.Sprintf("%s%s", ms.Ip, ms.Port)
	e.Logger.Fatal(e.Start(url))
}
