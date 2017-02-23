package user

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var DeleteUserCommand = cli.Command{
	Name:      "delete",
	ShortName: "d",
	Usage:     "Delete a user account on the targeted Predix UAA instance",
	ArgsUsage: "USER_NAME",
	Before: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		target, context, instance := uaac.Targets.GetCurrent()
		scim := lib.ScimFactory.New(target, context)

		userName := c.Args()[0]
		userNameColorized := terminal.EntityNameColor(userName)
		uaaNameColorized := terminal.EntityNameColor(instance.Name)

		user, err := scim.GetUser(userName)
		if err == nil {
			if user == nil {
				global.UI.Failed("User %s not found.", userNameColorized)
			}
		} else {
			global.UI.Failed(err.Error())
		}

		global.UI.Say("Deleting user %s on service instance %s", userNameColorized, uaaNameColorized)
		err = scim.DeleteUser(user.ID)
		if err != nil {
			global.UI.Failed(err.Error())
		}
		global.UI.Ok()

		return nil
	},
	BashComplete: func(c *cli.Context) {
		if c.NArg() == 0 {
			global.UI.Say("Enter: USER_NAME")
			global.UI.Say("_")
		}
	},
}
