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

func main2() {
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
	// path := scanner.Text()
	scanner.Scan()
	// method := scanner.Text()
	scanner.Scan()
	summary := scanner.Text()
	summary = strings.Replace(summary, "summary: ", "", 1)
	// keep scanning until we get to the tag
	for strings.TrimSpace(scanner.Text()) != "tags:" {
		scanner.Scan()
	}
	scanner.Scan()
	tag := scanner.Text()

	filename := fmt.Sprintf(
		"%s/%s/%s/%s.yaml",
		ROOT_DIR,
		PATHS_DIR,
		sanitize(tag),
		sanitize(summary))

	fmt.Println("filename", filename)
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`-`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`:`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`'`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`/`).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`{`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`}`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`\(`).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`\)`).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(s, func(r string) string {
		return strings.ToLower(r)
	})
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(` `).ReplaceAllString(s, "_")
	s = regexp.MustCompile(`__`).ReplaceAllString(s, "_")
	// gets rid of any non-alphanumeric characters or underscores
	s = regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(s, "")

	return s
}
