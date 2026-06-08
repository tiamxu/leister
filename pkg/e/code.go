package e

const (
	// 通用
	ErrSuccess      = 0
	ErrParamInvalid = 9001
	ErrInternal     = 9999

	// Jenkins
	ErrJenkinsCreate = 1001
	ErrJenkinsUpdate = 1002
	ErrJenkinsList   = 1003

	// GitLab
	ErrGitlabGet = 2001
	ErrGitlabGen = 2002

	// Docker
	ErrDockerBuild     = 3001
	ErrDockerPush      = 3002
	ErrDockerLogin     = 3003
	ErrDockerCredential = 3004

	// Kube
	ErrKubeGet    = 4001
	ErrKubeCreate = 4002
	ErrKubeRestart = 4003
)
