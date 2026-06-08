package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tiamxu/leister/config"
)

// Client API 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建 API 客户端实例
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.API.BaseURL, "/"),
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.API.Timeout) * time.Second,
		},
	}
}

// apiEnvelope 统一响应外壳 {code, msg, data}
type apiEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data,omitempty"`
}

// postJSON 发送 POST 请求，解析 {code, msg, data} 响应
// 业务成功时把 data 反序列化到 resp；resp 可为 nil
func (c *Client) postJSON(ctx context.Context, path string, req interface{}, resp interface{}) error {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", httpResp.StatusCode, string(body))
	}

	var env apiEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if env.Code != 200 {
		return fmt.Errorf("api error [%d]: %s", env.Code, env.Msg)
	}

	if resp != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, resp); err != nil {
			return fmt.Errorf("failed to decode data: %w", err)
		}
	}
	return nil
}

// getJSON 发送 GET 请求，解析 {code, msg, data} 响应
func (c *Client) getJSON(ctx context.Context, path string, resp interface{}) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", httpResp.StatusCode, string(body))
	}

	var env apiEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if env.Code != 200 {
		return fmt.Errorf("api error [%d]: %s", env.Code, env.Msg)
	}

	if resp != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, resp); err != nil {
			return fmt.Errorf("failed to decode data: %w", err)
		}
	}
	return nil
}
