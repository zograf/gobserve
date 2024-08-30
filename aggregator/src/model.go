package aggregator

import "github.com/labstack/echo/v4"

type Aggregator struct {
	Component *Component
}

type Component struct {
	Info   *ServiceInfo
	SRIP   string
	SRPort string
}

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type CustomContext struct {
	echo.Context
	Sr Aggregator
}

type LogData struct {
	ServiceCount int
	Data         map[string][]LogEntry
}

type LogEntry struct {
	RequestTimestamp  string            `json:"request_timestamp"`
	ClientIP          string            `json:"client_ip"`
	Method            string            `json:"method"`
	URL               string            `json:"url"`
	Protocol          string            `json:"protocol"`
	RequestHeaders    map[string]string `json:"request_headers"`
	RequestBody       string            `json:"request_body"`
	ResponseTimestamp string            `json:"response_timestamp"`
	StatusCode        int               `json:"status_code"`
	ResponseHeaders   map[string]string `json:"response_headers"`
	ResponseBody      string            `json:"response_body"`
	DurationMs        int64             `json:"duration_ms"`
}
