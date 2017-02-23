package client

import (
	"regexp"
	"strconv"

	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

type CommandFlagDestinations struct {
	Name                string
	Scopes              string
	Grants              string
	Authorities         string
	AccessTokenTimeout  string
	RefreshTokenTimeout string
	RedirectURI         string
	AutoApprove         string
	SignupRedirect      string
	Secret              string
}

var clientCommandFlagDestinations = CommandFlagDestinations{}
var clientFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "name",
		Usage:       "The `name` of the client",
		Destination: &clientCommandFlagDestinations.Name,
	},
	cli.StringFlag{
		Name:        "scope",
		Usage:       "Comma separated list of `scopes` on the client",
		Destination: &clientCommandFlagDestinations.Scopes,
	},
	cli.StringFlag{
		Name:        "authorized_grant_types",
		Usage:       "Comma separated list of `grant-types` for the client",
		Destination: &clientCommandFlagDestinations.Grants,
	},
	cli.StringFlag{
		Name:        "authorities",
		Usage:       "Comma separated list of `authorities` on the client",
		Destination: &clientCommandFlagDestinations.Authorities,
	},
	cli.StringFlag{
		Name:        "access_token_validity",
		Usage:       "Validity of the access token in seconds",
		Destination: &clientCommandFlagDestinations.AccessTokenTimeout,
	},
	cli.StringFlag{
		Name:        "refresh_token_validity",
		Usage:       "Validity of the refresh token in seconds",
		Destination: &clientCommandFlagDestinations.RefreshTokenTimeout,
	},
	cli.StringFlag{
		Name:        "redirect_uri",
		Usage:       "Comma separated list of `uri` to redirect to on grant/deny of access request",
		Destination: &clientCommandFlagDestinations.RedirectURI,
	},
	cli.StringFlag{
		Name:        "autoapprove",
		Usage:       "Comma separated list of auto approve `scopes` on the client",
		Destination: &clientCommandFlagDestinations.AutoApprove,
	},
	cli.StringFlag{
		Name:        "signup_redirect_url",
		Usage:       "The `url` to redirect to when signing up",
		Destination: &clientCommandFlagDestinations.SignupRedirect,
	},
}

func getCreateClientFlags() []cli.Flag {
	return append(clientFlags, cli.StringFlag{
		Name:        "secret, s",
		Usage:       "The `secret` for the client",
		Destination: &clientCommandFlagDestinations.Secret,
	})
}

func getUpdateClientFlags() []cli.Flag {
	return clientFlags
}

var splitBy = regexp.MustCompile(`[\s,]+`)

func split(s string) []string {
	if s != "" {
		return splitBy.Split(s, -1)
	}
	return nil
}

func atoi(s string) int {
	if s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			global.UI.Failed("Invalid value. %s", err.Error())
		}
		return v
	}
	return 0
}

func replace(a *[]string, s string) {
	v := split(s)
	if v != nil {
		*a = v
	}
}
