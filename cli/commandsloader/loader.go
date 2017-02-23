package commandsloader

import (
	"github.build.ge.com/adoption/predix-cli/cli/commands"
	"github.build.ge.com/adoption/predix-cli/cli/commands/cache"
	"github.build.ge.com/adoption/predix-cli/cli/commands/cf"
	"github.build.ge.com/adoption/predix-cli/cli/commands/createservice"
	"github.build.ge.com/adoption/predix-cli/cli/commands/serviceinfo"
	"github.build.ge.com/adoption/predix-cli/cli/commands/uaa"
)

func Load() {
	_ = commands.LoginCommand
	_ = cache.CacheCommand
	_ = cf.CfCommand
	_ = createservice.CreateServiceCommand
	_ = serviceinfo.ServiceInfoCommand
	_ = uaa.UaaCommand
}
