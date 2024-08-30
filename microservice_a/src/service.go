package microservice

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func New() *Microservice {
	ip := os.Getenv("IP")
	p := os.Getenv("PORT")
	srIp := os.Getenv("SERVICE_REGISTRY_IP")
	srPort := os.Getenv("SERVICE_REGISTRY_PORT")
	name := os.Getenv("NAME")

	ms := &Microservice{
		Component: &Component{
			Info: &ServiceInfo{
				Ip:   ip,
				Port: p,
				Name: name,
			},
			SRIP:   srIp,
			SRPort: srPort,
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
	//e.POST("/register", register)
	e.GET("/", getTest)

	err := register(*ms)
	if err != nil {
		fmt.Println("[*] Failed to register service: %w", err)
	}

	url := fmt.Sprintf("%s%s", ms.Component.Info.Ip, ms.Component.Info.Port)
	e.Logger.Fatal(e.Start(url))
}
