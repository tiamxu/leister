package types

// GitlabProjectRequest 获取 GitLab 单个项目请求
type GitlabProjectRequest struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

// GitlabProjectResponse GitLab 项目响应
type GitlabProjectResponse struct {
	Status  string       `json:"status"`
	Project *ProjectInfo `json:"project"`
}

// ProjectInfo 项目信息
type ProjectInfo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Group         string `json:"group"`
	HTTPURLToRepo string `json:"http_url_to_repo"`
	SSHURLToRepo  string `json:"ssh_url_to_repo"`
}

// GitlabGenRequest 生成 GitLab 组项目数据请求
type GitlabGenRequest struct {
	Group string `json:"group"`
}

// GitlabGenResponse 生成结果响应
type GitlabGenResponse struct {
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Projects []*ProjectInfo `json:"projects"`
}
