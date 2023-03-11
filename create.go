package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func init() {
	functionsMu.Lock()
	defer functionsMu.Unlock()

	functions["create"] = funcCreate
}

func funcCreate() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: swb create <folder>\n")
		os.Exit(1)
	}

	_path := os.Args[2]

	Mkdir(_path)
	conf := writeDefaultConf(path.Join(_path, "config.yml"))

	Mkdir(path.Join(_path, conf.Templates))
	Mkdir(path.Join(_path, conf.Output))
	Mkdir(path.Join(_path, conf.Input))
	Mkdir(path.Join(_path, conf.Input, "00-test"))

	// write example templates:
	WriteFile(path.Join(_path, conf.Templates, "index.md"),
		`{{range .}}
{{if .List}}
---

# [{{.Title}}]({{.Path}}/)

{{.Desc}}
{{end}}
{{end}}
`)

	//create default template
	WriteFile(path.Join(_path, conf.Templates, "filter.lua"), "")
	// write default pandoc.html

	cmd := exec.Command("pandoc",
		"-o"+path.Join(_path, conf.Templates, "pandoc.html"),
		"-D", "html")

	err := cmd.Run()
	errorFatal("Failed to run pandoc", err)

	WriteFile(path.Join(_path, conf.Input, "00-test", "index.md"),
		`# Headline

Lorem ipsum dolor sit!`)

	WriteFile(path.Join(_path, conf.Input, "00-test", "meta.yml"),
		`title: "Title"
desc: "Description"

list: true`)
}

func WriteFile(path, content string) {
	err := os.WriteFile(path, []byte(content), 0o0755)
	errorFatal("Failed to create "+path, err)
}

func Mkdir(path string) {
	err := os.MkdirAll(path, 0o755)
	errorFatal("Failed to mkdir "+path, err)
}
