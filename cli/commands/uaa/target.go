package uaa

import (
	"strconv"

	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/urfave/cli"
)

var targetFlags = []cli.Flag{}

var TargetCommand = cli.Command{
	Name:      "target",
	ShortName: "t",
	Usage:     "Set or view the targeted Predix UAA instance and context",
	ArgsUsage: "[ID]",
	Flags:     targetFlags,
	Before: func(c *cli.Context) error {
		if c.NArg() > 1 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 1 {
			id, err := strconv.Atoi(c.Args()[0])
			if err != nil {
				global.UI.Failed("Invalid ID.")
			}
			uaac.Targets.SetCurrentForID(id)
		}

		uaac.Targets.PrintCurrent()
		return nil
	},
}
