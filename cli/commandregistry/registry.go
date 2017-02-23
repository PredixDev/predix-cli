package commandregistry

import "github.com/urfave/cli"

type registry struct {
	Commands []cli.Command
}

var Registry = NewRegistry()

func NewRegistry() *registry {
	return &registry{
		Commands: []cli.Command{},
	}
}

func GetCommands() []cli.Command {
	return Registry.Commands
}

func AddCommand(cmd cli.Command) {
	Registry.Commands = append(Registry.Commands, cmd)
}
