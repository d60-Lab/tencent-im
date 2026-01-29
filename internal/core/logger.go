/**
 * @Author: d60-Lab
 * @Date: 2026/1/29
 * @Desc: 日志接口定义
 */

package core

import (
	"context"
	"log"
	"os"
)

// Logger 日志接口，用户可以实现此接口来自定义日志输出
type Logger interface {
	// Debug 输出调试日志
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	// Info 输出信息日志
	Info(ctx context.Context, msg string, fields map[string]interface{})
	// Warn 输出警告日志
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	// Error 输出错误日志
	Error(ctx context.Context, msg string, fields map[string]interface{})
}

// defaultLogger 默认日志实现（使用标准库 log）
type defaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger 创建默认日志实现
func NewDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "[TencentIM] ", log.LstdFlags),
	}
}

func (l *defaultLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *defaultLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Printf("[INFO] %s %v", msg, fields)
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Printf("[WARN] %s %v", msg, fields)
}

func (l *defaultLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Printf("[ERROR] %s %v", msg, fields)
}

// noopLogger 空日志实现（不输出任何日志）
type noopLogger struct{}

// NewNoopLogger 创建空日志实现
func NewNoopLogger() Logger {
	return &noopLogger{}
}

func (l *noopLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {}
func (l *noopLogger) Info(ctx context.Context, msg string, fields map[string]interface{})  {}
func (l *noopLogger) Warn(ctx context.Context, msg string, fields map[string]interface{})  {}
func (l *noopLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {}
