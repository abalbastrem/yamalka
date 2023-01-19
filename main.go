package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const ROOT_DIR = "docs"
const PATHS_DIR = "paths"
const DIR_SEP = "/"

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
			processPath(scanner)
		}

		if scanner.Text() == "components:" {
			break
		}
	}
}

func processPath(scanner *bufio.Scanner) {
	path := scanner.Text()
	scanner.Scan()
	method := scanner.Text()
	scanner.Scan()
	scanner.Scan()
	scanner.Scan()
	tagDir := scanner.Text()

	filename := fmt.Sprintf(
		"%s/%s/%s/%s_%s.yaml",
		ROOT_DIR,
		PATHS_DIR,
		sanitize(tagDir),
		sanitize(path),
		sanitize(method))
	fmt.Println("filename", filename)
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`-`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`:`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`/`).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`{`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`}`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(s, func(r string) string {
		return strings.ToLower(r)
	})
	s = strings.TrimSpace(s)

	return s
}
