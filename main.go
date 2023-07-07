package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	// Start exploring from the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		return
	}

	// Explore the directory recursively
	err = exploreDirectory(currentDir, "")
	if err != nil {
		fmt.Println("Failed to explore directory:", err)
	}
}

func exploreDirectory(dirPath, indent string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())

		fmt.Println(indent + file.Name())

		if file.IsDir() {
			err = exploreDirectory(filePath, indent+"  ")
			if err != nil {
				fmt.Printf("Failed to explore directory '%s': %v\n", filePath, err)
			}
		}
	}

	return nil
}
