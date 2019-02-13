package main

import (
	"github.com/RivenZoo/backbone/logger"
	"io"
	"os"
)

func main() {
	parseFlagConfig()

	if config.debug {
		logger.SetLogLevel(logger.DEBUG)
	}

	files := []string{}
	if config.inputDir != "" {
		files = listHttpAPIFiles(config.inputDir)
	} else {
		files = append(files, config.inputFile)
	}
	for _, fpath := range files {
		handleSourceFile(fpath)
	}
}

func listHttpAPIFiles(inputDir string) []string {
	return nil
}

func handleSourceFile(filePath string) {
	g := newHttpAPIGenerator(*genOption)
	g.ParseFile(filePath)
	if err := g.ParseHttpAPIMarkers(); err != nil {
		logger.Errorf("ParseHttpAPIMarkers error %v", err)
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
