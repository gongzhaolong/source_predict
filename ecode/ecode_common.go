package ecode

import "fmt"

type ErrCode struct {
	code int    // 错误码
	msg  string // 错误信息
	err  error  // debug 错误信息
}

func (err *ErrCode) Error() string {
	str := fmt.Sprintf("%s：%s", err.msg, err.err)
	return str
}

func (err *ErrCode) Code() int {
	return err.code
}

func (err *ErrCode) String() string {
	return err.msg
}

func (err *ErrCode) Err() error {
	return err.err
}
