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
	"github.com/RivenZoo/sqlagent"
)

// @HttpAPI("/url", apiReq, apiResp)

/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/
`
	g := newHttpAPIGenerator(httpAPIGeneratorOption{
		imports: []string{"github.com/gin"},
	})
	err := g.parseCode("test.go", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)

	err = g.parseHttpAPIMarkers()
	assert.Nil(t, err)

	assert.True(t, len(g.markers) == 2)

	g.genHttpAPIDeclare()


}
