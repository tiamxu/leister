package gitlab

import (
	"github.com/urfave/cli/v2"
)

var (
	GetProjectCmd = cli.Command{
		Name:   "get",
		Usage:  "get gitlab project info console",
		Flags:  Flags,
		Before: InitFlags,
		Action: RunGetProject,
	}
	GenProjectDBCmd = cli.Command{
		Name:   "gen",
		Usage:  "generate gitlab project data to db",
		Flags:  Flags,
		Before: InitGenFlags,
		Action: RunGenProject,
	}
)
