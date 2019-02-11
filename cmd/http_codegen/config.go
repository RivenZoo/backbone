package main

import "flag"

var config *Config

type Config struct {
	inputFile string
	inputDir  string
}

func parseFlagConfig() {
	config = &Config{}
	inputFile := flag.String("input", "", "input go source file")
	dir := flag.String("inputDir", "", "input go source directory")
	flag.Parse()

	config.inputFile = *inputFile
	config.inputDir = *dir
}
