package main

import aggregator "github.com/zograf/gobserve/gateway/src"

func main() {
	agg := aggregator.New()
	agg.Run()
}
