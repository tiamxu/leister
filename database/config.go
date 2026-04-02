package database

type Item struct {
	CodeID        int    `json:"code_id"`
	AppName       string `json:"app_name"`
	AppGroup      string `json:"app_group"`
	AppType       string `json:"app_type"`
	SSHURLToRepo  string `json:"ssh_url_to_repo"`
	HTTPURLToRepo string `json:"http_url_to_repo"`
}
