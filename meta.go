package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Meta struct {
	Path  string `yaml:""`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`

	List bool `yaml:"list"`
}

func readMeta(path string) *Meta {
	f, err := os.Open(path)
	errorFatal("Error reading meta "+path, err)

	dec := yaml.NewDecoder(f)

	m := new(Meta)
	errorFatal("Error decoding meta "+path, dec.Decode(m))

	return m
}
