/**
 * @Author: fuxiao
 * @Email: 576101059@qq.com
 * @Date: 2021/8/27 11:31 上午
 * @Desc: HTTP 客户端核心实现
 */

package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/d60-Lab/tencent-im/internal/enum"
	"github.com/d60-Lab/tencent-im/internal/sign"
	"github.com/d60-Lab/tencent-im/internal/types"
)

const (
	defaultBaseUrl     = "https://adminapiger.im.qcloud.com"
	defaultBackupUrl   = "https://console.tim.qq.com"
	defaultVersion     = "v4"
	defaultContentType = "json"
	defaultExpiration  = 3600
	defaultTimeout     = 30 * time.Second
)

var invalidResponse = NewError(enum.InvalidResponseCode, "invalid response")

type Client interface {
	// Get GET请求
	Get(serviceName string, command string, data interface{}, resp interface{}) error
	// Post POST请求
	Post(serviceName string, command string, data interface{}, resp interface{}) error
	// Put PUT请求
	Put(serviceName string, command string, data interface{}, resp interface{}) error
	// Patch PATCH请求
	Patch(serviceName string, command string, data interface{}, resp interface{}) error
	// Delete DELETE请求
	Delete(serviceName string, command string, data interface{}, resp interface{}) error
}

type client struct {
	httpClient      *http.Client
	opt             *Options
	userSig         string
	userSigExpireAt int64
	baseUrl         string
	backupUrl       string
	logger          Logger
}

type Options struct {
	AppId      int           // 应用SDKAppID，可在即时通信 IM 控制台 的应用卡片中获取。
	AppSecret  string        // 密钥信息，可在即时通信 IM 控制台 的应用详情页面中获取，具体操作请参见 获取密钥
	UserId     string        // 用户ID
	Expiration int           // UserSig过期时间（秒）
	BaseUrl    string        // 可选：自定义 API 基础 URL
	BackupUrl  string        // 可选：备用 URL
	Timeout    time.Duration // 可选：请求超时时间，默认 30 秒
	Logger     Logger        // 可选：自定义日志实现，默认使用标准输出
	Debug      bool          // 可选：是否开启调试模式
}

func NewClient(opt *Options) Client {
	if opt.BaseUrl == "" {
		opt.BaseUrl = defaultBaseUrl
	}
	if opt.BackupUrl == "" {
		opt.BackupUrl = defaultBackupUrl
	}
	if opt.Timeout == 0 {
		opt.Timeout = defaultTimeout
	}
	if opt.Logger == nil {
		if opt.Debug {
			opt.Logger = NewDefaultLogger()
		} else {
			opt.Logger = NewNoopLogger()
		}
	}

	rand.Seed(time.Now().UnixNano())

	c := &client{
		opt:        opt,
		baseUrl:    opt.BaseUrl,
		backupUrl:  opt.BackupUrl,
		logger:     opt.Logger,
		httpClient: &http.Client{Timeout: opt.Timeout},
	}

	return c
}

// Get GET请求
func (c *client) Get(serviceName string, command string, data interface{}, resp interface{}) error {
	return c.request(http.MethodGet, serviceName, command, data, resp)
}

// Post POST请求
func (c *client) Post(serviceName string, command string, data interface{}, resp interface{}) error {
	return c.request(http.MethodPost, serviceName, command, data, resp)
}

// Put PUT请求
func (c *client) Put(serviceName string, command string, data interface{}, resp interface{}) error {
	return c.request(http.MethodPut, serviceName, command, data, resp)
}

// Patch PATCH请求
func (c *client) Patch(serviceName string, command string, data interface{}, resp interface{}) error {
	return c.request(http.MethodPatch, serviceName, command, data, resp)
}

// Delete DELETE请求
func (c *client) Delete(serviceName string, command string, data interface{}, resp interface{}) error {
	return c.request(http.MethodDelete, serviceName, command, data, resp)
}

// request Request请求
func (c *client) request(method, serviceName, command string, data, resp interface{}) error {
	ctx := context.Background()
	url := c.buildUrl(c.baseUrl, serviceName, command)

	// 序列化请求数据
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			c.logger.Error(ctx, "Failed to marshal request data", map[string]interface{}{
				"error": err,
			})
			return err
		}
		body = bytes.NewBuffer(jsonData)

		c.logger.Debug(ctx, "Request data", map[string]interface{}{
			"method": method,
			"url":    url,
			"body":   string(jsonData),
		})
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		c.logger.Error(ctx, "Failed to create request", map[string]interface{}{
			"error": err,
		})
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(ctx, "Failed to send request", map[string]interface{}{
			"error": err,
			"url":   url,
		})
		return err
	}
	defer httpResp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		c.logger.Error(ctx, "Failed to read response", map[string]interface{}{
			"error": err,
		})
		return err
	}

	c.logger.Debug(ctx, "Response received", map[string]interface{}{
		"status": httpResp.StatusCode,
		"body":   string(respBody),
	})

	// 解析响应
	if err = json.Unmarshal(respBody, resp); err != nil {
		c.logger.Error(ctx, "Failed to unmarshal response", map[string]interface{}{
			"error": err,
			"body":  string(respBody),
		})
		return err
	}

	// 检查业务错误
	if r, ok := resp.(types.ActionBaseRespInterface); ok {
		if r.GetActionStatus() == enum.FailActionStatus {
			return NewError(r.GetErrorCode(), r.GetErrorInfo())
		}

		if r.GetErrorCode() != enum.SuccessCode {
			return NewError(r.GetErrorCode(), r.GetErrorInfo())
		}
	} else if r, ok := resp.(types.BaseRespInterface); ok {
		if r.GetErrorCode() != enum.SuccessCode {
			return NewError(r.GetErrorCode(), r.GetErrorInfo())
		}
	} else {
		return invalidResponse
	}

	return nil
}

// buildUrl 构建一个请求URL
func (c *client) buildUrl(baseUrl, serviceName, command string) string {
	format := "%s/%s/%s/%s?sdkappid=%d&identifier=%s&usersig=%s&random=%d&contenttype=%s"
	random := rand.Int31()
	userSig := c.getUserSig()
	return fmt.Sprintf(format, baseUrl, defaultVersion, serviceName, command, c.opt.AppId, c.opt.UserId, userSig, random, defaultContentType)
}

// getUserSig 获取签名
func (c *client) getUserSig() string {
	now, expiration := time.Now(), c.opt.Expiration

	if expiration <= 0 {
		expiration = defaultExpiration
	}

	if c.userSig == "" || c.userSigExpireAt <= now.Unix() {
		c.userSig, _ = sign.GenUserSig(c.opt.AppId, c.opt.AppSecret, c.opt.UserId, expiration)
		c.userSigExpireAt = now.Add(time.Duration(expiration) * time.Second).Unix()
	}

	return c.userSig
}
