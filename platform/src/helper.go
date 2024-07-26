package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

func writeConfig(config *PlatformConfig) error {
	jsonString, err := json.Marshal(config)
	if err != nil {
		return err
	}

	bytes := []byte(jsonString)

	err = os.WriteFile(CONFIG_PATH, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readConfig() (*PlatformConfig, error) {
	file, err := os.Open(CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to open the config file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %v", err)
	}

	var config PlatformConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return &config, nil
}

func isDir(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("couldn't find the file")
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("not a directory")
	}
	return nil
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

func saveCompose(cf *ComposeFile) error {
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

func clean() error {
	config, err := readConfig()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	i := config.ServiceCounter - 1
	for i > 0 {
		fileName := fmt.Sprintf("%s%cp%d", wd, os.PathSeparator, i)
		cmd := exec.Command("rm", "-r", fileName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		fileName = fmt.Sprintf("%s%cms%d", wd, os.PathSeparator, i)
		cmd = exec.Command("rm", "-r", fileName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		i--
	}

	return nil
}

func makeHealthCheck(port string) HealthCheck {
	url := fmt.Sprintf("http://localhost%s/health", port)
	hc := HealthCheck{
		Test:     []string{"CMD", "curl", "-f", url},
		Timeout:  "5s",
		Retries:  5,
		Interval: "5s",
	}
	return hc
}

func makeService(port int, dependsOn, dockerfilePath, srIp string, srPort int) Service {
	// TODO: Figure out a logic behind proxy port
	portStr := fmt.Sprintf(":%d", port)
	s := Service{
		Build: BuildConf{
			Context:    "./",
			Dockerfile: dockerfilePath,
		},
		Environment: map[string]string{
			"PORT":                  portStr,
			"IP":                    IP,
			"SERVICE_REGISTRY_PORT": fmt.Sprintf(":%d", srPort),
			"SERVICE_REGISTRY_IP":   srIp,
		},
		Ports: []string{
			fmt.Sprintf("%d:%d", port, port),
		},
		HealthCheck: makeHealthCheck(portStr),
		DependsOn: map[string]DependsOnCondition{
			dependsOn: {Condition: "service_healthy"},
		},
	}
	return s
}
