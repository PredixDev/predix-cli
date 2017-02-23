package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.com/urfave/cli"
)

func RunCLICommand(cmd cli.Command, args []string, ui *testterm.FakeUI) {
	set := flagSet(cmd, args)

	defer func() {
		errMsg := recover()
		if errMsg != nil && errMsg != testterm.QuietPanic {
			panic(errMsg)
		}
	}()

	context := cli.NewContext(nil, set, nil)
	vals := reflect.ValueOf(cmd.Action).Call([]reflect.Value{reflect.ValueOf(context)})
	err, _ := vals[0].Interface().(error)
	if err != nil {
		ui.Failed(err.Error())
	}
}

func Run(function func()) {
	defer func() {
		errMsg := recover()
		if errMsg != nil && errMsg != testterm.QuietPanic {
			panic(errMsg)
		}
	}()

	function()
}

func BashCompleteCLICommand(cmd cli.Command, args []string, ui *testterm.FakeUI) {
	set := flagSet(cmd, args)
	context := cli.NewContext(nil, set, nil)
	cmd.BashComplete(context)
}

func BeforeCLICommand(cmd cli.Command, args []string, ui *testterm.FakeUI) error {
	set := flagSet(cmd, args)
	context := cli.NewContext(nil, set, nil)

	defer func() {
		errMsg := recover()
		if errMsg != nil && errMsg != testterm.QuietPanic {
			panic(errMsg)
		}
	}()

	return cmd.Before(context)
}

func Context(cmd cli.Command, args []string) *cli.Context {
	set := flagSet(cmd, args)
	return cli.NewContext(nil, set, nil)
}

func flagSet(cmd cli.Command, args []string) *flag.FlagSet {
	if cmd.SkipFlagParsing {
		index := len(args)
		for i, v := range args {
			if strings.HasPrefix(v, "-") {
				index = i
				break
			}
		}
		newArgs := append([]string{}, args[:index]...)
		newArgs = append(newArgs, "--")
		newArgs = append(newArgs, args[index:]...)
		args = newArgs
	}
	set := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	for _, f := range cmd.Flags {
		f.Apply(set)
	}
	err := set.Parse(args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	err = normalizeFlags(cmd.Flags, set)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return set
}

func copyFlag(name string, ff *flag.Flag, set *flag.FlagSet) {
	switch ff.Value.(type) {
	case *cli.StringSlice:
	default:
		_ = set.Set(name, ff.Value.String())
	}
}

func normalizeFlags(flags []cli.Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := strings.Split(f.GetName(), ",")
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if !visited[name] {
				copyFlag(name, ff, set)
			}
		}
	}
	return nil
}
