package main

import (
	"go/ast"
	"go/token"
)

type TypeModifier int

const (
	nilTypeModifier TypeModifier = iota
	StarTypeModifier
	ArrayTypeModifier
)

type paramOption struct {
	VarName string
	VarType string
	TypeModifier
}

func makeFuncParamField(option paramOption) *ast.Field {
	f := &ast.Field{}
	if option.VarName != "" {
		f.Names = []*ast.Ident{
			ast.NewIdent(option.VarName),
		}
	}
	var typeExpr ast.Expr
	switch option.TypeModifier {
	case StarTypeModifier:
		typeExpr = ast.Expr(&ast.StarExpr{
			X: ast.Expr(ast.NewIdent(option.VarType)),
		})
	case ArrayTypeModifier:
		typeExpr = ast.Expr(&ast.ArrayType{
			Elt: ast.NewIdent(option.VarType),
		})
	default:
		typeExpr = ast.Expr(ast.NewIdent(option.VarType))
	}
	f.Type = typeExpr
	return f
}

func makeFuncParam(params []paramOption) *ast.FieldList {
	if len(params) == 0 {
		return nil
	}
	fields := &ast.FieldList{
		List: make([]*ast.Field, 0, len(params)),
	}
	for _, p := range params {
		f := makeFuncParamField(p)
		fields.List = append(fields.List, f)
	}
	return fields
}

func makeFuncDecl(pos token.Pos, methodName string,
	params []paramOption, retParam []paramOption,
	doc *ast.CommentGroup) *ast.FuncDecl {

	fnType := &ast.FuncType{
		Func:    pos,
		Params:  makeFuncParam(params),
		Results: makeFuncParam(retParam),
	}
	fnName := &ast.Ident{
		Name:    methodName,
		NamePos: pos + 1,
	}
	funcDecl := &ast.FuncDecl{
		Name: fnName,
		Doc:  doc,
		Body: &ast.BlockStmt{},
		Type: fnType,
	}
	return funcDecl
}
