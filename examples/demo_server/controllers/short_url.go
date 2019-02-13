package controllers

import (
	"github.com/RivenZoo/backbone/examples/demo_server/model"
	"github.com/RivenZoo/backbone/http/handler/error_code"
	"github.com/gin-gonic/gin"
	"hash/adler32"
	"strconv"
)

import (
	"fmt"
)

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
	return nil, nil
}

func formatURL(key string) string {
	return fmt.Sprintf("http://example.com/abbr/%s", key)
}
