package platform

type BuildConf struct {
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

type Service struct {
	Build       BuildConf         `yaml:"build,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
}

type ComposeFile struct {
	Services map[string]Service `yaml:"services"`
}

const COMPOSE_FILE_NAME string = "docker-compose.yml"
const CONTEXT string = "../"
const SERVICE_REGISTRY_DOCKERFILE string = "./service_registry/Dockerfile"
const SERVICE_REGISTRY_PORT int = 7777
const IP string = "localhost"
