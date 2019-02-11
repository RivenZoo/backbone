package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/printer"
	"os"
	"testing"
)

func TestMakeFuncParamField(t *testing.T) {
	src := `
package controller

// @HttpAPI("/url", apiReq, apiResp)

/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/
func handleApiReq2()	{}
`
	sa, err := ParseSourceCode("test", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)
	t.Log(sa.node)

	f1 := makeFuncDecl(sa.node.Comments[0].End()+1, "foo", []paramOption{
		paramOption{
			VarName:      "data",
			VarType:      "byte",
			TypeModifier: ArrayTypeModifier,
		},
	}, []paramOption{
		paramOption{
			VarType: "error",
		},
	}, sa.node.Comments[0])
	f2 := makeFuncDecl(0, "foo2", []paramOption{
		paramOption{
			VarName:      "req",
			VarType:      "testReq",
			TypeModifier: StarTypeModifier,
		},
	}, []paramOption{
		paramOption{
			VarName:      "resp",
			VarType:      "testResp",
			TypeModifier: StarTypeModifier,
		},
		paramOption{
			VarName: "err",
			VarType: "error",
		},
	}, nil)

	//t.Logf("pos %v %v %v", f1.Type.Pos(), f1.Name.Pos(), f1.Doc.Pos())
	t.Logf("f1 %#v", f1)
	t.Logf("f1 %v %v %v", f1.Type.Pos(), f1.Name.Pos(), f1.Doc.Pos())
	//t.Logf("pos %v %v %v", f2.Type.Pos(), f2.Name.Pos(), f2.Doc.Pos())
	fns := []*ast.FuncDecl{f1, f2}
	for _, fn := range fns {
		sa.node.Decls = append(sa.node.Decls, fn)
	}
	printer.Fprint(os.Stdout, sa.fSet, sa.node)
}
