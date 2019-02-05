package error_code

var OK = NewErrorCode(0, "ok")

// error code format: [4|5][000][00][00]
// client error
var (
	ErrBadRequest = NewErrorCode(40000001, "bad request")
)

// server error
var (
	ErrServerError = NewErrorCode(50000001, "server error")
)
