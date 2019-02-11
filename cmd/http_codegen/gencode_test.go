package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/format"
	"go/printer"
	"reflect"
	"testing"
)

func TestGenHttpAPIFuncDeclare(t *testing.T) {
	src := `
package controller

// this is a package comment

// @HttpAPI("/url", apiReq, apiResp)


/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/

`
	sa, err := ParseSourceCode("test", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)
	inspectNode(t, sa.node)

	markers, err := ParseHttpAPIMarkers(sa)

	genHttpAPIHandleFunc(sa, markers)

	inspectNode(t, sa.node)
	sourceBuf := bytes.NewBuffer(make([]byte, 0))
	printer.Fprint(sourceBuf, sa.fSet, sa.node)
	sourceCode, err := format.Source(sourceBuf.Bytes())
	assert.Nil(t, err)
	t.Log(string(sourceCode))
}

func inspectNode(t *testing.T, fNode *ast.File) {
	for _, decl := range fNode.Decls {
		t.Logf("%v", reflect.TypeOf(decl))
		fn := decl.(*ast.FuncDecl)
		t.Logf("name: %v", fn.Name)
		t.Logf("type: %v", fn.Type)
		for _, p := range fn.Type.Params.List {
			t.Logf("param: %v", p)
		}
		t.Logf("body: %v", fn.Body)
		t.Logf("recv: %v", fn.Recv)
		t.Logf("doc: %p,%v,%v", fn.Doc, fn.Doc, fn.Doc.Text())
	}
	for _, comment := range fNode.Comments {
		t.Logf("commentptr: %p", comment)
		t.Logf("comment: %v, %s", comment, comment.Text())
	}
}

func TestGenHttpAPITypeDecl(t *testing.T) {
	src := `
package controller

// this is a package comment

// @HttpAPI("/url", apiReq, apiResp)

/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/

`
	sa, err := ParseSourceCode("test", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)
	inspectNode(t, sa.node)

	markers, err := ParseHttpAPIMarkers(sa)

	genHttpAPIHandleFunc(sa, markers)
	genHttpAPITypeDecl(sa, markers)

	sourceBuf := bytes.NewBuffer(make([]byte, 0))
	printer.Fprint(sourceBuf, sa.fSet, sa.node)
	sourceCode, err := format.Source(sourceBuf.Bytes())
	assert.Nil(t, err)
	t.Log(string(sourceCode))
}
