package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHttpAPIMarkers(t *testing.T) {
	src := `
package controller

// @HttpAPI("/url", apiReq, apiResp)

/*
* @HttpAPI("/url2", apiReq2, apiResp2)
*/
`
	sa, err := ParseSourceCode("test", bytes.NewReader([]byte(src)))
	assert.Nil(t, err)
	t.Log(sa.node)

	markers, err := ParseHttpAPIMarkers(sa)
	assert.Nil(t, err)
	assert.True(t, len(markers) == 2)

	assert.Equal(t, `"/url"`, markers[0].URL)
	assert.Equal(t, `apiReq`, markers[0].RequestType)
	assert.Equal(t, `apiResp`, markers[0].ResponseType)
}
