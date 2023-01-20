package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const ROOT_DIR = "docs"
const PATHS_DIR = "paths"
const DIR_SEP = "/"

const DOC_ORIGINAL = "doc"
const DOC_NEW = "main"
const DOC_EXTENSION = ".yaml"

func main() {

	doc_original := make(map[string]interface{})
	doc_new := make(map[string]interface{})
	ymlContent, err := os.ReadFile(ROOT_DIR + DIR_SEP + DOC_ORIGINAL + DOC_EXTENSION)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(ymlContent, &doc_original)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(ymlContent, &doc_new)
	if err != nil {
		panic(err)
	}

	for path, vMethod := range doc_original["paths"].(map[string]interface{}) {
		var method, filename, tag string
		var v interface{}

		fmt.Println()
		fmt.Println("PATH\t", path)
		for method, v = range vMethod.(map[string]interface{}) {
			if method == "$ref" { // already processed
				continue
			}
			method = sanitize(method)
			fmt.Println("METHOD\t", method)
			for k, v := range v.(map[string]interface{}) {
				if k == "summary" {
					filename = sanitize(v.(string))
					fmt.Println("SUMMARY\t", filename)
				}
				if k == "tags" {
					tagArr, _ := v.([]interface{})
					tag = sanitize(tagArr[0].(string))
					fmt.Println("TAG\t", tag)
				}
			}
			docPaths := doc_new["paths"].(map[string]interface{})
			docPath := docPaths[path].(map[string]interface{})
			delete(docPath, method)
			docPath["$ref"] = PATHS_DIR + DIR_SEP + tag + DIR_SEP + filename + DOC_EXTENSION
		}

		dir := fmt.Sprintf(
			"%s/%s/%s",
			ROOT_DIR,
			PATHS_DIR,
			tag)

		fmt.Println("FILENAME", dir+DIR_SEP+filename+DOC_EXTENSION)

		pathContents := replaceRefs(vMethod.(map[string]interface{}))
		// writes out subdivided files
		streamOut(dir, filename, pathContents)
	}

	// writes out the main doc
	fmt.Println("\nMAIN DOC IN ", ROOT_DIR+DIR_SEP+DOC_NEW+DOC_EXTENSION)
	newDoc := NewDocument(doc_new)
	streamOut(ROOT_DIR, DOC_NEW, newDoc)

	fmt.Println("\nDONE")
}

func replaceRefs(a any) any {
	val, ok := a.(map[string]interface{})
	if ok {
		for k, v := range val {
			if k == "$ref" {
				val[k] = strings.Replace(v.(string), "#/", "../../"+DOC_EXTENSION+"#/", 1)
				return val
			}
			val[k] = replaceRefs(v)
		}
		return val
	}

	return a
}

func streamOut(dir string, filename string, data any) {
	// creates dir if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	// create writer to file
	f, err := os.Create(dir + DIR_SEP + filename + DOC_EXTENSION)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := yaml.NewEncoder(f)
	e.SetIndent(2)
	e.Encode(data)
	defer e.Close()
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
	// gets rid of final non-alphanumeric character, if any
	s = regexp.MustCompile(`[^a-zA-Z0-9]$`).ReplaceAllString(s, "")

	return s
}

type Document struct {
	Openapi    string `yaml:"openapi"`
	Info       any    `yaml:"info"`
	Servers    any    `yaml:"servers"`
	Paths      any    `yaml:"paths"`
	Tags       any    `yaml:"tags"`
	Components any    `yaml:"components"`
}

func NewDocument(data map[string]interface{}) *Document {
	return &Document{
		Openapi:    data["openapi"].(string),
		Info:       data["info"].(map[string]interface{}),
		Servers:    data["servers"].([]interface{}),
		Paths:      data["paths"].(map[string]interface{}),
		Tags:       data["tags"].([]interface{}),
		Components: data["components"].(map[string]interface{}),
	}
}
