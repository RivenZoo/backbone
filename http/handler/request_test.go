package handler

import (
	"bytes"
	"errors"
	"github.com/RivenZoo/backbone/http/handler/error_code"
	"github.com/RivenZoo/backbone/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type testReq struct {
	Ping string `json:"ping"`
}

type testResp struct {
	Pong string `json:"pong"`
}

func newTestReq() interface{} {
	return &testReq{}
}

func processTestReq(c *gin.Context, req interface{}) (resp interface{}, err error) {
	tr := req.(*testReq)
	return &testResp{tr.Ping}, nil
}

func postProcessTestReq(c *gin.Context, resp interface{}, err error) () {
	logger.Logf("[INFO] resp %v, error %v", resp, err)
}

func TestHandleRequest(t *testing.T) {
	key := "body-key"
	p := &RequestProcessor{
		NewReqFunc:  newTestReq,
		ProcessFunc: processTestReq,
	}
	req := &testReq{"test"}
	expectResp := &testResp{"test"}

	data, err := json.Marshal(req)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(data))

	resp, err := handleRequest(c, p)
	assert.Nil(t, err)
	assert.EqualValues(t, expectResp, resp)
	t.Log(resp, err)

	c, _ = gin.CreateTestContext(w)

	c.Set(key, data)
	p.BodyContextKey = key

	resp, err = handleRequest(c, p)
	assert.Nil(t, err)
	assert.EqualValues(t, expectResp, resp)
	t.Log(resp, err)
}

func TestHandleResponse(t *testing.T) {
	p := &RequestProcessor{
		NewReqFunc:  newTestReq,
		ProcessFunc: processTestReq,
	}
	req := &testReq{"test"}
	expectResp := &testResp{"test"}

	data, err := json.Marshal(req)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(data))

	resp, err := handleRequest(c, p)
	assert.Nil(t, err)
	assert.EqualValues(t, expectResp, resp)
	t.Log(resp, err)

	handleResponse(c, resp, p.ResponseEncoder)
	retData, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)
	expect, err := defaultResponseEncoder(expectResp)
	assert.Nil(t, err)
	assert.EqualValues(t, expect, retData)

	t.Log(string(retData))
}

func TestHandleError(t *testing.T) {
	p := &RequestProcessor{
		NewReqFunc: newTestReq,
		ProcessFunc: func(c *gin.Context, req interface{}) (resp interface{}, err error) {
			logger.Logf("req %v", req)
			return nil, errors.New("server error")
		},
	}
	req := &testReq{"test"}

	data, err := json.Marshal(req)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(data))

	resp, err := handleRequest(c, p)
	assert.NotNil(t, err)
	t.Log(resp, err)

	handleError(c, err, p.ErrorEncoder)

	retData, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)
	expect, err := defaultErrorEncoder(c, error_code.ErrServerError)
	assert.Nil(t, err)
	assert.EqualValues(t, expect, retData)

	t.Log(string(retData))
}

func TestRequestPostProcess(t *testing.T) {
	p := &RequestProcessor{
		NewReqFunc:      newTestReq,
		ProcessFunc:     processTestReq,
		PostProcessFunc: postProcessTestReq,
	}
	req := &testReq{"test"}
	expectResp := &testResp{"test"}

	data, err := json.Marshal(req)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(data))

	resp, err := handleRequest(c, p)
	assert.Nil(t, err)
	assert.EqualValues(t, expectResp, resp)
	t.Log(resp, err)

	handlePostRequest(c, resp, err, p.PostProcessFunc)
}
