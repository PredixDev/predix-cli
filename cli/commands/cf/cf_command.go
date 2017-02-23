package cf

import (
	"fmt"
	"sort"
	"strings"

	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/commandregistry"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/urfave/cli"
)

type ParamLookupFunc func(map[string]string) []string

type CommandParameter struct {
	Name        string
	Description string
}

type CommandMetadata struct {
	Name        string
	ShortName   string
	Parameters  []string
	Arguments   func(int, []string) []string
	PostExec    func()
	ParamLookup map[string]ParamLookupFunc
	BoolParams  []string
}

func stringArrayContains(arr []string, val string) bool {
	if arr == nil {
		return false
	}
	for _, item := range arr {
		if strings.Compare(item, val) == 0 {
			return true
		}
	}
	return false
}

func findArgsAndParams(boolParams []string, args []string) (n int, unmatchedParam string, arguments []string, parameters map[string]string) {
	n = 0
	unmatchedParam = ""
	unmatchedArg := false

	arguments = []string{}
	parameters = map[string]string{}
	var bashCompletionFlag = fmt.Sprintf("--%s", cli.BashCompletionFlag.Name)
	for i := len(args) - 1; i >= 0; i-- {
		if args[i] == "--" || args[i] == bashCompletionFlag {
			continue
		}
		if strings.HasPrefix(args[i], "-") {
			if stringArrayContains(boolParams, args[i]) {
				unmatchedArg = false
				parameters[args[i]] = "true"
			} else if unmatchedArg {
				unmatchedArg = false
				parameters[args[i]] = arguments[len(arguments)-1]
				arguments = arguments[:len(arguments)-1]
				n--
			} else if unmatchedParam == "" {
				unmatchedParam = args[i]
			}
		} else {
			unmatchedArg = true
			arguments = append(arguments, args[i])
			n++
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(arguments)))
	return n, unmatchedParam, arguments, parameters
}

var cfCommands = []string{"-v", "-h", "--help"}
var cfCommandLookup = map[string]CommandMetadata{}

var CfCommand = cli.Command{
	Name:            "cf",
	Usage:           "Run Cloud Foundry CLI commands on the Predix Platform",
	Action:          Action,
	BashComplete:    cfBashComplete,
	SkipFlagParsing: true,
}

func Action(c *cli.Context) error {
	var bashCompletionFlag = fmt.Sprintf("--%s", cli.BashCompletionFlag.Name)
	var cfArgs = []interface{}{}
	args := c.Args()
	for i := 0; i < c.NArg(); i++ {
		if args[i] == "--" {
			continue
		}
		if args[i] == bashCompletionFlag {
			cfBashComplete(c)
			return nil
		}
		cfArgs = append(cfArgs, args[i])
	}
	var err = global.Sh.Interactive().Command("cf", cfArgs...).Run()
	if err != nil {
		return cli.NewExitError("", 1)
	}
	if c.NArg() > 0 {
		command := cfCommandLookup[c.Args()[0]]
		if command.PostExec != nil {
			command.PostExec()
		}
	}
	return nil
}

func cfBashComplete(c *cli.Context) {
	var bashCompletionFlag = fmt.Sprintf("--%s", cli.BashCompletionFlag.Name)
	if c.NArg() > 1 || (c.NArg() == 1 && c.Args()[0] != bashCompletionFlag) {
		printCompletionArgsForCommand(c)
	} else {
		for _, cmd := range cfCommands {
			global.UI.Say(cmd)
		}
	}
}

func printCompletionArgsForCommand(c *cli.Context) {
	command := cfCommandLookup[c.Args()[0]]
	noOfCommandArgs, unmatchedParam, arguments, parameters := findArgsAndParams(command.BoolParams, c.Args().Tail())

	var fetchedCompletionArgs = []string{}
	if unmatchedParam == "" {
		if command.Arguments != nil {
			fetchedCompletionArgs = command.Arguments(noOfCommandArgs, arguments)
		}
	} else if command.ParamLookup != nil {
		paramLookupFunc := command.ParamLookup[unmatchedParam]
		if paramLookupFunc != nil {
			fetchedCompletionArgs = paramLookupFunc(parameters)
		}
	}
	for _, arg := range fetchedCompletionArgs {
		global.UI.Say(arg)
	}
	if len(fetchedCompletionArgs) == 0 && unmatchedParam == "" {
		for _, param := range command.Parameters {
			if parameters[param] == "" {
				global.UI.Say(param)
			}
		}
	}
}

func cfLookupOrgs(map[string]string) []string {
	return cf.Lookup.Orgs()
}

func cfLookupSpaces(params map[string]string) []string {
	return cf.Lookup.Spaces(params)
}

func cfLookupStacks(map[string]string) []string {
	return cf.Lookup.Stacks()
}

func cfLookupMarketplaceServices(map[string]string) []string {
	return cf.Lookup.MarketplaceServices()
}

func init() {
	commandregistry.AddCommand(CfCommand)
}
