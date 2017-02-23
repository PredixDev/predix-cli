package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.build.ge.com/adoption/cli-lib/panicprinter"
	"github.build.ge.com/adoption/cli-lib/shell"
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.build.ge.com/adoption/predix-cli/cli/commands/cf"
	"github.build.ge.com/adoption/predix-cli/cli/commandsloader"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func Main(args []string) {
	userDir, err := homedir.Dir()
	if err == nil {
		global.Env.ConfigDir = filepath.Join(userDir, ".predix")
		_ = os.MkdirAll(global.Env.ConfigDir, os.FileMode(0700))
	}
	global.Env.CfHomeDir = cfHomeDir()

	app := cli.NewApp()
	app.Writer = Writer
	app.EnableBashCompletion = true
	app.Name = global.Name
	app.Version = fmt.Sprintf("predix-labs-%s-beta", global.Version)
	app.Usage = "A command line tool to interact with the Predix platform"
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		noCfBypass, err := strconv.ParseBool(os.Getenv("PREDIX_NO_CF_BYPASS"))
		if err != nil || !noCfBypass {
			if err != nil && !c.Bool(cli.BashCompletionFlag.Name) {
				global.UI.Say("'%s' is not a registered Predix CLI command. Trying to run it as a CF CLI command.", terminal.CommandColor(cmd))
				global.UI.Say(terminal.AdvisoryColor("Note: To disable this message, set the PREDIX_NO_CF_BYPASS environment variable to 'false'"))
			}
			cli.HandleExitCoder(cf.Action(c))
		} else {
			global.UI.Say("'%s' is not a registered Predix CLI command.", terminal.CommandColor(cmd))
			os.Exit(1)
		}
	}

	defer panicprinter.HandlePanics("Predix CLI", app.Version, "https://github.com/PredixDev/predix-cli",
		"https://github.com/PredixDev/predix-cli/issues", global.UI)

	commandsloader.Load()
	app.Commands = commandregistry.GetCommands()

	uaac.Targets.LoadConfig()

	_ = app.Run(args)
}

func cfHomeDir() string {
	var cfDir string

	if os.Getenv("CF_HOME") != "" {
		cfDir = os.Getenv("CF_HOME")

		if _, err := os.Stat(cfDir); os.IsNotExist(err) {
			return ""
		}
	} else {
		dir, err := homedir.Dir()
		if err != nil {
			return ""
		}
		cfDir = dir
	}

	return cfDir
}

func init() {
	global.Sh = shell.NewShell()

	teePrinter := terminal.NewTeePrinter(Writer)
	global.UI = terminal.NewUI(os.Stdin, Writer, teePrinter)
}
