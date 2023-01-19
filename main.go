package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const ROOT_DIR = "docs"
const PATHS_DIR = "paths"
const DIR_SEP = "/"

func main() {

	doc_original := make(map[string]interface{})
	ymlContent, err := ioutil.ReadFile("./docs/doc.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(ymlContent, &doc_original)
	if err != nil {
		panic(err)
	}

	for path, vMethod := range doc_original["paths"].(map[string]interface{}) {
		var method, summary, tag string
		var v interface{}

		fmt.Println()
		path = sanitize(path)
		fmt.Println(path)
		for method, v = range vMethod.(map[string]interface{}) {
			if method == "$ref" { // already processed
				continue
			}
			method = sanitize(method)
			fmt.Println(method)
			for k, v := range v.(map[string]interface{}) {
				if k == "summary" {
					summary = sanitize(v.(string))
					fmt.Println("SUMMARY", summary)
				}
				if k == "tags" {
					tagArr, _ := v.([]interface{})
					tag = sanitize(tagArr[0].(string))
					fmt.Println("TAG", tag)
				}
			}
		}

		filename := fmt.Sprintf(
			"%s/%s/%s/%s.yaml",
			ROOT_DIR,
			PATHS_DIR,
			tag,
			summary)

		fmt.Println("FILENAME", filename)

		writeOutPathData(vMethod.(map[string]interface{}), filename)
		// TODO - remove the path from the original doc
		// delete(doc_original["paths"].(map[string]interface{}), path)
		// TODO - write out the original doc
		// writeOutPathData(doc_original, "./docs/doc.yaml")
		// TODO - tabs as two spaces
	}
}

func writeOutPathData(pathData map[string]interface{}, filename string) {
	// TODO bugfix
	yml, err := yaml.Marshal(pathData)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, yml, 0644)
	if err != nil {
		panic(err)
	}
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
