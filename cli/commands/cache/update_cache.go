package cache

import (
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

var updateCacheCommand = cli.Command{
	Name:  "update",
	Usage: "Update the Predix CLI cache",
	Action: func(c *cli.Context) error {
		cf.Cache.PurgeCurrent()
		global.UI.Say("Caching Orgs... ")
		orgs := cf.Lookup.Orgs()
		global.UI.Ok()
		global.UI.Say("Caching Spaces... ")
		for _, org := range orgs {
			cf.Lookup.Spaces(map[string]string{
				"-o": org,
			})
		}
		global.UI.Ok()
		global.UI.Say("Caching Stacks... ")
		cf.Lookup.Stacks()
		global.UI.Ok()
		global.UI.Say("Caching Services... ")
		services := cf.Lookup.MarketplaceServices()
		global.UI.Ok()
		global.UI.Say("Caching Service Plans... ")
		for _, service := range services {
			cf.Lookup.MarketplaceServicePlans(service)
		}
		global.UI.Ok()
		global.UI.Say("Caching Service Instances... ")
		cf.Lookup.ServiceInstances()
		global.UI.Ok()
		global.UI.Say("Caching Apps... ")
		cf.Lookup.Apps()
		global.UI.Ok()
		return nil
	},
}

func UpdateCommand() cli.Command {
	return updateCacheCommand
}
