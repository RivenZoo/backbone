package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"reflect"
	"testing"
)

func TestFilterStmt(t *testing.T) {
	src := `package controller

func foo(r *gin.Engine) {
	r.Post("/url", handler)
	r.Post("/url2", handler2)
}`
	sa, err := ParseSourceCode("tmp", bytes.NewBufferString(src))
	assert.Nil(t, err)

	fooDecl := filterGlobalFunc(sa.node.Decls, "foo")
	assert.NotNil(t, fooDecl)

	ast.Inspect(fooDecl, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selectExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if ok {
			t.Logf("%v.%v", selectExpr.X, selectExpr.Sel.Name)
		}
		for _, arg := range callExpr.Args {
			t.Logf("node %v, %v", reflect.TypeOf(arg), arg)
		}
		return true
	})

}
