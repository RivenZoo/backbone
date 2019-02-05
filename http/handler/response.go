package handler

var defaultResponseEncoder = ResponseBodyEncodeFunc(func(v interface{}) ([]byte, error) {
	return json.Marshal(v)
})

// SetDefaultResponseBodyEncoder set default ResponseBodyEncodeFunc.
// If you use NewRequestHandleFunc and RequestProcessor.ResponseEncoder not set, this encoder will be used.
func SetDefaultResponseBodyEncoder(encoder ResponseBodyEncodeFunc) {
	defaultResponseEncoder = encoder
}
