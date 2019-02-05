package handler

var defaultRequestBodyDecodeFunc = RequestBodyDecodeFunc(func(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
})

// SetRequestBodyDecoder set default RequestBodyDecodeFunc.
// If you use NewRequestHandleFunc and RequestProcessor.RequestDecoder not set, this decoder will be used.
func SetRequestBodyDecoder(decoder RequestBodyDecodeFunc) {
	defaultRequestBodyDecodeFunc = decoder
}
