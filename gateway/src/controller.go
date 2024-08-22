package gateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func forwardRequest(c echo.Context) error {
	cc := c.(*CustomContext)

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

	infos, err := cc.Sr.GetInfos()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var url string
	val, found := infos[name]
	if found {
		url = fmt.Sprintf("http://%s%s/%s", val.Ip, val.Port, endpoint)
		err = forward(c, url)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
	} else {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service not found",
		})
	}

	return c.NoContent(204)
}

func forward(c echo.Context, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	req := c.Request()

	forwardedReq, err := http.NewRequest(req.Method, url, req.Body)
	forwardedReq = forwardedReq.WithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to create forwarded request")
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
		return fmt.Errorf("failed to forward request")
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
		return fmt.Errorf("failed to copy response")
	}

	return nil
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
