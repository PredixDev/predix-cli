package user

import (
	"github.com/urfave/cli"
)

var UserCommand = cli.Command{
	Name:      "user",
	ShortName: "u",
	Usage:     "Manage user accounts on the targeted Predix UAA instance",
}

func init() {
	UserCommand.Subcommands = []cli.Command{
		CreateUserCommand,
		GetUserCommand,
		UpdateUserCommand,
		DeleteUserCommand,
	}
}
