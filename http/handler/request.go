package handler

import (
	"github.com/RivenZoo/backbone/logger"
	"github.com/gin-gonic/gin"
)

type NewRequestFunc func() interface{}

type RequestBodyDecodeFunc func(data []byte, v interface{}) error

type ResponseBodyEncodeFunc func(v interface{}) ([]byte, error)

type ErrorResponseEncodeFunc func(c *gin.Context, err error) ([]byte, error)

type RequestProcessFunc func(c *gin.Context, req interface{}) (resp interface{}, err error)

type RequestPostProcessFunc func(c *gin.Context, resp interface{}, err error)

type RequestProcessor struct {
	NewReqFunc  NewRequestFunc
	ProcessFunc RequestProcessFunc

	// optional
	// if BodyContextKey is set, get request body from gin.Context.
	BodyContextKey string
	// RequestDecoder if not set, use default json RequestBodyDecodeFunc
	RequestDecoder RequestBodyDecodeFunc
	// ResponseEncoder if not set, use default json ResponseBodyEncodeFunc
	ResponseEncoder ResponseBodyEncodeFunc
	// ResponseEncoder if not set, use default ErrorResponseEncodeFunc
	ErrorEncoder ErrorResponseEncodeFunc
	// PostProcessFunc if not set, skip post process.
	PostProcessFunc RequestPostProcessFunc
}

func NewRequestHandleFunc(p *RequestProcessor) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := handleRequest(c, p)
		if err != nil {
			handlePostRequest(c, resp, err, p.PostProcessFunc)
			handleError(c, err, p.ErrorEncoder)
			return
		}
		handlePostRequest(c, resp, err, p.PostProcessFunc)
		handleResponse(c, resp, p.ResponseEncoder)
	}
}

func handleRequest(c *gin.Context, p *RequestProcessor) (resp interface{}, err error) {
	data, err := getRequestBody(c, p.BodyContextKey)
	if err != nil {
		logger.Logf("[ERROR] getRequestBody error %v", err)
		return nil, err
	}
	req, err := decodeRequest(data, p)
	if err != nil {
		logger.Logf("[ERROR] decodeRequest error %v", err)
		return nil, err
	}
	resp, err = p.ProcessFunc(c, req)
	if err != nil {
		return nil, err
	}
	return
}

func handlePostRequest(c *gin.Context, resp interface{}, err error, postHandler RequestPostProcessFunc) {
	if postHandler != nil {
		postHandler(c, resp, err)
	}
}

func decodeRequest(data []byte, p *RequestProcessor) (req interface{}, err error) {
	req = p.NewReqFunc()
	decoder := p.RequestDecoder
	if decoder == nil {
		decoder = defaultRequestBodyDecodeFunc
	}
	err = decoder(data, req)
	return req, err
}

func getRequestBody(c *gin.Context, bodyKey string) ([]byte, error) {
	if bodyKey != "" {
		v, exists := c.Get(bodyKey)
		if !exists {
			return c.GetRawData()
		}
		data, ok := v.([]byte)
		if !ok {
			return c.GetRawData()
		}
		return data, nil
	}
	return c.GetRawData()
}

func handleError(c *gin.Context, err error, encoder ErrorResponseEncodeFunc) {
	if encoder == nil {
		encoder = defaultErrorEncoder
	}
	data, err := encoder(c, err)
	if err != nil {
		logger.Logf("[ERROR] EncodeError error %v", err)
		return
	}
	c.Writer.Write(data)
}

func handleResponse(c *gin.Context, resp interface{}, encoder ResponseBodyEncodeFunc) {
	if encoder == nil {
		encoder = defaultResponseEncoder
	}
	data, err := encoder(resp)
	if err != nil {
		logger.Logf("[ERROR] Encode error %v", err)
		return
	}
	c.Writer.Write(data)
}
