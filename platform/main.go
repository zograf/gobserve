package main

import (
	"flag"
	"log"
	"os"

	platform "github.com/zograf/gobserve/platform/src"
)

func main() {
	isInit := flag.Bool("init", false, "Initializes a docker compose for deployment. This will delete all existing docker files.")
	isCompose := flag.Bool("run", false, "Composes the docker files provided.")
	microservicePath := flag.String("add", "", "Provide a path to the docker file. The docker file is then prepared for deployment.")

	flag.Parse()

	if !*isInit && !*isCompose && *microservicePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *isInit {
		err := platform.Init()
		if err != nil {
			log.Fatal(err)
		}
	}
	if *microservicePath != "" {
		err := platform.Add(*microservicePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *isCompose {
		platform.Run()
	}
}
