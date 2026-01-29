/**
 * @Author: d60-Lab
 * @Date: 2026/1/29
 * @Desc: 核心客户端单元测试
 */

package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		opt  *Options
		want bool
	}{
		{
			name: "valid options with defaults",
			opt: &Options{
				AppId:     1400000000,
				AppSecret: "test-secret",
				UserId:    "admin",
			},
			want: true,
		},
		{
			name: "valid options with custom base url",
			opt: &Options{
				AppId:     1400000000,
				AppSecret: "test-secret",
				UserId:    "admin",
				BaseUrl:   "https://custom.example.com",
			},
			want: true,
		},
		{
			name: "valid options with debug mode",
			opt: &Options{
				AppId:     1400000000,
				AppSecret: "test-secret",
				UserId:    "admin",
				Debug:     true,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.opt)
			if (got != nil) != tt.want {
				t.Errorf("NewClient() = %v, want %v", got != nil, tt.want)
			}
		})
	}
}

func TestClient_Request(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ActionStatus":"OK","ErrorCode":0,"ErrorInfo":""}`))
	}))
	defer ts.Close()

	client := NewClient(&Options{
		AppId:     1400000000,
		AppSecret: "test-secret",
		UserId:    "admin",
		BaseUrl:   ts.URL,
		Timeout:   5 * time.Second,
	})

	// 简单测试 - 确保客户端可以正常工作
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestOptions_Defaults(t *testing.T) {
	opt := &Options{
		AppId:     1400000000,
		AppSecret: "test-secret",
		UserId:    "admin",
	}

	client := NewClient(opt).(*client)

	// 验证默认值
	if client.baseUrl != defaultBaseUrl {
		t.Errorf("Expected baseUrl = %s, got %s", defaultBaseUrl, client.baseUrl)
	}

	if client.backupUrl != defaultBackupUrl {
		t.Errorf("Expected backupUrl = %s, got %s", defaultBackupUrl, client.backupUrl)
	}

	if client.logger == nil {
		t.Error("Expected logger to be set, got nil")
	}
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name   string
		logger Logger
	}{
		{
			name:   "default logger",
			logger: NewDefaultLogger(),
		},
		{
			name:   "noop logger",
			logger: NewNoopLogger(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.logger == nil {
				t.Error("Logger should not be nil")
			}

			// 测试所有日志方法不会panic
			tt.logger.Debug(nil, "test", nil)
			tt.logger.Info(nil, "test", nil)
			tt.logger.Warn(nil, "test", nil)
			tt.logger.Error(nil, "test", nil)
		})
	}
}
