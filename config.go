package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type Config struct {
	Templates string `yaml:"templates"`
	Output    string `yaml:"output"`
	Input     string `yaml:"input"`
}

func readConf(_path string) *Config {
	_path = path.Join(_path, "config.yml")
	f, err := os.Open(_path)
	errorFatal("Error reading config from "+_path, err)
	defer f.Close()

	dec := yaml.NewDecoder(f)

	m := new(Config)
	errorFatal("Error decoding config from "+_path, dec.Decode(m))

	return m
}

func writeDefaultConf(path string) *Config {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o755)
	errorFatal("Error reading config from "+path, err)
	defer f.Close()

	dec := yaml.NewEncoder(f)

	conf := &Config{
		Templates: "templates/",
		Output:    "out/",
		Input:     "src/",
	}
	errorFatal("Error encoding config to "+path, dec.Encode(conf))

	return conf
}
