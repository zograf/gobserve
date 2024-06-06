package main

import (
	serviceregistry "github.com/zograf/gobserve/service_registry/pkg"
)

func main() {
	sr := serviceregistry.New()
	sr.Run()
}
