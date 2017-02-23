package uaa

import (
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.build.ge.com/adoption/predix-cli/cli/commands/uaa/client"
	"github.build.ge.com/adoption/predix-cli/cli/commands/uaa/user"
	"github.com/urfave/cli"
)

var UaaCommand = cli.Command{
	Name:  "uaa",
	Usage: "Manage Predix UAA instance",
}

func init() {
	UaaCommand.Subcommands = []cli.Command{
		LoginCommand,
		TargetCommand,
		TargetsCommand,
		client.ClientsCommand,
		client.ClientCommand,
		user.UsersCommand,
		user.UserCommand,
	}
	commandregistry.AddCommand(UaaCommand)
}
