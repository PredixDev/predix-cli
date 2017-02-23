package user

import (
	"regexp"

	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

type CommandFlagDestinations struct {
	GivenName  string
	FamilyName string
	Emails     string
	Phones     string
	Password   string
}

var userCommandFlagDestinations = CommandFlagDestinations{}
var userFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "given_name",
		Usage:       "Given `name` of the user",
		Destination: &userCommandFlagDestinations.GivenName,
	},
	cli.StringFlag{
		Name:        "family_name",
		Usage:       "Family `name` of the user",
		Destination: &userCommandFlagDestinations.FamilyName,
	},
	cli.StringFlag{
		Name:        "emails",
		Usage:       "Comma separated list of email `addresses` of the user",
		Destination: &userCommandFlagDestinations.Emails,
	},
	cli.StringFlag{
		Name:        "phones",
		Usage:       "Comma separated list of `phone numbers` of the user",
		Destination: &userCommandFlagDestinations.Phones,
	},
}

func getCreateUserFlags() []cli.Flag {
	return append(userFlags, cli.StringFlag{
		Name:        "password, p",
		Usage:       "The `user password`",
		Destination: &userCommandFlagDestinations.Password,
	})
}

func getUpdateUserFlags() []cli.Flag {
	return userFlags
}

var splitBy = regexp.MustCompile(`[\s,]+`)

func split(s string) []lib.Value {
	if s != "" {
		arr := splitBy.Split(s, -1)
		values := make([]lib.Value, len(arr))
		for i, v := range arr {
			values[i] = lib.Value{
				Value: v,
			}
		}
		return values
	}
	return nil
}

func replace(a *[]lib.Value, s string) {
	v := split(s)
	if v != nil {
		*a = v
	}
}
