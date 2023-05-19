package gitlab

import (
	"github.com/koding/multiconfig"
	"github.com/tiamxu/leister/database"
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
