package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"
)

func Init() error {
	cf := &ComposeFile{
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

	err := clean()
	if err != nil {
		return err
	}

	err = saveCompose(cf)
	if err != nil {
		return err
	}

	err = writeConfig(&PlatformConfig{
		ServiceCounter: 1,
	})

	return err
}

func Run() {
	cmd := exec.Command("docker-compose", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("failed to docker compose: %v", err)
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Fatalf("docker compose failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	err = cmd.Process.Signal(syscall.SIGSTOP)
	if err != nil {
		log.Fatalf("failed to SIGSTOP: %v", err)
	}

	downCmd := exec.Command("docker-compose", "down")
	downCmd.Stdout = os.Stdout
	downCmd.Stderr = os.Stderr
	err = downCmd.Run()
	if err != nil {
		log.Fatalf("failed to stop docker compose: %v", err)
	}
}

func Add(path string) error {
	if err := isDir(path); err != nil {
		return fmt.Errorf("failed to open the directory: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory")
	}

	config, err := readConfig()
	if err != nil {
		return fmt.Errorf("failed to read the config file: %v", err)
	}

	newPath := fmt.Sprintf("%s%cp%d", wd, os.PathSeparator, config.ServiceCounter)
	cmd := exec.Command("cp", "--recursive", path, newPath)
	cmd.Run()

	config.ServiceCounter++
	err = writeConfig(config)
	return err
}

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
		i--
	}

	return nil
}
