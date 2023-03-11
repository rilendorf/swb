package main

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"text/template"

	"flag"
	"fmt"
)

func main() {
	templates := flag.String("templates", "templates/", "overwrite directory containing templates")
	output := flag.String("out", "out/", "overwrite output directory")
	input := flag.String("src", "src/", "overwrite source directory")

	flag.Parse()

	dir, err := os.ReadDir(*input)
	errorFatal("Error reading input dir "+*input, err)

	var entries []*Meta

	// read all the projects
	for _, e := range dir {
		//ignore files
		if !e.IsDir() {
			continue
		}

		//open metadata file
		meta := readMeta(path.Join(*input, e.Name(), "meta.yml"))
		if !meta.List {
			continue
		}

		meta.Path = path.Join(e.Name())

		if meta.List {
			entries = append(entries, meta)
		}
	}

	index := openTemplate(os.DirFS(*templates), "index", "index.md")

	/* index */
	tmpindex := path.Join(*output, "index.tmp.md")

	fmt.Printf("Writing main %s\n", tmpindex)
	indexFile, err := os.OpenFile(tmpindex, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	errorFatal("Error opening "+tmpindex, err)
	defer indexFile.Close()

	err = index.Execute(indexFile, entries)
	errorFatal("Error executing index template", err)

	// do the markdown
	fmt.Printf("Writing main index.html\n")
	pandoc(path.Join(*output, "index.html"),
		path.Join(*output, "index.tmp.md"),
		fmt.Sprintf("--template=%s", path.Join(*templates, "pandoc.html")),
		fmt.Sprintf("--lua-filter=%s", path.Join(*templates, "filter.lua")),
		"--css=/style.css",
		"-f", "markdown+smart",
		"--to=html5",
	)

	/* generate all markdown files */
	for _, entry := range entries {
		dirpath := path.Join(*output, entry.Path)

		err = os.MkdirAll(dirpath, 0755)
		if err != nil {
			fmt.Printf("Error during mkdir '%s':%s\n", dirpath, err)
			continue
		}

		dir, err := os.ReadDir(path.Join(*input, entry.Path))
		if err != nil {
			fmt.Printf("Error reading dir '%s':%s\n", dirpath, err)
			continue
		}

		for _, file := range dir {
			if file.IsDir() {
				continue
			}

			name := file.Name()
			ext := filepath.Ext(name)
			name = name[0 : len(name)-len(ext)]

			extMap := map[string]struct{}{
				".md": struct{}{},
				//"htm":  struct{}{},
			}

			if _, ok := extMap[ext]; !ok {
				continue
			}

			fmt.Printf("Writing %s: %s\n", entry.Path, file.Name())

			pandoc(path.Join(*output, entry.Path, name+".html"),
				path.Join(*input, entry.Path, file.Name()),
				fmt.Sprintf("--template=%s", path.Join(*templates, "pandoc.html")),
				fmt.Sprintf("--lua-filter=%s", path.Join(*templates, "filter.lua")),
				"--css=/style.css",
				"-f", "markdown+smart",
				"--to=html5",
			)
		}
	}
}

func openTemplate(fs fs.FS, name, file string) *template.Template {
	f, err := fs.Open(file)
	errorFatal("Error opening template file "+file, err)
	defer f.Close()

	tmpl := template.New(name)
	slice, err := io.ReadAll(f)
	errorFatal("Error reading "+file, err)

	tmpl, err = tmpl.Parse(string(slice))
	errorFatal("Error parsing "+file, err)

	return tmpl
}

func errorFatal(f string, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		_, file = filepath.Split(file)

		fmt.Printf("%s:%d - %s: %s!\n", file, line, f, err)
		os.Exit(1)
	}
}
