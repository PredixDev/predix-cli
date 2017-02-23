package uaa

import (
	"strings"

	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var loginFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "ca-cert",
		Usage: "Use the given CA certificate `file` to validate the target's SSL certificate",
	},
	cli.BoolFlag{
		Name:  "skip-ssl-validation",
		Usage: "Do not attempt to validate the target's SSL certificate",
	},
	cli.StringFlag{
		Name:  "secret, s",
		Usage: "The `secret` for the client",
	},
	cli.StringFlag{
		Name:  "scope",
		Usage: "Comma separated list of `scopes` to request",
	},
	cli.StringFlag{
		Name:  "password, p",
		Usage: "The user `password`",
	},
}

var LoginCommand = cli.Command{
	Name:      "login",
	ShortName: "l",
	Usage:     "Login and target the specified Predix UAA instance",
	ArgsUsage: "UAA_INSTANCE CLIENT_ID [USERNAME]",
	Flags:     loginFlags,
	Before: func(c *cli.Context) error {
		if c.NArg() < 2 || c.NArg() > 3 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		caCertFile := c.String("ca-cert")
		skipSslVerify := c.Bool("skip-ssl-validation")
		clientSecret := c.String("secret")
		password := c.String("password")

		args := c.Args()
		instance, predixUaaServiceInfo := helpers.Uaa.FetchInstance(args[0])

		url := helpers.ServiceInfo.ResolveJSONPath(predixUaaServiceInfo, constants.PredixUaaURI).(string)
		infoClient := lib.InfoClientFactory.New(url, skipSslVerify, caCertFile)
		err := infoClient.Server()
		if err != nil {
			global.UI.Failed(err.Error())
		}

		var scopes []string = nil
		if c.String("scope") != "" {
			scopes = strings.Split(c.String("scope"), ",")
		}
		if clientSecret == "" {
			clientSecret = global.UI.AskForPassword("Client Secret")
		}

		issuer := lib.TokenIssuerFactory.New(url, args[1], clientSecret, skipSslVerify, caCertFile)

		var tr *lib.TokenResponse
		if c.NArg() == 2 {
			tr, err = issuer.ClientCredentialsGrant(scopes)
		} else {
			if password == "" {
				password = global.UI.AskForPassword("User Password")
			}
			tr, err = issuer.PasswordGrant(args[2], password, scopes)
		}
		if err != nil {
			global.UI.Failed(err.Error())
		}

		uaac.Targets.SetCurrent(url, instance.URL, skipSslVerify, caCertFile, tr)
		uaac.Targets.PrintCurrent()
		return nil
	},
	BashComplete: func(c *cli.Context) {
		nArg := c.NArg()
		completionArgs := []string{}
		if nArg == 0 {
			completionArgs = cf.Lookup.PredixUaaInstances()
		} else if nArg == 1 {
			completionArgs = []string{"Enter: CLIENT_ID", "_"}
		}
		for _, arg := range completionArgs {
			global.UI.Say(arg)
		}
		if len(completionArgs) == 0 {
			if nArg == 2 {
				global.UI.Say("[USERNAME]")
			}
			helpers.Completions.PrintFlags(c, loginFlags)
		}
	},
}
