package client

import (
	"context"

	"github.com/tiamxu/leister/types"
)

// GetGitlabProject 获取 GitLab 项目信息
func (c *Client) GetGitlabProject(ctx context.Context, req *types.GitlabProjectRequest) (*types.GitlabProjectResponse, error) {
	resp, err := c.postJSON(ctx, "/api/gitlab/project", req, &types.GitlabProjectResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*types.GitlabProjectResponse), nil
}

// GenGitlabProjects 生成 GitLab 项目数据
func (c *Client) GenGitlabProjects(ctx context.Context, req *types.GitlabGenRequest) (*types.GitlabGenResponse, error) {
	resp, err := c.postJSON(ctx, "/api/gitlab/gen", req, &types.GitlabGenResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*types.GitlabGenResponse), nil
}
