package main

import microservice "github.com/zograf/gobserve/microservice_a/src"

func main() {
	ms := microservice.New()
	ms.Run()
}
