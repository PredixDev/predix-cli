package commands

import (
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

var loginFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "u",
		Usage: "`Username`",
	},
	cli.StringFlag{
		Name:  "p",
		Usage: "`Password`",
	},
	cli.StringFlag{
		Name:  "o",
		Usage: "`Org`",
	},
	cli.StringFlag{
		Name:  "s",
		Usage: "`Space`",
	},
	cli.BoolFlag{
		Name:  "skip-ssl-validation",
		Usage: "Skip verification of the API endpoint. Not recommended!",
	},
	cli.BoolFlag{
		Name:  "sso",
		Usage: "Use a one-time password to login",
	},
	cli.BoolFlag{
		Name:  "us-west",
		Usage: "Login to the Predix US-West PoP",
	},
	cli.BoolFlag{
		Name:  "us-east",
		Usage: "Login to the Predix US-East PoP",
	},
	cli.BoolFlag{
		Name:  "japan",
		Usage: "Login to the Predix Japan PoP",
	},
	cli.BoolFlag{
		Name:  "uk",
		Usage: "Login to the Predix UK PoP",
	},
}
var LoginCommand = cli.Command{
	Name:      "login",
	ShortName: "l",
	Usage:     "Log user in to the Predix Platform",
	Flags:     loginFlags,
	Action: func(c *cli.Context) error {
		var choice string
		if c.Bool("us-west") {
			choice = "1"
		} else if c.Bool("us-east") {
			choice = "2"
		} else if c.Bool("japan") {
			choice = "3"
		} else if c.Bool("uk") {
			choice = "4"
		}

		if choice == "" {
			global.UI.Say("1. US-West")
			global.UI.Say("2. US-East")
			global.UI.Say("3. Japan")
			global.UI.Say("4. UK")
			choice = global.UI.Ask("Choose the PoP to set")
		}

		var api string
		switch choice {
		case "1":
			api = "https://api.system.aws-usw02-pr.ice.predix.io"
		case "2":
			api = "https://api.system.asv-pr.ice.predix.io"
		case "3":
			api = "https://api.system.aws-jp01-pr.ice.predix.io"
		case "4":
			api = "https://api.system.dc-uk01-pr.ice.predix.io"
		default:
			return cli.NewExitError("Invalid choice!", 1)
		}

		loginArgs := []string{"login", "-a", api}
		if c.IsSet("u") {
			loginArgs = append(loginArgs, []string{"-u", c.String("u")}...)
		}
		if c.IsSet("p") {
			loginArgs = append(loginArgs, []string{"-p", c.String("p")}...)
		}
		if c.IsSet("o") {
			loginArgs = append(loginArgs, []string{"-o", c.String("o")}...)
		}
		if c.IsSet("s") {
			loginArgs = append(loginArgs, []string{"-s", c.String("s")}...)
		}
		if c.Bool("skip-ssl-validation") {
			loginArgs = append(loginArgs, "--skip-ssl-validation")
		}
		if c.Bool("sso") {
			loginArgs = append(loginArgs, "--sso")
		}
		var err = global.Sh.Interactive().Command("cf", loginArgs).Run()
		if err != nil {
			global.UI.Failed("Error while logging user in")
		}
		return nil
	},
}

func init() {
	commandregistry.AddCommand(LoginCommand)
}
