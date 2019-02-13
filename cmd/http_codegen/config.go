package main

import (
	"flag"
	"github.com/RivenZoo/backbone/configutils"
	"github.com/RivenZoo/backbone/logger"
	"os"
)

var config *Config
var genOption *httpAPIGeneratorOption

type Config struct {
	inputFile     string
	inputDir      string
	genConfigFile string
	debug         bool
}

func parseFlagConfig() {
	config = &Config{}
	inputFile := flag.String("input", "", "input go source file")
	dir := flag.String("inputDir", "", "input go source directory")
	genCfgFile := flag.String("genCfg", "", "generator config file")
	debug := flag.Bool("debug", false, "debug flag")
	flag.Parse()

	config.inputFile = *inputFile
	config.inputDir = *dir
	config.genConfigFile = *genCfgFile
	config.debug = *debug

	if config.inputFile == "" && config.inputDir == "" {
		logger.Errorf("no input")
		os.Exit(-1)
	}
	if config.genConfigFile != "" {
		cfgTp := configutils.DetechFileConfigType(config.genConfigFile)
		if cfgTp == "" {
			cfgTp = configutils.ConfigTypeJSON
		}
		if err := configutils.UnmarshalFile(config.genConfigFile, &genOption, cfgTp); err != nil {
			logger.Errorf("UnmarshalFile %s error %v", config.genConfigFile, err)
			os.Exit(-1)
		}
	} else {
		genOption = &httpAPIGeneratorOption{}
	}
}
