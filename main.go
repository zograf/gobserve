package main

import "github.com/zograf/gobserve/simple_proxy"

func main() {
	go simple_proxy.Make_proxy()
	simple_proxy.Make_http_server()
}
