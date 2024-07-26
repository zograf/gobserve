package microservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func getTest(c echo.Context) error {
	//cc := c.(*CustomContext)

	return c.JSON(http.StatusOK, "Microservice A responsse")
}

func getServiceInfo(c echo.Context) error {
	cc := c.(*CustomContext)
	ms := cc.Ms
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://%s%s/serviceInfo", ms.ServiceRegistryIp, ms.ServiceRegistryPort)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to create HTTP request: %w", err).Error(),
		})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to send HTTP request: %w", err).Error(),
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Errorf("HTTP request failed: %w", err).Error(),
			})
		}

		bodyString := string(bodyBytes)
		return c.JSON(http.StatusOK, bodyString)
	}

	return c.JSON(resp.StatusCode, "Microservice A failed")
}

func register(c echo.Context) error {
	cc := c.(*CustomContext)
	ms := cc.Ms

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := ms.Info
	jsonPayload, err := json.Marshal(info)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to marshal JSON payload: %w", err).Error(),
		})
	}

	url := fmt.Sprintf("http://%s%s/serviceInfo", ms.ServiceRegistryIp, ms.ServiceRegistryPort)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to create HTTP request: %w", err).Error(),
		})
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to send HTTP request: %w", err).Error(),
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Errorf("HTTP request failed: %w", err).Error(),
			})
		}

		bodyString := string(bodyBytes)
		return c.JSON(http.StatusOK, bodyString)
	}

	return c.JSON(http.StatusOK, "Register successful")
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
