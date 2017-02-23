package createservice

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

func createTimeseries(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	predixUaaInstance, predixUaaServiceInfo := VerifyPredixUaaInstance(args)
	clientID, clientSecret, client, scim := VerifyClientIDAndSecret(c, predixUaaInstance, predixUaaServiceInfo)

	global.UI.Say("")
	instance := InstanceWithTrustedIssuerIDs(args, plan, predixUaaServiceInfo)
	predixTimeseriesServiceInfo := helpers.ServiceInfo.FetchFor(instance)

	auths := helpers.ServiceInfo.ResolveJSONPath(predixTimeseriesServiceInfo, constants.PredixTimeseriesIngestOauthScopes).([]interface{})
	auths = append(auths, helpers.ServiceInfo.ResolveJSONPath(predixTimeseriesServiceInfo, constants.PredixTimeseriesQueryOauthScopes).([]interface{})...)

	authorities := []string{}
	for _, v := range auths {
		authorities = append(authorities, v.(string))
	}

	CreateOrUpdateClient(clientID, clientSecret, client, scim, predixUaaInstance, nil, authorities)
	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func init() {
	CreateServiceHandlers[constants.PredixTimeseries] = createTimeseries
}
