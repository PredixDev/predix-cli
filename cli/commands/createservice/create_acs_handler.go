package createservice

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

func createAcs(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	predixUaaInstance, predixUaaServiceInfo := VerifyPredixUaaInstance(args)
	SetCurrentTarget(c, predixUaaInstance, predixUaaServiceInfo)
	clientID, clientSecret, client, scim := VerifyClientIDAndSecret(c, predixUaaInstance, predixUaaServiceInfo)

	global.UI.Say("")
	instance := InstanceWithTrustedIssuerIDs(args, plan, predixUaaServiceInfo)
	predixAcsServiceInfo := helpers.ServiceInfo.FetchFor(instance)

	authorities := []string{
		"acs.policies.read",
		"acs.policies.write",
		"acs.attributes.read",
		"acs.attributes.write",
		helpers.ServiceInfo.ResolveJSONPath(predixAcsServiceInfo, constants.PredixAcsOauthScopes).(string),
	}

	CreateOrUpdateClient(clientID, clientSecret, client, scim, predixUaaInstance, nil, authorities)
	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func init() {
	CreateServiceHandlers[constants.PredixAcs] = createAcs
}
