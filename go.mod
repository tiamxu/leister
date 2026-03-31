module github.com/tiamxu/leister

go 1.23.4

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/tiamxu/kit v0.0.0-20230519061052-8dd21a1d8a48
	github.com/xanzy/go-gitlab v0.83.0
)

replace github.com/tiamxu/kit => ../kit

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
