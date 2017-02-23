package createservice

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/urfave/cli"
)

var CreateServiceHandlers = map[string]func(*cli.Context){}

var createServiceFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "admin-secret, a",
		Usage: "The admin client `SECRET`",
	},
	cli.StringFlag{
		Name:  "client-id, c",
		Usage: "The `CLIENT-ID` to set scopes on",
	},
	cli.StringFlag{
		Name:  "client-secret, s",
		Usage: "The `CLIENT-SECRET` for the specified client-id",
	},
	cli.BoolFlag{
		Name:  "skip-ssl-validation",
		Usage: "Do not attempt to validate the target's SSL certificate",
	},
	cli.StringFlag{
		Name:  "ca-cert",
		Usage: "Use the given CA certificate `file` to validate the target's SSL certificate",
	},
}

var CreateServiceCommand = cli.Command{
	Name:            "create-service",
	ShortName:       "cs",
	Usage:           "Create a service instance",
	ArgsUsage:       "SERVICE PLAN SERVICE_INSTANCE <UAA_INSTANCE> <ASSET_INSTANCE> <TIMESERIES_INSTANCE> <ANALYTICS_CATALOG_INSTANCE>",
	Flags:           createServiceFlags,
	SkipFlagParsing: false,
	Before: func(c *cli.Context) error {
		if c.NArg() < 3 {
			return cli.NewExitError("Incorrect Usage", 1)
		} else if strings.Compare(c.Args()[0], constants.PredixUaa) != 0 && c.NArg() < 4 {
			return cli.NewExitError("Incorrect Usage", 1)
		} else if handler := CreateServiceHandlers[c.Args()[0]]; handler == nil {
			global.UI.Failed("No handler found for service '%s'", c.Args()[0])
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		global.UI.Say("")
		handler := CreateServiceHandlers[c.Args()[0]]
		if handler != nil {
			handler(c)
		} else {
			global.UI.Failed("No handler found for service '%s'", c.Args()[0])
		}
		return nil
	},
	BashComplete: func(c *cli.Context) {
		nArg := c.NArg()
		completionArgs := []string{}
		if nArg == 0 {
			completionArgs = cf.Lookup.MarketplaceServices()
		} else if nArg == 1 {
			completionArgs = cf.Lookup.MarketplaceServicePlans(c.Args()[0])
		} else if nArg == 2 {
			completionArgs = []string{"Enter: NAME_FOR_NEW_SERVICE_INSTANCE", "_"}
		} else if strings.Compare(c.Args()[0], constants.PredixUaa) != 0 {
			if nArg == 3 {
				completionArgs = cf.Lookup.PredixUaaInstances()
			} else if strings.Compare(c.Args()[0], constants.PredixAnalyticsRuntime) == 0 {
				if nArg == 4 {
					completionArgs = cf.Lookup.PredixAssetInstances()
				} else if nArg == 5 {
					completionArgs = cf.Lookup.PredixTimeseriesInstances()
				} else if nArg == 6 {
					completionArgs = cf.Lookup.PredixAnalyticsCatalogInstances()
				}
			}
		}
		for _, arg := range completionArgs {
			global.UI.Say(arg)
		}
		if len(completionArgs) == 0 {
			helpers.Completions.PrintFlags(c, createServiceFlags)
		}
	},
}

func AskForClientID(c *cli.Context) string {
	clientID := c.String("client-id")
	if clientID == "" {
		clientID = global.UI.Ask("Client ID")
	}
	return clientID
}

func AskForClientSecret(c *cli.Context) string {
	clientSecret := c.String("client-secret")
	if clientSecret == "" {
		clientSecret = global.UI.AskForVerifiedPassword("Client Secret")
	}
	return clientSecret
}

func VerifyServicePlan(c *cli.Context) (args cli.Args, plan *cf.Item) {
	args = c.Args()
	plan = cf.Lookup.MarketplaceServicePlanItem(args[0], args[1])

	if plan == nil {
		global.UI.Failed("Service plan not found")
	}
	return args, plan
}

func VerifyPredixUaaInstance(args cli.Args) (instance *cf.Item, serviceInfo map[string]interface{}) {
	return helpers.Uaa.FetchInstance(args[3])
}

func SetCurrentTarget(c *cli.Context, instance *cf.Item, serviceInfo map[string]interface{}) {
	admin := "admin"
	url := serviceInfo["uri"].(string)

	if !uaac.Targets.LookupAndSetCurrent(url, admin) {
		caCertFile := c.String("ca-cert")
		skipSslVerify := c.Bool("skip-ssl-validation")
		adminClientSecret := helpers.Uaa.AskForAdminClientSecret(c)

		issuer := lib.TokenIssuerFactory.New(url, admin, adminClientSecret, skipSslVerify, caCertFile)
		tr, err := issuer.ClientCredentialsGrant(nil)

		if err != nil {
			global.UI.Failed(err.Error())
		}

		uaac.Targets.SetCurrent(url, instance.URL, skipSslVerify, caCertFile, tr)
	}
	uaac.Targets.PrintCurrent()
}

func VerifyClientIDAndSecret(c *cli.Context, predixUaaInstance *cf.Item, predixUaaServiceInfo map[string]interface{}) (clientID, clientSecret string, client *lib.Client, scim lib.Scim) {
	target, context, instance := uaac.Targets.GetCurrent()
	if predixUaaInstance.URL != instance.URL {
		global.UI.Failed("Incorrect target UAA, should be %s", terminal.EntityNameColor(predixUaaInstance.Name))
	}
	scim = lib.ScimFactory.New(target, context)

	clientID = AskForClientID(c)
	clientIDColorized := terminal.EntityNameColor(clientID)
	uaaNameColorized := terminal.EntityNameColor(instance.Name)

	global.UI.Say("Checking if client %s exists on service instance %s", clientIDColorized, uaaNameColorized)
	client, err := scim.GetClient(clientID)

	if err == nil {
		if client != nil {
			global.UI.Say("Client %s exists. The required authorities will be added to it.", clientIDColorized)
		} else {
			global.UI.Say("Client %s does not exist. It will be created with the required authorities.", clientIDColorized)
			clientSecret = AskForClientSecret(c)
		}
	} else {
		global.UI.Failed(err.Error())
	}
	return clientID, clientSecret, client, scim
}

func SayCreatingServiceInstance(args cli.Args) {
	global.UI.Say("Creating service instance %s in org %s / space %s as %s", terminal.EntityNameColor(args[2]),
		terminal.EntityNameColor(cf.CurrentUserInfo().Org), terminal.EntityNameColor(cf.CurrentUserInfo().Space),
		terminal.EntityNameColor(cf.CurrentUserInfo().Name))
}

func InstanceWithTrustedIssuerIDs(args cli.Args, plan *cf.Item, predixUaaServiceInfo map[string]interface{}) (instance *cf.Item) {
	SayCreatingServiceInstance(args)

	instance, err := cf.Curl.PostItem("/v2/service_instances?accepts_incomplete=true",
		fmt.Sprintf(`{"name":"%s","space_guid":"%s","service_plan_guid":"%s","parameters":{"trustedIssuerIds":["%s"]}}`,
			args[2], cf.CurrentUserInfo().SpaceGUID, plan.GUID, helpers.ServiceInfo.ResolveJSONPath(predixUaaServiceInfo, constants.PredixUaaIssuerID)))

	if err != nil {
		global.UI.Failed(err.Error())
	} else {
		global.UI.Ok()
	}
	return instance
}

func InstanceWithParameters(args cli.Args, plan *cf.Item, parameters interface{}) (instance *cf.Item) {
	SayCreatingServiceInstance(args)

	parametersJSON, err := json.Marshal(parameters)
	if err == nil {
		instance, err = cf.Curl.PostItem("/v2/service_instances?accepts_incomplete=true",
			fmt.Sprintf(`{"name":"%s","space_guid":"%s","service_plan_guid":"%s","parameters":%s}`,
				args[2], cf.CurrentUserInfo().SpaceGUID, plan.GUID, parametersJSON))
	}

	if err != nil {
		global.UI.Failed(err.Error())
	} else {
		global.UI.Ok()
	}
	return instance
}

func CreateOrUpdateClient(clientID string, clientSecret string, client *lib.Client, scim lib.Scim, predixUaaInstance *cf.Item, scopes, authorities []string) {
	clientIDColorized := terminal.EntityNameColor(clientID)
	uaaNameColorized := terminal.EntityNameColor(predixUaaInstance.Name)
	clientExists := client != nil

	global.UI.Say("")

	if clientExists {
		global.UI.Say("Updating client %s on Predix UAA instance %s", clientIDColorized, uaaNameColorized)
	} else {
		global.UI.Say("Creating client %s on Predix UAA instance %s", clientIDColorized, uaaNameColorized)
		client = uaac.NewSimpleClient(clientID, clientSecret)
	}

	if scopes != nil {
		client.Scopes = append(client.Scopes, scopes...)
	}
	if authorities != nil {
		client.Authorities = append(client.Authorities, authorities...)
	}

	var err error
	if clientExists {
		err = scim.PutClient(client)
	} else {
		err = scim.CreateClient(client)
	}
	if err != nil {
		global.UI.Failed(err.Error())
	}

	global.UI.Ok()
}

func init() {
	commandregistry.AddCommand(CreateServiceCommand)
}
