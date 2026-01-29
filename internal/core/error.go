/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/27 1:12 下午
 * @Desc: 错误处理实现
 */

package core

import "fmt"

type Error interface {
	error
	Code() int
	Message() string
}

type respError struct {
	code    int
	message string
}

func NewError(code int, message string) Error {
	return &respError{
		code:    code,
		message: message,
	}
}

func (e *respError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

func (e *respError) Code() int {
	return e.code
}

func (e *respError) Message() string {
	return e.message
}
