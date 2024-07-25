package main

import (
	proxy "github.com/zograf/gobserve/proxy/src"
)

func main() {
	sr := proxy.New()
	sr.Run()
}
