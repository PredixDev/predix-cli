package uaa

import (
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/urfave/cli"
)

var TargetsCommand = cli.Command{
	Name:  "targets",
	Usage: "Display all Predix UAA targets and contexts",
	Action: func(c *cli.Context) error {
		uaac.Targets.PrintAll()
		return nil
	},
}
