package cache

import (
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.com/urfave/cli"
)

var CacheCommand = cli.Command{
	Name:  "cache",
	Usage: "Manage Predix CLI cache",
}

func init() {
	CacheCommand.Subcommands = []cli.Command{
		ConfigCommand(),
		ClearCommand(),
		UpdateCommand(),
	}
	commandregistry.AddCommand(CacheCommand)
}
