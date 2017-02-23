package client

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var createClientFlags = getCreateClientFlags()
var CreateClientCommand = cli.Command{
	Name:      "create",
	ShortName: "c",
	Usage:     "Register a client with the targeted Predix UAA instance",
	ArgsUsage: "CLIENT_ID",
	Flags:     createClientFlags,
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		var err error
		target, context, instance := uaac.Targets.GetCurrent()
		scim := lib.ScimFactory.New(target, context)

		client := lib.Client{}
		client.ID = c.Args()[0]
		client.Secret = clientCommandFlagDestinations.Secret
		if client.Secret == "" {
			client.Secret = global.UI.AskForVerifiedPassword("Client Secret")
		}
		client.Name = clientCommandFlagDestinations.Name
		client.Scopes = split(clientCommandFlagDestinations.Scopes)
		client.GrantTypes = split(clientCommandFlagDestinations.Grants)
		client.Authorities = split(clientCommandFlagDestinations.Authorities)
		client.AccessTokenTimeout = atoi(clientCommandFlagDestinations.AccessTokenTimeout)
		client.RefreshTokenTimeout = atoi(clientCommandFlagDestinations.RefreshTokenTimeout)
		client.RedirectURI = split(clientCommandFlagDestinations.RedirectURI)
		client.AutoApprove = split(clientCommandFlagDestinations.AutoApprove)
		client.SignupRedirect = clientCommandFlagDestinations.SignupRedirect

		global.UI.Say("Registering client %s on service instance %s", terminal.EntityNameColor(client.ID),
			terminal.EntityNameColor(instance.Name))
		err = scim.CreateClient(&client)
		if err != nil {
			global.UI.Failed(err.Error())
		}
		global.UI.Ok()

		return nil
	},
	BashComplete: func(c *cli.Context) {
		if c.NArg() == 0 {
			global.UI.Say("Enter: CLIENT_ID")
			global.UI.Say("_")
		} else {
			helpers.Completions.PrintFlags(c, createClientFlags)
		}
	},
}
