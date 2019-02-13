package main

import (
	"go/ast"
	"go/token"
)

func genHttpAPIFuncDeclare(targetAst *SourceAst, unDeclaredMarkers []*HttpAPIMarker) {
	for _, m := range unDeclaredMarkers {
		methodName := httpAPIMethodName(m.RequestType)
		params := []paramOption{
			paramOption{
				VarName:      "c",
				VarType:      "gin.Context",
				TypeModifier: StarTypeModifier,
			},
			paramOption{
				VarName:      "req",
				VarType:      m.RequestType,
				TypeModifier: StarTypeModifier,
			},
		}
		retParams := []paramOption{
			paramOption{
				VarName:      "resp",
				VarType:      m.ResponseType,
				TypeModifier: StarTypeModifier,
			},
			paramOption{
				VarName: "err",
				VarType: "error",
			},
		}
		funcDecl := makeFuncDecl(1, methodName,
			params, retParams, nil)
		targetAst.node.Decls = append(targetAst.node.Decls, funcDecl)
	}
	targetAst.node.Comments = nil
}

// genHttpAPIHandleFunc generate api declare if api not exists.
func genHttpAPIHandleFunc(targetAst *SourceAst, markers []*HttpAPIMarker) {
	declaredFuncs := filterDelcaredFuncNames(targetAst.node.Decls)
	unDeclaredMarkers := filterFuncUndeclaredHttpAPIMarkers(markers, declaredFuncs)
	genHttpAPIFuncDeclare(targetAst, unDeclaredMarkers)
}

func filterDeclaredTypeNames(decls []ast.Decl) []string {
	declaredTypeNames := make([]string, 0)
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				tpSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				declaredTypeNames = append(declaredTypeNames, tpSpec.Name.Name)
			}
		}
	}
	return declaredTypeNames
}

type undeclaredHttpAPIMarkerType struct {
	marker       *HttpAPIMarker
	requestType  string // "" means no need to declare
	responseType string // "" means no need to declare
}

func filterUndeclaredHttpAPITypes(markers []*HttpAPIMarker, declaredTypes []string) []undeclaredHttpAPIMarkerType {
	unDeclaredTypes := make([]undeclaredHttpAPIMarkerType, 0)
	findUndeclared := func(tp string) bool {
		for _, declared := range declaredTypes {
			if tp == declared {
				return true
			}
		}
		return false
	}
	for _, m := range markers {
		tp := undeclaredHttpAPIMarkerType{
			marker: m,
		}
		found := findUndeclared(m.RequestType)
		if !found {
			tp.requestType = m.RequestType
		}
		found = findUndeclared(m.ResponseType)
		if !found {
			tp.responseType = m.ResponseType
		}
		if tp.requestType != "" || tp.responseType != "" {
			unDeclaredTypes = append(unDeclaredTypes, tp)
		}
	}
	return unDeclaredTypes
}

func genTypeDeclares(targetAst *SourceAst, undeclaredTypes []undeclaredHttpAPIMarkerType) {
	decls := make([]ast.Decl, 0)
	for _, ut := range undeclaredTypes {
		if ut.requestType != "" {
			typeDecl := makeStructTypeDecl(1, ut.requestType)
			if typeDecl != nil {
				decls = append(decls, typeDecl)
			}
		}
		if ut.responseType != "" {
			typeDecl := makeStructTypeDecl(2, ut.responseType)
			if typeDecl != nil {
				decls = append(decls, typeDecl)
			}
		}
	}
	decls = append(decls, targetAst.node.Decls...)
	targetAst.node.Decls = decls
}

// genHttpAPITypeDecl generate request/response type declare if not exists.
func genHttpAPITypeDecl(targetAst *SourceAst, markers []*HttpAPIMarker) {
	typeNames := filterDeclaredTypeNames(targetAst.node.Decls)
	undeclaredTypes := filterUndeclaredHttpAPITypes(markers, typeNames)
	genTypeDeclares(targetAst, undeclaredTypes)
}

func copyCommentGroup(cg *ast.CommentGroup) *ast.CommentGroup {
	if cg == nil {
		return nil
	}
	ret := &ast.CommentGroup{
		List: []*ast.Comment{},
	}
	for _, c := range cg.List {
		ret.List = append(ret.List, &ast.Comment{
			Text: c.Text,
		})
	}
	return ret
}


