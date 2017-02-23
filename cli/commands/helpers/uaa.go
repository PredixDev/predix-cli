package helpers

import (
	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

type UaaInterface interface {
	FetchInstance(instanceName string) (instance *cf.Item, serviceInfo map[string]interface{})
	AskForAdminClientSecret(c *cli.Context) string
}

type uaa struct{}

var Uaa UaaInterface = uaa{}

func (o uaa) FetchInstance(instanceName string) (instance *cf.Item, serviceInfo map[string]interface{}) {
	instance = cf.Lookup.PredixUaaInstanceItem(instanceName)
	if instance != nil {
		serviceInfo = ServiceInfo.FetchFor(instance)
	} else {
		global.UI.Failed("Predix UAA service instance %s not found", terminal.EntityNameColor(instanceName))
	}
	return instance, serviceInfo
}

func (o uaa) AskForAdminClientSecret(c *cli.Context) string {
	adminClientSecret := c.String("admin-secret")
	if adminClientSecret == "" {
		adminClientSecret = global.UI.AskForVerifiedPassword("Admin Client Secret")
	}
	return adminClientSecret
}
