package handler

import "github.com/RivenZoo/backbone/http/handler/error_code"

var defaultErrorEncoder = ErrorResponseEncodeFunc(func(err error) ([]byte, error) {
	if err == nil {
		return []byte(error_code.OK.Error()), nil
	}
	switch code := err.(type) {
	case error_code.ErrorCode:
		return []byte(code.Error()), nil
	default:
	}
	return []byte(error_code.ErrServerError.Error()), nil
})

// SetDefaultErrorEncoder set default ErrorResponseEncodeFunc.
// If you use NewRequestHandleFunc and RequestProcessor.ErrorEncoder not set, this encoder will be used.
func SetDefaultErrorEncoder(encoder ErrorResponseEncodeFunc) {
	defaultErrorEncoder = encoder
}
