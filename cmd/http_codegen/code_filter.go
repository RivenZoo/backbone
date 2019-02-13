package main

import (
	"go/ast"
	"go/token"
	"strings"
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

func filterUnInitRouters(sa *SourceAst, initFuncName, registerFuncName string, markers []*HttpAPIMarker) []initRouterStmtInfo {
	initDecl := filterGlobalFunc(sa.node.Decls, initFuncName)
	if initDecl == nil {
		// no init func
		ret := make([]initRouterStmtInfo, 0, len(markers))
		for _, m := range markers {
			ret = append(ret, initRouterStmtInfo{
				marker:          m,
				afterLine:       sa.fSet.Position(sa.node.End()).Line,
				handlerFuncName: httpAPIMethodName(m.RequestType),
			})
		}
		return ret
	}

	registerUrls := make([]string, 0)
	ast.Inspect(initDecl, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selectExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if selectExpr.Sel.Name != registerFuncName {
			return true
		}
		if len(callExpr.Args) > 0 {
			urlLit, ok := callExpr.Args[0].(*ast.BasicLit)
			if ok {
				registerUrls = append(registerUrls, urlLit.Value)
			}
		}
		return true
	})
	findRegisterStmt := func(url string) bool {
		for _, u := range registerUrls {
			if strings.Trim(url, `"`) == strings.Trim(u, `"`) {
				return true
			}
		}
		return false
	}
	ret := make([]initRouterStmtInfo, 0, len(markers))
	for _, m := range markers {
		if !findRegisterStmt(m.URL) {
			ret = append(ret, initRouterStmtInfo{
				afterLine:       sa.fSet.Position(initDecl.End()).Line - 1,
				marker:          m,
				handlerFuncName: httpAPIMethodName(m.RequestType),
			})
		}
	}
	return ret
}

type funcCodeInfo struct {
	funcName  string
	beginLine int
	endLine   int
}

func filterGlobalFunc(decls []ast.Decl, funcName string) *ast.FuncDecl {
	for _, decl := range decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == funcName {
			return funcDecl
		}
	}
	return nil
}

func getFuncCodeInfo(fset *token.FileSet, funcDecl *ast.FuncDecl) *funcCodeInfo {
	if fset == nil || funcDecl == nil {
		return nil
	}
	return &funcCodeInfo{
		funcName:  funcDecl.Name.Name,
		beginLine: fset.Position(funcDecl.Pos()).Line,
		endLine:   fset.Position(funcDecl.End()).Line,
	}
}
