package microservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

func register(ms Microservice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := ms.Info
	jsonPayload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	url := fmt.Sprintf("http://%s%s/serviceInfo", ms.ServiceRegistryIp, ms.ServiceRegistryPort)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("HTTP request failed: %w", err)
		}
		fmt.Println("[*] Register successful!")
	}

	return nil
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
