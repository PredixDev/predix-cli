package client

import (
	"encoding/json"

	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var GetClientCommand = cli.Command{
	Name:      "get",
	ShortName: "g",
	Usage:     "Get a client registered with the targeted Predix UAA instance",
	ArgsUsage: "CLIENT_ID",
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		target, context, _ := uaac.Targets.GetCurrent()
		scim := lib.ScimFactory.New(target, context)

		clientID := c.Args()[0]
		clientIDColorized := terminal.EntityNameColor(clientID)

		client, err := scim.GetClient(clientID)
		if err == nil {
			if client == nil {
				global.UI.Failed("Client %s not found.", clientIDColorized)
			}
		} else {
			global.UI.Failed(err.Error())
		}

		clientJSON, _ := json.MarshalIndent(client, "", " ")
		global.UI.Say(string(clientJSON))
		return nil
	},
	BashComplete: func(c *cli.Context) {
		if c.NArg() == 0 {
			global.UI.Say("Enter: CLIENT_ID")
			global.UI.Say("_")
		}
	},
}
