package client

import (
	"github.com/urfave/cli"
)

var ClientCommand = cli.Command{
	Name:      "client",
	ShortName: "cl",
	Usage:     "Manage clients registered with the targeted Predix UAA instance",
}

func init() {
	ClientCommand.Subcommands = []cli.Command{
		CreateClientCommand,
		GetClientCommand,
		UpdateClientCommand,
		DeleteClientCommand,
	}
}
