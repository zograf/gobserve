package main

import (
	serviceregistry "github.com/zograf/gobserve/service_registry/src"
)

func main() {
	sr := serviceregistry.New()
	sr.Run()
}
