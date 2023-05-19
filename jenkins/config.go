package jenkins

import (
	"github.com/koding/multiconfig"
	"github.com/tiamxu/leister/config"
)

const configPath = "config/config.yaml"

var (
	cfg *config.Config
)

func loadConfig() {
	cfg = new(config.Config)
	multiconfig.MustLoadWithPath(configPath, cfg)
}
