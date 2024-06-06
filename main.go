package main

import (
	"github.com/zograf/gobserve/proxy"
	serviceregistry "github.com/zograf/gobserve/service_registry"
)

func main() {
	sr := serviceregistry.New()
	go sr.Run()

	proxy := proxy.New()
	proxy.Run()
}
