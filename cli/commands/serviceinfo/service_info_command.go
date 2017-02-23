package serviceinfo

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

var ServiceInfoCommand = cli.Command{
	Name:      "service-info",
	ShortName: "si",
	Usage:     "List info for a service instance",
	ArgsUsage: "APP_NAME SERVICE_INSTANCE",
	Before: func(c *cli.Context) error {
		if c.NArg() != 2 {
			return cli.NewExitError("Incorrect Usage", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		appName := c.Args()[0]
		instanceName := c.Args()[1]
		appNameColorized := terminal.EntityNameColor(appName)
		instanceNameColorized := terminal.EntityNameColor(instanceName)

		app := cf.Lookup.AppForName(appName)
		if app == nil {
			global.UI.Failed("App %s not found", appNameColorized)
		}
		instance := cf.Lookup.InstanceForInstanceName(instanceName)
		if instance == nil {
			global.UI.Failed("Service instance %s not found", instanceNameColorized)
		}
		helpers.ServiceInfo.PrintForAppAndServiceInstance(app, instance)
		return nil
	},
	BashComplete: func(c *cli.Context) {
		completionArgs := []string{}
		if c.NArg() == 0 {
			completionArgs = cf.Lookup.Apps()
		} else if c.NArg() == 1 {
			completionArgs = cf.Lookup.ServiceInstances()
		}
		for _, arg := range completionArgs {
			global.UI.Say(arg)
		}
	},
}

func init() {
	commandregistry.AddCommand(ServiceInfoCommand)
}
