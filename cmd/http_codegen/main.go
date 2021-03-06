package main

import (
	"bytes"
	"github.com/RivenZoo/backbone/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	parseFlagConfig()

	if config.debug {
		logger.SetLogLevel(logger.DEBUG)
	} else {
		logger.SetLogLevel(logger.INFO)
	}

	files := []string{}
	if config.inputDir != "" {
		files = listHttpAPIFiles(config.inputDir)
	} else {
		files = append(files, config.inputFile)
	}

	logger.Debugf("input files %v", files)
	for _, fpath := range files {
		handleSourceFile(fpath)
	}
}

func listHttpAPIFiles(inputDir string) []string {
	fileInfos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		logger.Errorf("read path %s error %v", inputDir, err)
		return nil
	}

	files := make([]string, 0, len(fileInfos))
	for _, info := range fileInfos {
		if info.IsDir() {
			continue
		}
		path := filepath.Join(inputDir, info.Name())
		if strings.HasSuffix(path, "_handlers.go") ||
			strings.HasSuffix(path, "_urls.go") ||
			strings.HasSuffix(path, "_test.go") {
			continue
		}
		files = append(files, path)
	}

	return files
}

func sourceModifiedSinceLastGen(filePath string, g *HttpAPIGenerator) bool {
	info, err := os.Lstat(filePath)
	if err != nil {
		// ignore error, treat as modified
		logger.Debugf("Lstat file %s error %v", filePath, err)
		return true
	}
	logger.Debugf("file %s mod time %s", filePath, info.ModTime())

	handlerFile := g.apiHandlerFileName(filePath)
	handlerFileInfo, err := os.Lstat(handlerFile)
	if err != nil {
		// ignore error, treat as modified
		logger.Debugf("Lstat file %s error %v", handlerFile, err)
		return true
	}
	logger.Debugf("file %s mod time %s", handlerFile, handlerFileInfo.ModTime())

	if info.ModTime().After(handlerFileInfo.ModTime()) {
		return true
	}

	initRouterFile := g.httpRouterInitFilename()
	initRouterFileInfo, err := os.Lstat(initRouterFile)
	if err != nil {
		// ignore error, treat as modified
		logger.Debugf("Lstat file %s error %v", handlerFile, err)
		return true
	}
	logger.Debugf("file %s mod time %s", initRouterFile, initRouterFileInfo.ModTime())

	if info.ModTime().After(initRouterFileInfo.ModTime()) {
		return true
	}
	return false
}

func handleSourceFile(filePath string) {
	g := newHttpAPIGenerator(*genOption)
	if err := g.ParseFile(filePath); err != nil {
		logger.Errorf("parse file %s error %v", filePath, err)
		return
	}
	if err := g.ParseHttpAPIMarkers(); err != nil {
		logger.Errorf("ParseHttpAPIMarkers error %v", err)
		return
	}
	if len(g.markers) == 0 {
		return
	}

	if !sourceModifiedSinceLastGen(filePath, g) {
		logger.Debugf("source file %s not modified since last generate", filePath)
		return
	}

	g.GenHttpAPIDeclare()
	debugOutput("api_declare", g.sourceFileOutput)
	outputCode(filePath, g.OutputAPIDeclare)

	g.GenHttpAPIHandler()
	debugOutput("api_handler", g.handlerOutput)
	outputCode(g.apiHandlerFileName(filePath), g.OutputAPIHandler)

	g.GenInitHttpAPIRouter()
	debugOutput("init_router", g.routerInitOutput)
	outputCode(g.httpRouterInitFilename(), g.OutputInitHttpAPIRouter)
}

func outputCode(filePath string, writeFunc func(w io.Writer) error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := writeFunc(buf)
	if err != nil {
		logger.Errorf("writeFunc %s error %v", filePath, err)
		return
	}
	err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		logger.Errorf("write to file %s error %v", filePath, err)
		return
	}
}

func debugOutput(name string, outputs []generatedOutput) {
	if !config.debug {
		return
	}
	logger.Debugf("%s output %d", name, len(outputs))
	for _, output := range outputs {
		logger.Debugf("%s:%d %s", name, output.afterLine, output.buffer.String())
	}
}
