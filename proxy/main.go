package main

import (
	proxy "github.com/zograf/gobserve/proxy/pkg"
)

func main() {
	sr := proxy.New()
	sr.Run()
}
