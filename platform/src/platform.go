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

	err := clean()
	if err != nil {
		return err
	}

	sr, err := makeServiceRegistry()
	if err != nil {
		return err
	}

	cf.Services["service_registry"] = *sr

	err = saveCompose(cf)
	if err != nil {
		return err
	}

	_, err = makeServiceByName("gateway", 0, GATEWAY_PORT)
	if err != nil {
		return err
	}

	_, err = makeServiceByName("aggregator", 0, AGGREGATOR_PORT)
	if err != nil {
		return err
	}

	err = writeConfig(&PlatformConfig{
		ServiceCounter:  1,
		NextProxyPort:   10001,
		NextServicePort: 40001,
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

	proxyName, err := makeServiceByName("proxy", config.ServiceCounter, config.NextProxyPort)
	if err != nil {
		return err
	}

	err = makeMicroservice(path, proxyName, config.ServiceCounter, config.NextServicePort, config.NextProxyPort)
	if err != nil {
		return err
	}

	config.ServiceCounter++
	config.NextProxyPort++
	config.NextServicePort++
	err = writeConfig(config)
	return err
}

func Clean() error {
	return clean()
}
