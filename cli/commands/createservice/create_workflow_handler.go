package createservice

import (
	"fmt"

	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

type WorkflowParameters struct {
	TrustedIssuers []string `json:"trustedIssuerIds"`
	ClientID       string   `json:"clientId"`
	ClientSecret   string   `json:"clientSecret"`
}

func createWorkflow(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	parameters := WorkflowParameters{}

	predixUaaInstance, predixUaaServiceInfo := VerifyPredixUaaInstance(args)
	parameters.TrustedIssuers = []string{helpers.ServiceInfo.ResolveJSONPath(predixUaaServiceInfo, constants.PredixUaaIssuerID).(string)}

	clientID, clientSecret, client, scim := VerifyClientIDAndSecret(c, predixUaaInstance, predixUaaServiceInfo)
	if client != nil {
		clientSecret = AskForClientSecret(c)
	}
	parameters.ClientID = clientID
	parameters.ClientSecret = clientSecret

	global.UI.Say("")
	instance := InstanceWithParameters(args, plan, parameters)

	workflowServiceInfo := helpers.ServiceInfo.FetchFor(instance)
	scopes := []string{
		"azuqua.zones.read",
	}
	authorities := []string{
		fmt.Sprintf("azuqua.zones.%s.user", helpers.ServiceInfo.ResolveJSONPath(workflowServiceInfo, constants.WorkflowHTTPHeaderValue)),
	}

	CreateOrUpdateClient(clientID, clientSecret, client, scim, predixUaaInstance, scopes, authorities)
	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func init() {
	CreateServiceHandlers[constants.Workflow] = createWorkflow
}
