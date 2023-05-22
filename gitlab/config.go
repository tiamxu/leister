package gitlab

import (
	"github.com/koding/multiconfig"
	"github.com/tiamxu/leister/database"
	"github.com/xanzy/go-gitlab"
)

const configPath = "config/config.yaml"

var (
	cfg *Config
)

type Config struct {
	Jenkins `yaml:"jenkins"`
	Gitlab  `yaml:"gitlab"`
	DB      *database.Config `yaml:"db"`
}
type Jenkins struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Gitlab struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

func loadConfig() {
	cfg = new(Config)
	multiconfig.MustLoadWithPath(configPath, cfg)
}

func Connect(cfg *Config) (*gitlab.Client, error) {
	git, err := gitlab.NewClient(cfg.Gitlab.Token, gitlab.WithBaseURL(cfg.Gitlab.Url))
	if err != nil {
		return git, err
	}
	return git, nil
}
