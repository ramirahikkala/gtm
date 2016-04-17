package command

import (
	"flag"
	"fmt"
	"os"

	"edgeg.io/gtm/metric"
	"github.com/mitchellh/cli"
)

type GitCommit struct {
}

func NewCommit() (cli.Command, error) {
	return GitCommit{}, nil
}

func (r GitCommit) Help() string {
	return `
	Log time for git tracked files and set the file's tracked time to zero.

	gtm commit [--dry-run] [--debug]
	`
}

func (r GitCommit) Run(args []string) int {
	commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
	dryRun := commitFlags.Bool(
		"dry-run",
		true,
		"Do not log time but show time logged for all files")
	debug := commitFlags.Bool(
		"debug",
		false,
		"Print debug statements to the console")
	commitFlags.Parse(os.Args[2:])

	m, err := metric.Process(*dryRun, *debug)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Println(m)
	return 0
}

func (r GitCommit) Synopsis() string {
	return `
	Log time for git tracked files
	`
}
