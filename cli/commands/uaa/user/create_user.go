package user

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var createUserFlags = getCreateUserFlags()
var CreateUserCommand = cli.Command{
	Name:      "create",
	ShortName: "c",
	Usage:     "Create a user account on the targeted Predix UAA instance",
	ArgsUsage: "USER_NAME",
	Flags:     createUserFlags,
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

		user := lib.User{}
		user.UserName = c.Args()[0]
		user.Password = userCommandFlagDestinations.Password
		if user.Password == "" {
			user.Password = global.UI.AskForVerifiedPassword("Password")
		}
		if userCommandFlagDestinations.GivenName != "" || userCommandFlagDestinations.FamilyName != "" {
			user.Name = lib.Name{
				GivenName:  userCommandFlagDestinations.GivenName,
				FamilyName: userCommandFlagDestinations.FamilyName,
			}
		}
		user.Emails = split(userCommandFlagDestinations.Emails)
		user.Phones = split(userCommandFlagDestinations.Phones)

		global.UI.Say("Creating user %s on service instance %s", terminal.EntityNameColor(user.UserName),
			terminal.EntityNameColor(instance.Name))
		err = scim.CreateUser(&user)
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
			helpers.Completions.PrintFlags(c, createUserFlags)
		}
	},
}
