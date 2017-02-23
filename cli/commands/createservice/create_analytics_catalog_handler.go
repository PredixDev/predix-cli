package createservice

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

func createAnalyticsCatalog(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	predixUaaInstance, predixUaaServiceInfo := VerifyPredixUaaInstance(args)
	clientID, clientSecret, client, scim := VerifyClientIDAndSecret(c, predixUaaInstance, predixUaaServiceInfo)

	global.UI.Say("")
	instance := InstanceWithTrustedIssuerIDs(args, plan, predixUaaServiceInfo)
	predixAnalyticsCatalogServiceInfo := helpers.ServiceInfo.FetchFor(instance)

	authorities := []string{helpers.ServiceInfo.ResolveJSONPath(predixAnalyticsCatalogServiceInfo, constants.PredixAnalyticsCatalogOauthScopes).(string)}

	CreateOrUpdateClient(clientID, clientSecret, client, scim, predixUaaInstance, nil, authorities)
	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func init() {
	CreateServiceHandlers[constants.PredixAnalyticsCatalog] = createAnalyticsCatalog
}
