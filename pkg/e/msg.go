package e

// GetMsg 根据错误码返回中文错误信息
func GetMsg(code int) string {
	switch code {
	case ErrSuccess:
		return "成功"
	case ErrParamInvalid:
		return "参数无效"
	case ErrInternal:
		return "内部错误"

	case ErrJenkinsCreate:
		return "创建Jenkins任务失败"
	case ErrJenkinsUpdate:
		return "更新Jenkins任务失败"
	case ErrJenkinsList:
		return "获取Jenkins项目列表失败"

	case ErrGitlabGet:
		return "获取GitLab项目信息失败"
	case ErrGitlabGen:
		return "生成GitLab项目数据失败"

	case ErrDockerBuild:
		return "Docker镜像构建失败"
	case ErrDockerPush:
		return "Docker镜像推送失败"
	case ErrDockerLogin:
		return "Docker登录失败"
	case ErrDockerCredential:
		return "Docker凭证未配置"

	case ErrKubeGet:
		return "获取K8s部署失败"
	case ErrKubeCreate:
		return "创建K8s部署失败"
	case ErrKubeRestart:
		return "重启K8s部署失败"

	default:
		return "未知错误"
	}
}
