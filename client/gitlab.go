package client

import (
	"context"

	"github.com/tiamxu/leister/types"
)

// GetGitlabProject 获取 GitLab 项目信息
func (c *Client) GetGitlabProject(ctx context.Context, req *types.GitlabProjectRequest) (*types.GitlabProjectResponse, error) {
	resp := &types.GitlabProjectResponse{}
	if err := c.postJSON(ctx, "/api/gitlab/project", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GenGitlabProjects 生成 GitLab 项目数据
func (c *Client) GenGitlabProjects(ctx context.Context, req *types.GitlabGenRequest) (*types.GitlabGenResponse, error) {
	resp := &types.GitlabGenResponse{}
	if err := c.postJSON(ctx, "/api/gitlab/gen", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
