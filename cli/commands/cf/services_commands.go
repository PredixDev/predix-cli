package cf

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
)

func cfArg1ServiceInstance(nArgs int, arguments []string) (completionArgs []string) {
	completionArgs = []string{}
	if nArgs == 0 {
		completionArgs = cf.Lookup.ServiceInstances()
	}
	return completionArgs
}

func cfArg1MarketplaceServiceArg2ServicePlan(nArgs int, arguments []string) (completionArgs []string) {
	completionArgs = []string{}
	if nArgs == 0 {
		completionArgs = cf.Lookup.MarketplaceServices()
	} else if nArgs == 1 {
		completionArgs = cf.Lookup.MarketplaceServicePlans(arguments[0])
	}
	return completionArgs
}

func cfArg1AppNameArg2ServiceInstance(nArgs int, arguments []string) (completionArgs []string) {
	completionArgs = []string{}
	if nArgs == 0 {
		completionArgs = cf.Lookup.Apps()
	} else if nArgs == 1 {
		completionArgs = cf.Lookup.ServiceInstances()
	}
	return completionArgs
}

var cfServicesCommands = []CommandMetadata{
	{
		Name:       "marketplace",
		ShortName:  "m",
		Parameters: []string{"-s"},
		ParamLookup: map[string]ParamLookupFunc{
			"-s": cfLookupMarketplaceServices,
		},
	},
	{
		Name:      "services",
		ShortName: "s",
	},
	{
		Name:       "service",
		Parameters: []string{"--guid"},
		BoolParams: []string{"--guid"},
		Arguments:  cfArg1ServiceInstance,
	},
	{
		Name:       "create-service",
		ShortName:  "cs",
		Parameters: []string{"-c", "-t"},
		Arguments:  cfArg1MarketplaceServiceArg2ServicePlan,
	},
	{
		Name:       "update-service",
		Parameters: []string{"-p", "-c", "-t"},
		Arguments:  cfArg1ServiceInstance,
	},
	{
		Name:       "delete-service",
		ShortName:  "ds",
		Parameters: []string{"-f"},
		BoolParams: []string{"-f"},
		Arguments:  cfArg1ServiceInstance,
	},
	{
		Name:      "rename-service",
		Arguments: cfArg1ServiceInstance,
	},
	{
		Name:       "create-service-key",
		ShortName:  "csk",
		Parameters: []string{"-c"},
		Arguments:  cfArg1ServiceInstance,
	},
	{
		Name:      "service-keys",
		ShortName: "sk",
		Arguments: cfArg1ServiceInstance,
	},
	{
		Name:       "service-key",
		Parameters: []string{"--guid"},
		BoolParams: []string{"--guid"},
		Arguments:  cfArg1ServiceInstance,
	},
	{
		Name:       "delete-service-key",
		ShortName:  "dsk",
		Parameters: []string{"-f"},
		BoolParams: []string{"-f"},
		Arguments:  cfArg1ServiceInstance,
	},
	{
		Name:       "bind-service",
		ShortName:  "bs",
		Parameters: []string{"-c"},
		Arguments:  cfArg1AppNameArg2ServiceInstance,
	},
	{
		Name:      "unbind-service",
		ShortName: "us",
		Arguments: cfArg1AppNameArg2ServiceInstance,
	},
	{
		Name:       "create-user-provided-service",
		ShortName:  "cups",
		Parameters: []string{"-p", "-l"},
	},
	{
		Name:       "update-user-provided-service",
		ShortName:  "uups",
		Parameters: []string{"-p", "-l"},
		Arguments:  cfArg1ServiceInstance,
	},
}

func init() {
	for _, cmd := range cfServicesCommands {
		cfCommands = append(cfCommands, cmd.Name)
		cfCommandLookup[cmd.Name] = cmd
		if cmd.ShortName != "" {
			cfCommands = append(cfCommands, cmd.ShortName)
			cfCommandLookup[cmd.ShortName] = cmd
		}
	}
}
