package helpers

import (
	"strings"

	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

type CompletionsInterface interface {
	PrintFlags(c *cli.Context, flags []cli.Flag)
}

type completions struct{}

var Completions CompletionsInterface = completions{}

func (o completions) PrintFlags(c *cli.Context, flags []cli.Flag) {
	if flags != nil {
		for _, flag := range flags {
			name := strings.Split(flag.GetName(), ",")[0]
			if !c.IsSet(name) {
				names := strings.Split(flag.GetName(), ",")
				for _, name := range names {
					global.UI.Say(prefixed(strings.TrimSpace(name)))
				}
			}
		}
	}
}

func prefixed(name string) string {
	return prefixFor(name) + name
}

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}
	return prefix
}
