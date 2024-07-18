package main

import (
	"flag"
	"os"

	platform "github.com/zograf/gobserve/platform/src"
)

func main() {
	isInit := flag.Bool("init", false, "Initializes a docker compose for deployment. This will delete all existing docker files.")
	isCompose := flag.Bool("run", false, "Composes the docker files provided.")
	dockerFile := flag.String("add", "", "Provide a path to the docker file. The docker file is then prepared for deployment.")

	flag.Parse()

	if !*isInit && !*isCompose && *dockerFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *isInit {
		platform.Init()
	}
	if *dockerFile != "" {
		// Add the docker file path to the compose
		os.Exit(1)
	}
	if *isCompose {
		platform.Run()
	}
}
