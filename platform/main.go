package main

import (
	"flag"
	"log"
	"os"

	platform "github.com/zograf/gobserve/platform/src"
)

func main() {
	isInit := flag.Bool("init", false, "Initializes a docker compose for deployment. This will delete all existing docker files. Cleans the directory before init.")
	isCompose := flag.Bool("run", false, "Composes the docker files provided.")
	isClean := flag.Bool("clean", false, "Cleans the platform directory.")
	microservicePath := flag.String("add", "", "Provide a path to the docker file. The docker file is then prepared for deployment.")

	flag.Parse()

	if !*isInit && !*isCompose && !*isClean && *microservicePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *isClean {
		err := platform.Clean()
		if err != nil {
			log.Fatal(err)
		}
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
