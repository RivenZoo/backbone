package main

import (
	"go/ast"
	"go/token"
)

func filterDelcaredFuncNames(decls []ast.Decl) []string {
	declaredFuncs := make([]string, 0)
	for _, decl := range decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		declaredFuncs = append(declaredFuncs, funcDecl.Name.Name)
	}
	return declaredFuncs
}

func filterFuncUndeclaredHttpAPIMarkers(markers []*HttpAPIMarker, declaredFuncs []string) []*HttpAPIMarker {
	unDeclaredMarkers := make([]*HttpAPIMarker, 0, len(markers)-len(declaredFuncs))
	findUndeclared := func(fn string) bool {
		for _, declared := range declaredFuncs {
			if fn == declared {
				return true
			}
		}
		return false
	}
	for _, m := range markers {
		methodName := httpAPIMethodName(m.RequestType)
		found := findUndeclared(methodName)
		if !found {
			unDeclaredMarkers = append(unDeclaredMarkers, m)
		}
	}
	return unDeclaredMarkers
}

func filterUnImportedPackage(imports []*ast.ImportSpec, pkgs []importInfo) (unImported []importInfo) {
	unImported = make([]importInfo, 0)

	findImported := func(pkg string) bool {
		for _, imp := range imports {
			if pkg == imp.Path.Value {
				return true
			}
		}
		return false
	}
	for _, pkg := range pkgs {
		if !findImported(pkg.PkgPath) {
			unImported = append(unImported, pkg)
		}
	}
	return
}

func filterUndeclaredHandlers(sa *SourceAst, markers []*HttpAPIMarker) []apiHandlerDefineInfo {
	lastLine := 0
	globalVars := []string{}
	for _, decl := range sa.node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			varSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			globalVars = append(globalVars, varSpec.Names[0].Name)
			lastLine = sa.fSet.Position(varSpec.End()).Line
		}
	}
	findDeclaredVar := func(varName string) bool {
		for _, v := range globalVars {
			if varName == v {
				return true
			}
		}
		return false
	}

	ret := make([]apiHandlerDefineInfo, 0)
	for i, m := range markers {
		vn := httpAPIHandlerVarName(m.RequestType)
		if !findDeclaredVar(vn) {
			ret = append(ret, apiHandlerDefineInfo{
				marker:        m,
				varName:       vn,
				afterLine:     lastLine + i,
				apiMethodName: httpAPIMethodName(m.RequestType),
			})
		}
	}
	return ret
}
