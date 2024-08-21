package proxy

import "github.com/labstack/echo/v4"

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type CustomContext struct {
	echo.Context
	Sr Registry
}

type Registry interface {
	GetInfos() (map[string]*ServiceInfo, error)
	AddServiceInfo(si *ServiceInfo) error
	GetProxiedService() *ServiceInfo
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

const LOG_FILE string = "proxy.log"
