package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"fmt"
	"sync"
)

var (
	functions   = make(map[string]func())
	functionsMu sync.RWMutex
)

func init() {
	functionsMu.Lock()
	defer functionsMu.Unlock()

	functions["help"] = func() { helpS(functions, &functionsMu) }
}

func main() {
	if len(os.Args) < 2 {
		helpS(functions, &functionsMu)
		os.Exit(1)
	}

	f, ok := functions[os.Args[1]]
	if !ok {
		helpS(functions, &functionsMu)
		os.Exit(1)
	}

	f()
}

func openTemplate(name, file string) *template.Template {
	f, err := os.Open(file)
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

func helpS(m map[string]func(), mut *sync.RWMutex) {
	fmt.Printf("Usage: swb <subcommand>\n\nsubcommand := {%s }\n", keysS(m, mut))
}

func keysS(m map[string]func(), mut *sync.RWMutex) (str string) {
	mut.RLock()
	defer mut.RUnlock()

	for k := range m {
		str += " " + k
	}

	return
}
