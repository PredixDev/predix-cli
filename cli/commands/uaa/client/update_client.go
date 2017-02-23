package client

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var updateClientFlags = getUpdateClientFlags()
var UpdateClientCommand = cli.Command{
	Name:      "update",
	ShortName: "u",
	Usage:     "Update a client registered with the targeted Predix UAA instance",
	ArgsUsage: "CLIENT_ID",
	Flags:     updateClientFlags,
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		target, context, instance := uaac.Targets.GetCurrent()
		scim := lib.ScimFactory.New(target, context)

		clientID := c.Args()[0]
		clientIDColorized := terminal.EntityNameColor(clientID)
		uaaNameColorized := terminal.EntityNameColor(instance.Name)

		global.UI.Say("Checking if client %s exists on service instance %s", clientIDColorized, uaaNameColorized)
		client, err := scim.GetClient(clientID)
		if err == nil {
			if client == nil {
				global.UI.Failed("Client %s not found.", clientIDColorized)
			}
		} else {
			global.UI.Failed(err.Error())
		}
		global.UI.Ok()

		if clientCommandFlagDestinations.Name != "" {
			client.Name = clientCommandFlagDestinations.Name
		}
		replace(&client.Scopes, clientCommandFlagDestinations.Scopes)
		replace(&client.GrantTypes, clientCommandFlagDestinations.Grants)
		replace(&client.Authorities, clientCommandFlagDestinations.Authorities)
		timeout := atoi(clientCommandFlagDestinations.AccessTokenTimeout)
		if timeout != 0 {
			client.AccessTokenTimeout = timeout
		}
		timeout = atoi(clientCommandFlagDestinations.RefreshTokenTimeout)
		if timeout != 0 {
			client.RefreshTokenTimeout = timeout
		}
		replace(&client.RedirectURI, clientCommandFlagDestinations.RedirectURI)
		replace(&client.AutoApprove, clientCommandFlagDestinations.AutoApprove)
		if clientCommandFlagDestinations.SignupRedirect != "" {
			client.SignupRedirect = clientCommandFlagDestinations.SignupRedirect
		}

		global.UI.Say("")
		global.UI.Say("Updating client %s on service instance %s", clientIDColorized, uaaNameColorized)
		err = scim.PutClient(client)
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
			helpers.Completions.PrintFlags(c, updateClientFlags)
		}
	},
}
