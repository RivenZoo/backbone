package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/format"
	"testing"
)

func TestGenHttpAPIDeclare(t *testing.T) {
	src := `
package controller

import (
	"github.com/RivenZoo/sqlagent"
)

// @HttpAPI("/url", apiReq, apiResp)

/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/
`
	g := newHttpAPIGenerator(httpAPIGeneratorOption{
		imports: []importInfo{{"github.com/gin", ""},
			{"github.com/request/header", ""}},
		commonAPIDefinition: commonHttpAPIDefinition{
			CommonRequestFields:  "header.RequestHeader",
			CommonResponseFields: "header.ResponseHeader",
			CommonFuncStmt:       "// set common code snippet here",
		},
	})
	err := g.parseCode("test.go", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)

	err = g.parseHttpAPIMarkers()
	assert.Nil(t, err)

	assert.True(t, len(g.markers) == 2)

	g.genHttpAPIDeclare()
	codeBuf := bytes.NewBuffer(make([]byte, 0))
	g.outputAPIDeclare(codeBuf)

	code, err := format.Source(codeBuf.Bytes())
	assert.Nil(t, err)
	t.Log(string(code))
}
