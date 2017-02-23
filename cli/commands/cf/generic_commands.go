package cf

var cfGenericCommands = []CommandMetadata{
	{
		Name: "help",
	},
	{
		Name: "version",
	},
	{
		Name:       "login",
		ShortName:  "l",
		Parameters: []string{"-a", "-u", "-p", "-o", "-s", "--sso", "--skip-ssl-validation"},
		BoolParams: []string{"--sso", "--skip-ssl-validation"},
		ParamLookup: map[string]ParamLookupFunc{
			"-o": cfLookupOrgs,
			"-s": cfLookupSpaces,
		},
	},
	{
		Name:      "logout",
		ShortName: "lo",
	},
	{
		Name:      "passwd",
		ShortName: "pw",
	},
	{
		Name:       "target",
		ShortName:  "t",
		Parameters: []string{"-o", "-s"},
		ParamLookup: map[string]ParamLookupFunc{
			"-o": cfLookupOrgs,
			"-s": cfLookupSpaces,
		},
	},
	{
		Name:       "api",
		Parameters: []string{"--unset", "--skip-ssl-validation"},
		BoolParams: []string{"--unset", "--skip-ssl-validation"},
	},
	{
		Name: "auth",
	},
}

func init() {
	for _, cmd := range cfGenericCommands {
		cfCommands = append(cfCommands, cmd.Name)
		cfCommandLookup[cmd.Name] = cmd
		if cmd.ShortName != "" {
			cfCommands = append(cfCommands, cmd.ShortName)
			cfCommandLookup[cmd.ShortName] = cmd
		}
	}
}
