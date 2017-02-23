package cache

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.com/urfave/cli"
)

var clearCacheFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all",
		Usage: "Clear the Predix CLI cache for all targets",
	},
}

var clearCacheCommand = cli.Command{
	Name:  "clear",
	Usage: "Clear the Predix CLI cache for the current CF target or all targets",
	Flags: clearCacheFlags,
	Action: func(c *cli.Context) error {
		if c.Bool("all") {
			cf.Cache.PurgeAll()
		} else {
			cf.Cache.PurgeCurrent()
		}
		return nil
	},
	BashComplete: func(c *cli.Context) {
		helpers.Completions.PrintFlags(c, clearCacheFlags)
	},
}

func ClearCommand() cli.Command {
	return clearCacheCommand
}
