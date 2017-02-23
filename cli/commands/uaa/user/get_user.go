package user

import (
	"encoding/json"

	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var GetUserCommand = cli.Command{
	Name:      "get",
	ShortName: "g",
	Usage:     "Get a user account from the targeted Predix UAA instance",
	ArgsUsage: "USER_NAME",
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		target, context, _ := uaac.Targets.GetCurrent()
		scim := lib.ScimFactory.New(target, context)

		userName := c.Args()[0]
		userNameColorized := terminal.EntityNameColor(userName)

		user, err := scim.GetUser(userName)
		if err == nil {
			if user == nil {
				global.UI.Failed("User %s not found.", userNameColorized)
			}
		} else {
			global.UI.Failed(err.Error())
		}

		userJSON, _ := json.MarshalIndent(user, "", " ")
		global.UI.Say(string(userJSON))
		return nil
	},
	BashComplete: func(c *cli.Context) {
		if c.NArg() == 0 {
			global.UI.Say("Enter: USER_NAME")
			global.UI.Say("_")
		}
	},
}
