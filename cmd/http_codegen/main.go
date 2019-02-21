package main

import (
	"github.com/RivenZoo/backbone/logger"
	"io"
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
	files := []string{}
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Errorf("read path %s error %v", path, err)
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "_handlers.go") ||
			strings.HasSuffix(path, "_urls.go") ||
			strings.HasSuffix(path, "_test.go") {
			return nil
		}
		files = append(files, path)
		return nil
	})
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

	handlerFile := apiHandlerFileName(filePath)
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
	g.ParseFile(filePath)
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
	outputCode(apiHandlerFileName(filePath), g.OutputAPIHandler)

	g.GenInitHttpAPIRouter()
	debugOutput("init_router", g.routerInitOutput)
	outputCode(g.httpRouterInitFilename(), g.OutputInitHttpAPIRouter)
}

func outputCode(filePath string, writeFunc func(w io.Writer) error) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("open file %s error %v", filePath, err)
		return
	}
	err = writeFunc(f)
	if err != nil {
		logger.Errorf("writeFunc %s error %v", filePath, err)
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
