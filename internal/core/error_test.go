/**
 * @Author: d60-Lab
 * @Date: 2026/1/29
 * @Desc: 错误处理单元测试
 */

package core

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
		want    Error
	}{
		{
			name:    "success code",
			code:    0,
			message: "success",
			want:    NewError(0, "success"),
		},
		{
			name:    "error code",
			code:    10002,
			message: "invalid parameter",
			want:    NewError(10002, "invalid parameter"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.code, tt.message)
			if got.Code() != tt.code {
				t.Errorf("NewError().Code() = %v, want %v", got.Code(), tt.code)
			}
			if got.Message() != tt.message {
				t.Errorf("NewError().Message() = %v, want %v", got.Message(), tt.message)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	err := NewError(10002, "test error")
	expected := "code: 10002, message: test error"

	if err.Error() != expected {
		t.Errorf("Error.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestError_IsError(t *testing.T) {
	err := NewError(10002, "test error")

	// 测试 error 接口
	var e error = err
	if e == nil {
		t.Error("Expected error to implement error interface")
	}
}

func TestError_TypeAssertion(t *testing.T) {
	var err error = NewError(10002, "test error")

	// 测试类型断言
	if timErr, ok := err.(Error); ok {
		if timErr.Code() != 10002 {
			t.Errorf("Code() = %v, want 10002", timErr.Code())
		}
	} else {
		t.Error("Type assertion failed")
	}
}

func TestError_NotEqual(t *testing.T) {
	err1 := NewError(10002, "error1")
	err2 := NewError(10003, "error2")
	stdErr := errors.New("standard error")

	if err1.Error() == err2.Error() {
		t.Error("Different errors should not be equal")
	}

	if err1.Error() == stdErr.Error() {
		t.Error("TIM error should not equal standard error")
	}
}
