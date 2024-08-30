package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

func getLogs(c echo.Context) error {
	cc := c.(*CustomContext)

	infos, err := cc.Sr.GetInfos()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	ret := LogData{
		ServiceCount: 0,
		Data:         make(map[string][]LogEntry),
	}

	for _, info := range infos {
		if info.Name == cc.Sr.Component.Info.Name {
			continue
		}
		url := fmt.Sprintf("http://%s%s/log", info.Ip, info.Port)
		logData, err := getLogEntries(url)
		if err != nil {
			continue
		}
		ret.ServiceCount += 1
		ret.Data[info.Name] = logData
	}

	return c.JSON(http.StatusOK, ret)
}

func getLogEntries(url string) ([]LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create log request")
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do log request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var ret []LogEntry
	err = json.Unmarshal(body, &ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}
