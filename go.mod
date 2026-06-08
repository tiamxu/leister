module github.com/tiamxu/leister

go 1.25.0

require (
	github.com/spf13/cobra v1.10.2
	github.com/tiamxu/kit v0.0.0-20230519061052-8dd21a1d8a48
)

replace github.com/tiamxu/kit => ../kit

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
)
