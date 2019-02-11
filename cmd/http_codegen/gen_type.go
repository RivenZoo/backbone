package main

import (
	"go/ast"
	"go/token"
)

func makeStructTypeDecl(pos token.Pos, typeName string) *ast.GenDecl {
	tpDecl := &ast.GenDecl{
		TokPos: pos,
		Tok:    token.TYPE,
		Specs: []ast.Spec{
			ast.Spec(&ast.TypeSpec{
				Name: ast.NewIdent(typeName),
				Type: ast.Expr(&ast.StructType{
					Fields:     &ast.FieldList{},
				}),
			}),
		},
	}
	return tpDecl
}
