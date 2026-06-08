package types

// JenkinsJobRequest 创建/更新 Jenkins 任务请求
type JenkinsJobRequest struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

// JenkinsJobResponse Jenkins 任务响应（创建/更新/批量）
type JenkinsJobResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// JenkinsProjectItem Jenkins 项目列表项（用于批量创建时取真实项目）
type JenkinsProjectItem struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

// JenkinsProjectListResponse 项目列表响应
type JenkinsProjectListResponse struct {
	Status   string                `json:"status"`
	Projects []*JenkinsProjectItem `json:"projects"`
}
