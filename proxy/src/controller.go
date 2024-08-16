package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func getAll(c echo.Context) error {
	cc := c.(*CustomContext)

	infos, err := cc.Sr.GetInfos()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, infos)
}

func getByName(c echo.Context) error {
	name := c.Param("name")
	cc := c.(*CustomContext)

	infos, err := cc.Sr.GetInfos()
	if err != nil {
		return err
	}

	value, exists := infos[name]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service info with that name was not found",
		})
	}

	return c.JSON(http.StatusOK, value)
}

func proxyPass(c echo.Context) error {
	cc := c.(*CustomContext)

	service := cc.Sr.GetProxiedService()
	if service == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service not running",
		})
	}

	path := c.Request().URL.Path
	splitPath := strings.SplitN(path[1:], "/", 2)

	var name, endpoint string
	if len(splitPath) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Bad path formatting",
		})
	} else if len(splitPath) == 1 {
		name = splitPath[0]
		endpoint = ""
	} else {
		name = splitPath[0]
		endpoint = splitPath[1]
	}

	if name == service.Name {
		url := fmt.Sprintf("http://%s%s/%s", service.Ip, service.Port, endpoint)
		err := forwardRequest(c, url)
		if err != nil {
			return err
		}
	} else {
		infos, err := cc.Sr.GetInfos()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		val, found := infos[name]
		if !found {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Bad path",
			})
		}
		url := fmt.Sprintf("http://%s%s/%s", val.Ip, val.Port, endpoint)
		err = forwardRequest(c, url)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusNoContent, nil)
}

func forwardRequest(c echo.Context, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	req := c.Request()

	forwardedReq, err := http.NewRequest(req.Method, url, req.Body)
	forwardedReq.WithContext(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create forwarded request")
	}

	for key, values := range req.Header {
		for _, value := range values {
			forwardedReq.Header.Add(key, value)
		}
	}

	client := http.Client{}
	fmt.Printf("[*] Forwarding to: %s\n", url)
	resp, err := client.Do(forwardedReq)
	if err != nil {
		return fmt.Errorf("Failed to forward request")
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}
	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to copy response")
	}

	return nil
}

func register(c echo.Context) error {
	si := &ServiceInfo{}

	if err := c.Bind(si); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	cc := c.(*CustomContext)
	if err := cc.Sr.AddServiceInfo(si); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, si)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
