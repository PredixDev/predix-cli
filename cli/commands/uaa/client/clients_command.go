package client

import (
	"encoding/json"

	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var clientsFlags = []cli.Flag{
	cli.IntFlag{
		Name:  "start",
		Usage: "Show results starting at this `index`",
		Value: 0,
	},
	cli.IntFlag{
		Name:  "count",
		Usage: "The `count` of results to show",
		Value: 100,
	},
}
var ClientsCommand = cli.Command{
	Name:  "clients",
	Usage: "List client registrations on the targeted Predix UAA instance",
	Flags: clientsFlags,
	Action: func(c *cli.Context) error {
		start := 0
		count := 100
		if !c.IsSet("start") {
			start = c.Int("start")
		}
		if !c.IsSet("count") {
			count = c.Int("count")
		}

		target, context, _ := uaac.Targets.GetCurrent()
		scim := lib.ScimFactory.New(target, context)
		clients, err := scim.GetClients(start, count)

		if err != nil {
			global.UI.Failed(err.Error())
		}

		clientsJSON, _ := json.MarshalIndent(clients, "", " ")
		global.UI.Say(string(clientsJSON))
		return nil
	},
	BashComplete: func(c *cli.Context) {
		helpers.Completions.PrintFlags(c, clientsFlags)
	},
}
