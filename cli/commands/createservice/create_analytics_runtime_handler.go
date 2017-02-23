package createservice

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

type AnalyticsRuntimeServiceParameters struct {
	TrustedIssuers    []string `json:"trustedIssuerIds"`
	DependentServices struct {
		Asset            string `json:"predixAssetZoneId,omitempty"`
		Timeseries       string `json:"predixTimeseriesZoneId,omitempty"`
		AnalyticsCatalog string `json:"predixAnalyticsCatalogZoneId,omitempty"`
	} `json:"dependentServices,omitempty"`
	TrustedClient struct {
		ID     string `json:"clientId,omitempty"`
		Secret string `json:"clientSecret,omitempty"`
	} `json:"trustedClientCredential,omitempty"`
}

func createAnalyticsRuntime(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	parameters := AnalyticsRuntimeServiceParameters{}

	predixUaaInstance, predixUaaServiceInfo := VerifyPredixUaaInstance(args)
	parameters.TrustedIssuers = []string{helpers.ServiceInfo.ResolveJSONPath(predixUaaServiceInfo, constants.PredixUaaIssuerID).(string)}

	clientID, clientSecret, client, scim := VerifyClientIDAndSecret(c, predixUaaInstance, predixUaaServiceInfo)

	hasDependent := false
	nArg := c.NArg()
	if nArg >= 5 {
		parameters.DependentServices.Asset = fetchInstanceGUID(args[4], "Predix Asset")
		hasDependent = true
	}
	if nArg >= 6 {
		parameters.DependentServices.Timeseries = fetchInstanceGUID(args[5], "Predix Timeseries")
	}
	if nArg >= 7 {
		parameters.DependentServices.AnalyticsCatalog = fetchInstanceGUID(args[6], "Predix Analytics Catalog")
	}

	if hasDependent {
		if client != nil {
			clientSecret = AskForClientSecret(c)
			parameters.TrustedClient.ID = clientID
			parameters.TrustedClient.Secret = clientSecret
		} else {
			global.UI.Failed("A trusted client is required for the dependent services")
		}
	}

	global.UI.Say("")
	instance := InstanceWithParameters(args, plan, parameters)
	predixAnalyticsRuntimeServiceInfo := helpers.ServiceInfo.FetchFor(instance)

	authorities := []string{helpers.ServiceInfo.ResolveJSONPath(predixAnalyticsRuntimeServiceInfo, constants.PredixAnalyticsRuntimeOauthScopes).(string)}

	CreateOrUpdateClient(clientID, clientSecret, client, scim, predixUaaInstance, nil, authorities)
	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func fetchInstanceGUID(instanceName string, service string) string {
	instance := cf.Lookup.InstanceForInstanceName(instanceName)
	if instance == nil {
		global.UI.Failed("%s service instance %s not found", service, terminal.EntityNameColor(instanceName))
	}
	return instance.GUID
}

func init() {
	CreateServiceHandlers[constants.PredixAnalyticsRuntime] = createAnalyticsRuntime
}
