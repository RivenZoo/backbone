package main

import (
	"encoding/json"
	"flag"
	"github.com/RivenZoo/backbone/logger"
	"os"
	"path/filepath"
	"strings"
)

var config *Config

type Config struct {
	GitRepoURL  string
	GitVersion  string
	LocalPath   string
	ProjectName string
	GoModName   string
	OutputDir   string
	TmplArgs    map[string]string
	Debug       bool
}

func (c *Config) String() string {
	data, err := json.Marshal(config)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func parseFlagConfig() {
	config = &Config{
		TmplArgs: make(map[string]string),
	}
	gitRepo := flag.String("gitRepo", "", "seed project git repo url, if not set LocalPath will be read")
	gitVer := flag.String("gitVer", "master", "git version tag, if no set master will be used")
	localPath := flag.String("localPath", "", "seed project local dir path")
	modName := flag.String("modName", "", "project mod full name, if set go mod will be used, eg github.com/RivenZoo/backbone")
	projName := flag.String("project", "", "project name, if not set it will be parsed from mod name")
	outputDir := flag.String("output", "./", "project dir, default current dir")
	debug := flag.Bool("debug", false, "debug flag")
	flag.Parse()

	config.GitRepoURL = *gitRepo
	config.GitVersion = *gitVer
	config.LocalPath = *localPath
	config.GoModName = *modName
	config.ProjectName = *projName
	config.OutputDir = *outputDir
	config.Debug = *debug

	if config.GitRepoURL == "" && config.LocalPath == "" {
		logger.Errorf("no input seed project")
		os.Exit(-1)
	}
	if config.GoModName == "" && config.ProjectName == "" {
		logger.Errorf("no mod or project name")
		os.Exit(-1)
	}

	parseProjectName()
	parseTmplArgs()
	addFlagToTmplArgs()
}

const tmplArgsKeyPrefix = "-tmpl_"

// parseTmplArgs parse command line args -tmpl_{key}={value}, set {key}={value} to TmplArgs
func parseTmplArgs() {
	n := flag.NArg()
	for i := 0; i < n; i++ {
		arg := flag.Arg(i)
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			logger.Debugf("skip args %s", arg)
			continue
		}
		key, val := kv[0], kv[1]
		if strings.HasPrefix(key, tmplArgsKeyPrefix) {
			config.TmplArgs[strings.TrimPrefix(key, tmplArgsKeyPrefix)] = val
		}
	}
}

func parseProjectName() {
	if config.ProjectName != "" {
		return
	}
	config.ProjectName = filepath.Base(config.GoModName)
}

func addFlagToTmplArgs() {
	config.TmplArgs["modName"] = config.GoModName
	config.TmplArgs["project"] = config.ProjectName
}
