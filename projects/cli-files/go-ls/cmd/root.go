package cmd

import (
	"fmt"
	"log"
	"os"
)

func Execute() {
	/*
		0. pass dir as an argument
		1. read directory
		2. display content to sdtout
	*/

	arg := os.Args[1]

	if len(os.Args) > 2 {
		log.Fatal("Only one argument is allowed")
	}

	files, err := os.ReadDir(arg)

	if err != nil {
		log.Fatalf("Error reading directory %s", err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
