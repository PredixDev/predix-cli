package createservice

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

func createViews(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	_, predixUaaServiceInfo := VerifyPredixUaaInstance(args)

	global.UI.Say("")
	instance := InstanceWithTrustedIssuerIDs(args, plan, predixUaaServiceInfo)

	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func init() {
	CreateServiceHandlers[constants.PredixViews] = createViews
}
