package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// JenkinsJobRequest Jenkins 任务创建请求
type JenkinsJobRequest struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

// JenkinsJobResponse Jenkins 任务创建响应
type JenkinsJobResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GitlabProjectRequest GitLab 项目获取请求
type GitlabProjectRequest struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

// GitlabProjectResponse GitLab 项目获取响应
type GitlabProjectResponse struct {
	Status  string       `json:"status"`
	Project *ProjectInfo `json:"project"`
}

// ProjectInfo 项目信息
type ProjectInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Group          string `json:"group"`
	HTTPURLToRepo  string `json:"http_url_to_repo"`
	SSHURLToRepo   string `json:"ssh_url_to_repo"`
}

// GitlabGenRequest GitLab 项目生成请求
type GitlabGenRequest struct {
	Group string `json:"group"`
}

// GitlabGenResponse GitLab 项目生成响应
type GitlabGenResponse struct {
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Projects []*ProjectInfo `json:"projects"`
}

// CreateJenkinsJob 创建 Jenkins 任务
func (c *Client) CreateJenkinsJob(ctx context.Context, req *JenkinsJobRequest) (*JenkinsJobResponse, error) {
	resp, err := c.postJSON(ctx, "/api/jenkins/create", req, &JenkinsJobResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*JenkinsJobResponse), nil
}

// CreateJenkinsJobs 批量创建 Jenkins 任务
func (c *Client) CreateJenkinsJobs(ctx context.Context, reqs []*JenkinsJobRequest) (*JenkinsJobResponse, error) {
	resp, err := c.postJSON(ctx, "/api/jenkins/cts", reqs, &JenkinsJobResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*JenkinsJobResponse), nil
}

// UpdateJenkinsJob 更新 Jenkins 任务
func (c *Client) UpdateJenkinsJob(ctx context.Context, req *JenkinsJobRequest) (*JenkinsJobResponse, error) {
	resp, err := c.postJSON(ctx, "/api/jenkins/update", req, &JenkinsJobResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*JenkinsJobResponse), nil
}

// GetGitlabProject 获取 GitLab 项目信息
func (c *Client) GetGitlabProject(ctx context.Context, req *GitlabProjectRequest) (*GitlabProjectResponse, error) {
	resp, err := c.postJSON(ctx, "/api/gitlab/project", req, &GitlabProjectResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*GitlabProjectResponse), nil
}

// GenGitlabProjects 生成 GitLab 项目数据
func (c *Client) GenGitlabProjects(ctx context.Context, req *GitlabGenRequest) (*GitlabGenResponse, error) {
	resp, err := c.postJSON(ctx, "/api/gitlab/gen", req, &GitlabGenResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*GitlabGenResponse), nil
}

// postJSON 发送 POST 请求并解析 JSON 响应
func (c *Client) postJSON(ctx context.Context, path string, req interface{}, resp interface{}) (interface{}, error) {
	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer httpResp.Body.Close()

	// 检查响应状态码
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	// 解析响应体
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return resp, nil
}
