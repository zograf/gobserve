package platform

type BuildConf struct {
	Context    string   `yaml:"context,omitempty"`
	Dockerfile string   `yaml:"dockerfile,omitempty"`
	Args       []string `yaml:"args,omitempty"`
}

type Service struct {
	Build       BuildConf                     `yaml:"build,omitempty"`
	Environment map[string]string             `yaml:"environment,omitempty"`
	HealthCheck HealthCheck                   `yaml:"healthcheck,omitempty"`
	DependsOn   map[string]DependsOnCondition `yaml:"depends_on,omitempty"`
	Ports       []string                      `yaml:"ports,omitempty"`
	NetworkMode string                        `yaml:"network_mode,omitempty"`
}

type ComposeFile struct {
	Services map[string]Service `yaml:"services"`
}

type PlatformConfig struct {
	ServiceCounter  int `json:"service_counter"`
	NextProxyPort   int `json:"next_proxy_port"`
	NextServicePort int `json:"next_service_port"`
}

type HealthCheck struct {
	Test     []string `yaml:"test,flow"`
	Interval string   `yaml:"interval"`
	Timeout  string   `yaml:"timeout"`
	Retries  int      `yaml:"retries"`
}

type DependsOnCondition struct {
	Condition string `yaml:"condition"`
}

// TODO: Get these in a config?
const COMPOSE_FILE_NAME string = "docker-compose.yml"
const CONTEXT string = "./"
const SERVICE_REGISTRY_PORT int = 7777
const IP string = "127.0.0.1"
const CONFIG_PATH string = "./config.json"
const PROXY_PATH string = "../proxy"
const SERVICE_REGISTRY_PATH string = "../service_registry"
const GATEWAY_PATH string = "../gateway"
const GATEWAY_PORT int = 8080
