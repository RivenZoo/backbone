package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenHttpAPIDeclare(t *testing.T) {
	src := `
package controller

import (
	
)

// @HttpAPI("/url", apiReq, apiResp)

/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/
`
	g := newHttpAPIGenerator(httpAPIGeneratorOption{
		apiDefineFileImports: []importInfo{{"github.com/gin-gonic/gin", ""},
			{"github.com/request/header", ""}},
		commonAPIDefinition: commonHttpAPIDefinition{
			CommonRequestFields:  "header.RequestHeader",
			CommonResponseFields: "header.ResponseHeader",
			CommonFuncStmt:       "// set common code snippet here",
		},
	})
	err := g.ParseCode("test.go", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)

	err = g.ParseHttpAPIMarkers()
	assert.Nil(t, err)

	assert.True(t, len(g.markers) == 2)

	g.GenHttpAPIDeclare()
	codeBuf := bytes.NewBuffer(make([]byte, 0))
	g.OutputAPIDeclare(codeBuf)

	t.Log(codeBuf.String())
}
