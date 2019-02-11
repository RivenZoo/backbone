package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
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
