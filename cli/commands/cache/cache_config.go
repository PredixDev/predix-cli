package cache

import (
	"fmt"
	"strconv"
	"strings"

	"github.build.ge.com/adoption/predix-cli/cli/cache"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

var cacheConfigCommand = cli.Command{
	Name:  "config",
	Usage: "Get or set the cache config options",
	Subcommands: cli.Commands{
		{
			Name:      "all-timeouts",
			Usage:     "Get or set all the timeouts",
			ArgsUsage: "[timeout in seconds]",
			Before:    timeoutBeforeFunc,
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					timeout, _ := strconv.Atoi(c.Args()[0])
					timeouts := map[string]int{}
					for i := 0; i < len(cf.Caches); i++ {
						cf.Caches[i].Timeout = timeout
						timeouts[cf.Caches[i].Name] = timeout
					}
					cache.SaveTimeouts(timeouts)
				} else {
					for i := 0; i < len(cf.Caches); i++ {
						global.UI.Say("%-32s%d", fmt.Sprintf("%s-timeout", cf.Caches[i].Name), cf.Caches[i].Timeout)
					}
				}
				return nil
			},
		},
	},
}

func timeoutBeforeFunc(c *cli.Context) error {
	if c.NArg() == 1 {
		if _, err := strconv.Atoi(c.Args()[0]); err == nil {
			return nil
		}
		return cli.NewExitError("Invalid timeout!", 1)
	}
	return nil
}

func timeoutAction(c *cli.Context) error {
	cacheName := strings.Replace(c.Command.Name, "-timeout", "", -1)
	if c.NArg() == 1 {
		timeout, _ := strconv.Atoi(c.Args()[0])
		setTimeoutFor(cacheName, timeout)
	} else {
		for i := 0; i < len(cf.Caches); i++ {
			if strings.Compare(cacheName, cf.Caches[i].Name) == 0 {
				global.UI.Say(strconv.Itoa(cf.Caches[i].Timeout))
				break
			}
		}
	}
	return nil
}

func setTimeoutFor(cacheName string, timeout int) {
	cf.UpdateCacheTimeout(cacheName, timeout)
	cache.SaveTimeouts(map[string]int{
		cacheName: timeout,
	})
}

func ConfigCommand() cli.Command {
	for _, cache := range cf.Caches {
		command := cli.Command{
			Name:      fmt.Sprintf("%s-timeout", cache.Name),
			Usage:     fmt.Sprintf("Get or set the %s timeout", cache.Name),
			ArgsUsage: "[timeout in seconds]",
			Before:    timeoutBeforeFunc,
			Action:    timeoutAction,
		}
		cacheConfigCommand.Subcommands = append(cacheConfigCommand.Subcommands, command)
	}
	return cacheConfigCommand
}
