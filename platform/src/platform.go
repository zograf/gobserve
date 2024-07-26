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

	// Copy proxy
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
	proxyFolder := fmt.Sprintf("p%d", config.ServiceCounter)
	cmd := exec.Command("cp", "--recursive", path, newPath)
	cmd.Run()

	dockerfilePath := fmt.Sprintf(".%c%s%cDockerfile", os.PathSeparator, proxyFolder, os.PathSeparator)

	sp := makeService(9001, "service_registry", dockerfilePath, "service_registry", SERVICE_REGISTRY_PORT)

	cf, err := ReadCompose()
	if err != nil {
		return err
	}

	cf.Services[proxyFolder] = sp
	err = saveCompose(cf)
	if err != nil {
		return err
	}

	// Copy microservice
	if err := isDir(path); err != nil {
		return fmt.Errorf("failed to open the directory: %v", err)
	}

	wd, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory")
	}

	config, err = readConfig()
	if err != nil {
		return fmt.Errorf("failed to read the config file: %v", err)
	}

	newPath = fmt.Sprintf("%s%cms%d", wd, os.PathSeparator, config.ServiceCounter)
	serviceFolder := fmt.Sprintf("ms%d", config.ServiceCounter)
	cmd = exec.Command("cp", "--recursive", path, newPath)
	cmd.Run()

	dockerfilePath = fmt.Sprintf(".%c%s%cDockerfile", os.PathSeparator, serviceFolder, os.PathSeparator)

	s := makeService(1001, proxyFolder, dockerfilePath, proxyFolder, 9001)

	cf, err = ReadCompose()
	if err != nil {
		return err
	}
	cf.Services[serviceFolder] = s
	err = saveCompose(cf)
	if err != nil {
		return err
	}

	config.ServiceCounter++
	err = writeConfig(config)
	return err
}
