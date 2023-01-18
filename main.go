package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	filename := "./docs/doc.yaml"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// if regex matches, then print the line
		match, _ := regexp.MatchString(`^  /.*:$`, scanner.Text())
		if match {
			processPath(scanner.Text())
		}

		if scanner.Text() == "components:" {
			break
		}
	}
}

func processPath(path) {
	fmt.Println("processing", path)
}
