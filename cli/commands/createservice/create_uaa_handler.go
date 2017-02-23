package createservice

import (
	"fmt"

	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

func createUaa(c *cli.Context) {
	args, plan := VerifyServicePlan(c)
	adminClientSecret := helpers.Uaa.AskForAdminClientSecret(c)

	SayCreatingServiceInstance(args)

	instance, err := cf.Curl.PostItem("/v2/service_instances?accepts_incomplete=true",
		fmt.Sprintf("{\"name\":\"%s\",\"space_guid\":\"%s\",\"service_plan_guid\":\"%s\",\"parameters\":{\"adminClientSecret\":\"%s\"}}",
			args[2], cf.CurrentUserInfo().SpaceGUID, plan.GUID, adminClientSecret))

	if err != nil {
		global.UI.Failed(err.Error())
	}
	global.UI.Ok()

	cf.Cache.InvalidateServiceInstances()
	cf.Lookup.ServiceInstances()

	global.UI.Say("")
	helpers.ServiceInfo.PrintFor(instance)
}

func init() {
	CreateServiceHandlers[constants.PredixUaa] = createUaa
}
