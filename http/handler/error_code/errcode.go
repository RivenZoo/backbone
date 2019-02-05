package error_code

import "fmt"

type ErrorCode struct {
	ret int
	msg string
}

func NewErrorCode(ret int, msg string) ErrorCode {
	return ErrorCode{ret: ret, msg: msg}
}

func (c ErrorCode) OK() bool {
	return c.ret == 0
}

func (c ErrorCode) Error() string {
	return fmt.Sprintf(`{"ret": %d, "msg": "%s"}`, c.ret, c.msg)
}

func (c ErrorCode) ErrorInfo() (ret int, msg string) {
	return c.ret, c.msg
}
