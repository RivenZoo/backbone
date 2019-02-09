package controllers

import (
	"fmt"
	"github.com/RivenZoo/backbone/examples/demo_server/model"
	"github.com/RivenZoo/backbone/http/handler"
	"github.com/RivenZoo/backbone/http/handler/error_code"
	"github.com/gin-gonic/gin"
	"hash/adler32"
	"strconv"
)

var abbreviateURLProcessor = handler.NewRequestHandleFunc(&handler.RequestProcessor{
	NewReqFunc: func() interface{} {
		return &abbreviateURLReq{}
	},
	ProcessFunc: func(c *gin.Context, req interface{}) (resp interface{}, err error) {
		abbrReq := req.(*abbreviateURLReq)
		return handleAbbreviateURLReq(c, abbrReq)
	},
})

//go:generate http_codegen -input $GOFILE

// @HttpAPI("/url/abbr", abbreviateURLReq, abbreviateURLResp)
type abbreviateURLReq struct {
	URL string `json:"url"`
}

type abbreviateURLResp struct {
	URL string `json:"url"`
}

func handleAbbreviateURLReq(c *gin.Context, req *abbreviateURLReq) (resp *abbreviateURLResp, err error) {
	if req.URL == "" {
		return nil, error_code.ErrBadRequest
	}
	cs := adler32.Checksum([]byte(req.URL))
	s := strconv.FormatInt(int64(cs), 36)
	err = model.SetAbbreviateURL(s, req.URL)
	if err != nil {
		return nil, err
	}
	return &abbreviateURLResp{URL: formatURL(s)}, nil
}

func formatURL(key string) string {
	return fmt.Sprintf("http://example.com/abbr/%s", key)
}
