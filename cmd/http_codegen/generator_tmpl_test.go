package main

import (
	"bytes"
	"testing"
)

func TestGenHttpAPIDefinitionByTmpl(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	genHttpAPIDefinitionByTmpl(&HttpAPIMarker{
		FileScopeMarker: FileScopeMarker{
			Identity{Name: "test"},
		},
		RequestType:  "testReq",
		ResponseType: "testResp",
	}, buf, commonHttpAPIDefinition{})
	t.Log(buf.String())
}
