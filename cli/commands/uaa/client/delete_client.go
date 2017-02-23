package client

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var DeleteClientCommand = cli.Command{
	Name:      "delete",
	ShortName: "d",
	Usage:     "Delete a client registered with the targeted Predix UAA instance",
	ArgsUsage: "CLIENT_ID",
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

		global.UI.Say("Deleting client %s on service instance %s", clientIDColorized, uaaNameColorized)
		err := scim.DeleteClient(clientID)
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
		}
	},
}
