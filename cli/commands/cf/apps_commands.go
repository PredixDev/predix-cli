package cf

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
)

func cfArg1AppName(nArgs int, arguments []string) (completionArgs []string) {
	completionArgs = []string{}
	if nArgs == 0 {
		completionArgs = cf.Lookup.Apps()
	}
	return completionArgs
}

func cfArg1And2AppName(nArgs int, arguments []string) (completionArgs []string) {
	completionArgs = []string{}
	if nArgs == 0 || nArgs == 1 {
		completionArgs = cf.Lookup.Apps()
	}
	return completionArgs
}

var cfAppsCommands = []CommandMetadata{
	{
		Name:      "apps",
		ShortName: "a",
	},
	{
		Name:       "app",
		Parameters: []string{"--guid"},
		BoolParams: []string{"--guid"},
		Arguments:  cfArg1AppName,
	},
	{
		Name:       "push",
		ShortName:  "p",
		Parameters: []string{"-b", "-c", "-d", "-f", "-i", "-k", "-m", "-n", "-p", "-s", "-t", "--no-hostname", "--no-manifest", "--no-route", "--no-start", "--random-route"},
		BoolParams: []string{"--no-hostname", "--no-manifest", "--no-route", "--no-start", "--random-route"},
		Arguments:  cfArg1AppName,
		PostExec:   cf.Cache.InvalidateApps,
		ParamLookup: map[string]ParamLookupFunc{
			"-s": cfLookupStacks,
		},
	},
	{
		Name:       "scale",
		Parameters: []string{"-i", "-k", "-m", "-f"},
		BoolParams: []string{"-f"},
		Arguments:  cfArg1AppName,
	},
	{
		Name:       "delete",
		ShortName:  "d",
		Parameters: []string{"-f", "-r"},
		BoolParams: []string{"-f", "-r"},
		Arguments:  cfArg1AppName,
		PostExec:   cf.Cache.InvalidateApps,
	},
	{
		Name:      "rename",
		Arguments: cfArg1AppName,
		PostExec:  cf.Cache.InvalidateApps,
	},
	{
		Name:      "start",
		ShortName: "st",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "stop",
		ShortName: "sp",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "restart",
		ShortName: "rs",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "restage",
		ShortName: "rg",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "restart-app-instance",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "events",
		Arguments: cfArg1AppName,
	},
	{
		Name:       "files",
		ShortName:  "f",
		Parameters: []string{"-i"},
		Arguments:  cfArg1AppName,
	},
	{
		Name:       "logs",
		Parameters: []string{"--recent"},
		BoolParams: []string{"--recent"},
		Arguments:  cfArg1AppName,
	},
	{
		Name:      "env",
		ShortName: "e",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "set-env",
		ShortName: "se",
		Arguments: cfArg1AppName,
	},
	{
		Name:      "unset-env",
		Arguments: cfArg1AppName,
	},
	{
		Name: "stacks",
	},
	{
		Name:       "stack",
		Parameters: []string{"--guid"},
		BoolParams: []string{"--guid"},
		Arguments: func(nArgs int, arguments []string) (completionArgs []string) {
			completionArgs = []string{}
			if nArgs == 0 {
				completionArgs = cfLookupStacks(nil)
			}
			return completionArgs
		},
	},
	{
		Name:       "copy-source",
		Parameters: []string{"-o", "-s", "--no-restart"},
		BoolParams: []string{"--no-restart"},
		Arguments:  cfArg1And2AppName,
		ParamLookup: map[string]ParamLookupFunc{
			"-o": cfLookupOrgs,
			"-s": cfLookupSpaces,
		},
	},
	{
		Name:       "create-app-manifest",
		Parameters: []string{"-p"},
		Arguments:  cfArg1AppName,
	},
}

func init() {
	for _, cmd := range cfAppsCommands {
		cfCommands = append(cfCommands, cmd.Name)
		cfCommandLookup[cmd.Name] = cmd
		if cmd.ShortName != "" {
			cfCommands = append(cfCommands, cmd.ShortName)
			cfCommandLookup[cmd.ShortName] = cmd
		}
	}
}
