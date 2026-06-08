package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		baseURL: cfg.API.BaseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.API.Timeout) * time.Second,
		},
	}
}

// postJSON 发送 POST 请求并解析 JSON 响应
func (c *Client) postJSON(ctx context.Context, path string, req interface{}, resp interface{}) (interface{}, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	// 显式 close（不用 defer），并在错误分支前 drain body
	// 避免 HTTP/1.1 keep-alive 连接无法复用
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, httpResp.Body)
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp, nil
}

// getJSON 发送 GET 请求并解析 JSON 响应
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

	if httpResp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, httpResp.Body)
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
