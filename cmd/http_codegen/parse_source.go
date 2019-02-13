package main

import (
	"github.com/RivenZoo/backbone/logger"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

type SourceAst struct {
	filePath string
	fSet     *token.FileSet
	node     *ast.File
}

func ParseSourceFile(srcFile string) (*SourceAst, error) {
	sa := &SourceAst{
		filePath: srcFile,
		fSet:     token.NewFileSet(),
	}
	node, err := parser.ParseFile(sa.fSet, srcFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	sa.node = node
	return sa, nil
}

func ParseSourceCode(filePath string, code io.Reader) (*SourceAst, error) {
	sa := &SourceAst{
		filePath: filePath,
		fSet:     token.NewFileSet(),
	}
	node, err := parser.ParseFile(sa.fSet, filePath, code, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	sa.node = node
	return sa, nil
}

func readSource(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		logger.Errorf("open file %s error %v", filePath, err)
		return nil, err
	}
	return data, nil
}
