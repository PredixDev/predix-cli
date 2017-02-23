package user

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var updateUserFlags = getUpdateUserFlags()
var UpdateUserCommand = cli.Command{
	Name:      "update",
	ShortName: "u",
	Usage:     "Update a user account on the targeted Predix UAA instance",
	ArgsUsage: "USER_NAME",
	Flags:     updateUserFlags,
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

		global.UI.Say("Checking if user %s exists on service instance %s", userNameColorized, uaaNameColorized)
		user, err := scim.GetUser(userName)

		if err == nil {
			if user == nil {
				global.UI.Failed("User %s not found.", userNameColorized)
			}
		} else {
			global.UI.Failed(err.Error())
		}
		global.UI.Ok()

		if userCommandFlagDestinations.GivenName != "" {
			user.Name.GivenName = userCommandFlagDestinations.GivenName
		}
		if userCommandFlagDestinations.FamilyName != "" {
			user.Name.FamilyName = userCommandFlagDestinations.FamilyName
		}
		replace(&user.Emails, userCommandFlagDestinations.Emails)
		replace(&user.Phones, userCommandFlagDestinations.Phones)

		global.UI.Say("")
		global.UI.Say("Updating user %s on service instance %s", userNameColorized, uaaNameColorized)
		err = scim.PutUser(user)
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
		} else {
			helpers.Completions.PrintFlags(c, updateUserFlags)
		}
	},
}
