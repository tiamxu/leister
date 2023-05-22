package jenkins

import (
	"github.com/urfave/cli/v2"
)

var (
	CreateJobCmd = cli.Command{
		Name:   "create",
		Usage:  "create one jenkins job",
		Flags:  Flags,
		Before: InitFlags,
		Action: RunCreateJob,
	}
	CreateJobsCmd = cli.Command{
		Name:   "cts",
		Usage:  "create many jenkins jobs",
		Flags:  Flags,
		Before: InitCtxFlags,
		Action: RunCreateJobs,
	}
)
