package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpAPIGenerator_GenHttpAPIDeclare(t *testing.T) {
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
		commonAPIDefinition: commonHttpAPIDefinitionOption{
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

func TestHttpAPIGenerator_GenHttpAPIHandler(t *testing.T) {
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
		apiHandlerFileImports: []importInfo{
			{"github.com/RivenZoo/backbone/http/handler/error_code", ""},
		},
		commonHttpAPIHandler: commonHttpAPIHandlerOption{
			ErrorEncoder: `func(err error) ([]byte, error) {
	e, ok := err.(error_code.ErrorCode)
	if !ok {
		s := fmt.Sprintf("{\"resp_common\": %s}", error_code.ErrServerError.Error())
		return []byte(s), nil
	}
	s := fmt.Sprintf("{\"resp_common\": %s}", e.Error())
	return []byte(s), nil
}`,
		},
	})
	err := g.ParseCode("test.go", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)

	err = g.ParseHttpAPIMarkers()
	assert.Nil(t, err)

	assert.True(t, len(g.markers) == 2)

	g.GenHttpAPIHandler()
	codeBuf := bytes.NewBuffer(make([]byte, 0))
	err = g.OutputAPIHandler(codeBuf)
	assert.Nil(t, err)

	t.Log(codeBuf.String())
}
