package platform

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Init() {
	cf := &ComposeFile{
		Version:  "3",
		Services: make(map[string]Service),
	}

	serviceregistry := Service{
		Build: BuildConf{
			Context:    CONTEXT,
			Dockerfile: SERVICE_REGISTRY_DOCKERFILE,
		},
		Environment: map[string]string{
			"SERVICE_REGISTRY_PORT": fmt.Sprintf(":%d", SERVICE_REGISTRY_PORT),
			"SERVICE_REGISTRY_IP":   IP,
		},
		Ports: []string{
			fmt.Sprintf("%d:%d", SERVICE_REGISTRY_PORT, SERVICE_REGISTRY_PORT),
		},
	}

	cf.Services["service_registry"] = serviceregistry

	SaveCompose(cf)
}

func Run() {

}

func Add() {

}

func ReadCompose() (*ComposeFile, error) {
	data, err := os.ReadFile(COMPOSE_FILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to open the compose file")
	}

	var cf ComposeFile
	err = yaml.Unmarshal(data, &cf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data")
	}

	return &cf, nil
}

func SaveCompose(cf *ComposeFile) error {
	data, err := yaml.Marshal(cf)
	if err != nil {
		return fmt.Errorf("failed to marshal the compose file")
	}

	err = os.WriteFile(COMPOSE_FILE_NAME, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save the compose file")
	}

	return nil
}
