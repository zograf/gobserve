package main

import (
	gateway "github.com/zograf/gobserve/gateway/src"
)

func main() {
	gw := gateway.New()
	gw.Run()
}
