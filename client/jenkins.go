package client

import (
	"context"
	"net/url"

	"github.com/tiamxu/leister/types"
)

// CreateJenkinsJob 创建 Jenkins 任务
func (c *Client) CreateJenkinsJob(ctx context.Context, req *types.JenkinsJobRequest) (*types.JenkinsJobResponse, error) {
	resp, err := c.postJSON(ctx, "/api/jenkins/create", req, &types.JenkinsJobResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*types.JenkinsJobResponse), nil
}

// CreateJenkinsJobs 批量创建 Jenkins 任务
func (c *Client) CreateJenkinsJobs(ctx context.Context, reqs []*types.JenkinsJobRequest) (*types.JenkinsJobResponse, error) {
	resp, err := c.postJSON(ctx, "/api/jenkins/cts", reqs, &types.JenkinsJobResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*types.JenkinsJobResponse), nil
}

// UpdateJenkinsJob 更新 Jenkins 任务
func (c *Client) UpdateJenkinsJob(ctx context.Context, req *types.JenkinsJobRequest) (*types.JenkinsJobResponse, error) {
	resp, err := c.postJSON(ctx, "/api/jenkins/update", req, &types.JenkinsJobResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*types.JenkinsJobResponse), nil
}

// ListJenkinsProjects 列出 Jenkins 项目（用于批量创建时取真实列表）
func (c *Client) ListJenkinsProjects(ctx context.Context, group string) (*types.JenkinsProjectListResponse, error) {
	path := "/api/jenkins/projects"
	if group != "" {
		path += "?group=" + url.QueryEscape(group)
	}

	resp := &types.JenkinsProjectListResponse{}
	if err := c.getJSON(ctx, path, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
