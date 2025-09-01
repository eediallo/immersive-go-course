package cmd

import (
	"fmt"
	"log"
	"os"
)

func Execute() {
	/*
		1. read directory
		2. display content to sdtout
	*/

	files, err := os.ReadDir(".")

	if err != nil {
		log.Fatalf("Error reading directory %s", err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
