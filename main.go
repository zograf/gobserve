package main

import (
	"flag"
	"os"
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
		// Initialize an "empty" docker compose
		os.Exit(1)
	}
	if *dockerFile != "" {
		// Add the docker file path to the compose
		os.Exit(1)
	}
	if *isCompose {
		// Run the compose
		os.Exit(1)
	}
}
