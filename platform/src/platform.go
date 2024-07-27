package platform

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Init() error {
	cf := &ComposeFile{
		Services: make(map[string]Service),
	}

	portStr := fmt.Sprintf(":%d", SERVICE_REGISTRY_PORT)

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
		HealthCheck: makeHealthCheck(portStr),
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
	config, err := readConfig()
	if err != nil {
		return fmt.Errorf("failed to read the config file: %v", err)
	}

	proxyName, err := makeProxy(PROXY_PATH, config.ServiceCounter, 9001)
	if err != nil {
		return err
	}

	err = makeMicroservice(path, proxyName, config.ServiceCounter, 1001)
	if err != nil {
		return err
	}

	config.ServiceCounter++
	err = writeConfig(config)
	return err
}
