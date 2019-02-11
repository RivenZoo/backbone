package main

import (
	"github.com/RivenZoo/backbone/logger"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	parseFlagConfig()
	logger.Debugf("%v", *config)
	err := parseSourceFile2(config.inputFile)
	if err != nil {
		logger.Infof("parse error %v", err)
	}
}

func parseSourceFile2(srcFile string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, cg := range f.Comments {
		for _, c := range cg.List {
			logger.Infof("comment: %s", c.Text)
			logger.Infof("position: %v", fset.Position(c.Pos()))
		}
	}
	logger.Infof("identity: %v", f.Name)
	for _, decl := range f.Decls {
		switch sd := decl.(type) {
		case *ast.GenDecl:
			for _, sp := range sd.Specs {
				switch realSp := sp.(type) {
				case *ast.ImportSpec:
					logger.Infof("gen decl %v", realSp.Path)
				case *ast.ValueSpec:
					logger.Infof("gen decl %v", realSp.Names)
				case *ast.TypeSpec:
					logger.Infof("gen decl %v", realSp.Name)
				}
			}
		case *ast.FuncDecl:
			logger.Infof("func decl %v", sd.Name)
		}
	}
	return nil
}
