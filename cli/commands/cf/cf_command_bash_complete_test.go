package cf_test

import (
	"strings"

	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cf Command", func() {
	var (
		ui        *testterm.FakeUI
		sh        *testshell.FakeShell
		oldCurl   cf.CurlInterface
		oldLookup cf.LookupInterface
		curl      *cffakes.FakeCurlInterface
		lookup    *cffakes.FakeLookupInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		sh = &testshell.FakeShell{}
		curl = &cffakes.FakeCurlInterface{}
		lookup = &cffakes.FakeLookupInterface{}

		global.Env.NoCache = true
		global.UI = ui
		global.Sh = sh

		global.CurrentUserInfo = &global.UserInfo{
			API:       "FakeAPI",
			Name:      "FakeName",
			Org:       "FakeOrg",
			OrgGUID:   "FakeOrgGUID",
			OrgURL:    "FakeOrgURL",
			Space:     "FakeSpace",
			SpaceGUID: "FakeSpaceGUID",
			SpaceURL:  "FakeSpaceURL",
		}

		oldCurl = cf.Curl
		cf.Curl = curl

		oldLookup = cf.Lookup
		cf.Lookup = lookup
	})

	AfterEach(func() {
		cf.Lookup = oldLookup
		cf.Curl = oldCurl
	})

	Context("bash completion", func() {
		Describe("when action is called with --generate-bash-completion", func() {
			It("handles bash completion of cf command", func() {
				testcmd.RunCLICommand(CfCommand, []string{"--generate-bash-completion"}, ui)
				Expect(ui.Outputs).To(ConsistOf(cfCommandBashCompletions...))
			})

			Context("no data from cf", func() {
				var entries = []TableEntry{}
				for c, b := range cfCommandsParamsBashCompletion {
					entries = append(entries, Entry("shows completion for '"+c+" --generate-bash-completion' command", c, b))
				}
				DescribeTable("when a command is given", func(cmd string, params []interface{}) {
					testcmd.RunCLICommand(CfCommand, []string{cmd, "--generate-bash-completion"}, ui)
					Expect(ui.Outputs).To(ConsistOf(params...))
				}, entries...)

				entries = []TableEntry{}
				for c, b := range cfCommandsUnusedParamsBashCompletion {
					entries = append(entries, Entry("shows completion for parameters not used in '"+c+" --generate-bash-completion' command", c, b))
				}
				DescribeTable("when a command with parameters is given", func(cmd string, params []interface{}) {
					testcmd.RunCLICommand(CfCommand, append(strings.Split(cmd, " "), "--generate-bash-completion"), ui)
					Expect(ui.Outputs).To(ConsistOf(params...))
				}, entries...)
			})
		})

		It("displays all cf commands", func() {
			testcmd.BashCompleteCLICommand(CfCommand, []string{}, ui)
			Expect(ui.Outputs).To(ConsistOf(cfCommandBashCompletions...))
		})

		Context("no data from cf", func() {
			var entries = []TableEntry{}
			for c, b := range cfCommandsParamsBashCompletion {
				entries = append(entries, Entry("shows completion for '"+c+"' command", c, b))
			}
			DescribeTable("when a command is given", func(cmd string, params []interface{}) {
				testcmd.BashCompleteCLICommand(CfCommand, []string{cmd}, ui)
				Expect(ui.Outputs).To(ConsistOf(params...))
			}, entries...)

			entries = []TableEntry{}
			for c, b := range cfCommandsUnusedParamsBashCompletion {
				entries = append(entries, Entry("shows completion for parameters not used in '"+c+"' command", c, b))
			}
			DescribeTable("when a command with parameters is given", func(cmd string, params []interface{}) {
				testcmd.BashCompleteCLICommand(CfCommand, strings.Split(cmd, " "), ui)
				Expect(ui.Outputs).To(ConsistOf(params...))
			}, entries...)
		})

		Context("orgs present in cf", func() {
			BeforeEach(func() {
				lookup.OrgsReturns([]string{"Org1", "Org2"})
			})

			var entries = []TableEntry{}
			for _, c := range cfCommandsOrgsAndSpacesBashCompletion {
				entries = append(entries, Entry("shows orgs as completion for '"+c+" -o'", c))
			}
			DescribeTable("when orgs are expected for completion", func(cmd string) {
				testcmd.BashCompleteCLICommand(CfCommand, []string{cmd, "-o"}, ui)
				Expect(ui.Outputs).To(ConsistOf("Org1", "Org2"))
			}, entries...)
		})

		Context("spaces present in cf", func() {
			BeforeEach(func() {
				lookup.SpacesReturns([]string{"Space1", "Space2"})
			})

			var entries = []TableEntry{}
			for _, c := range cfCommandsOrgsAndSpacesBashCompletion {
				entries = append(entries, Entry("shows spaces from current org as completion for '"+c+" -s'", c))
			}
			DescribeTable("when no org specified and spaces are expected for completion", func(cmd string) {
				testcmd.BashCompleteCLICommand(CfCommand, []string{cmd, "-s"}, ui)
				Expect(ui.Outputs).To(ConsistOf("Space1", "Space2"))
			}, entries...)
		})

		Context("orgs and spaces present in cf", func() {
			BeforeEach(func() {
				lookup.OrgsReturns([]string{"Org1", "Org2"})
				lookup.SpacesReturns([]string{"Space1", "Space2"})
			})

			var entries = []TableEntry{}
			for _, c := range cfCommandsOrgsAndSpacesBashCompletion {
				entries = append(entries, Entry("shows spaces from specified org as completion for '"+c+" -o Org2 -s'", c))
			}
			DescribeTable("when org is specified and spaces are expected for completion", func(cmd string) {
				testcmd.BashCompleteCLICommand(CfCommand, []string{cmd, "-o", "Org2", "-s"}, ui)
				Expect(ui.Outputs).To(ConsistOf("Space1", "Space2"))
				Expect(lookup.SpacesArgsForCall(0)["-o"]).To(Equal("Org2"))
			}, entries...)
		})

		Context("apps present in cf", func() {
			BeforeEach(func() {
				lookup.AppsReturns([]string{"App1", "App2"})
			})

			var entries = []TableEntry{}
			for _, c := range cfCommandsAppsBashCompletion_NoArgs {
				entries = append(entries, Entry("shows apps as completion for '"+c+"'", c))
			}
			DescribeTable("when apps are expected for completion for first arg", func(cmd string) {
				testcmd.BashCompleteCLICommand(CfCommand, []string{cmd}, ui)
				Expect(ui.Outputs).To(ConsistOf("App1", "App2"))
			}, entries...)

			entries = []TableEntry{}
			for _, c := range cfCommandsAppsBashCompletion_WithArgs {
				entries = append(entries, Entry("shows apps as completion for '"+strings.Join(c, " ")+"'", c))
			}
			DescribeTable("when apps are expected for completion", func(cmd []string) {
				testcmd.BashCompleteCLICommand(CfCommand, cmd, ui)
				Expect(ui.Outputs).To(ConsistOf("App1", "App2"))
			}, entries...)
		})

		Context("service instances present in cf", func() {
			BeforeEach(func() {
				lookup.ServiceInstancesReturns([]string{"ServiceInstance1", "ServiceInstance2"})
			})

			var entries = []TableEntry{}
			for _, c := range cfCommandsServiceInstancesBashCompletion_NoArgs {
				entries = append(entries, Entry("shows service instances as completion for '"+c+"'", c))
			}
			DescribeTable("when service instances are expected for completion for first arg", func(cmd string) {
				testcmd.BashCompleteCLICommand(CfCommand, []string{cmd}, ui)
				Expect(ui.Outputs).To(ConsistOf("ServiceInstance1", "ServiceInstance2"))
			}, entries...)

			entries = []TableEntry{}
			for _, c := range cfCommandsServiceInstancesBashCompletion_WithArgs {
				entries = append(entries, Entry("shows service instances as completion for '"+strings.Join(c, " ")+"'", c))
			}
			DescribeTable("when service instances are expected for completion", func(cmd []string) {
				testcmd.BashCompleteCLICommand(CfCommand, cmd, ui)
				Expect(ui.Outputs).To(ConsistOf("ServiceInstance1", "ServiceInstance2"))
			}, entries...)
		})

		Context("services present in cf", func() {
			BeforeEach(func() {
				lookup.MarketplaceServicesReturns([]string{"Service1", "Service2"})
			})
			Describe("when bash completion required for 'marketplace -s' command", func() {
				It("shows services as completion", func() {
					testcmd.BashCompleteCLICommand(CfCommand, []string{"marketplace", "-s"}, ui)
					Expect(ui.Outputs).To(ConsistOf("Service1", "Service2"))
				})
			})
			Describe("when bash completion required for 'create-service' command", func() {
				It("shows services as completion", func() {
					testcmd.BashCompleteCLICommand(CfCommand, []string{"create-service"}, ui)
					Expect(ui.Outputs).To(ConsistOf("Service1", "Service2"))
				})
			})
		})

		Context("services and plans present in cf", func() {
			Describe("when bash completion required for 'create-service SERVICE' command", func() {
				It("shows plans as completion", func() {
					lookup.MarketplaceServicePlansReturns([]string{"Plan1", "Plan2"})

					testcmd.BashCompleteCLICommand(CfCommand, []string{"create-service", "SERVICE"}, ui)

					Expect(ui.Outputs).To(ConsistOf("Plan1", "Plan2"))
					Expect(lookup.MarketplaceServicePlansArgsForCall(0)).To(Equal("SERVICE"))
				})
				Context("when more than one service is found", func() {
					It("shows parameters as completion", func() {
						lookup.MarketplaceServicePlansReturns(nil)

						testcmd.BashCompleteCLICommand(CfCommand, []string{"create-service", "SERVICE"}, ui)

						Expect(ui.Outputs).To(ConsistOf("-c", "-t"))
					})
				})
				Context("when the specified service is not found and parameters are provided", func() {
					It("shows nothing as completion", func() {
						lookup.MarketplaceServicePlansReturns(nil)

						testcmd.BashCompleteCLICommand(CfCommand, []string{"create-service", "SERVICE", "-c", "'SOME_JSON'", "-t", "TAGS"}, ui)

						Expect(ui.Outputs).To(BeEmpty())
					})
				})
			})
		})
	})
})

var cfCommandBashCompletions = []interface{}{
	"-v",
	"-h",
	"--help",
	"apps",
	"a",
	"app",
	"push",
	"p",
	"scale",
	"delete",
	"d",
	"rename",
	"start",
	"st",
	"stop",
	"sp",
	"restart",
	"rs",
	"restage",
	"rg",
	"restart-app-instance",
	"events",
	"files",
	"f",
	"logs",
	"env",
	"e",
	"set-env",
	"se",
	"unset-env",
	"stacks",
	"stack",
	"copy-source",
	"create-app-manifest",
	"help",
	"version",
	"login",
	"l",
	"logout",
	"lo",
	"passwd",
	"pw",
	"target",
	"t",
	"api",
	"auth",
	"marketplace",
	"m",
	"services",
	"s",
	"service",
	"create-service",
	"cs",
	"update-service",
	"delete-service",
	"ds",
	"rename-service",
	"create-service-key",
	"csk",
	"service-keys",
	"sk",
	"service-key",
	"delete-service-key",
	"dsk",
	"bind-service",
	"bs",
	"unbind-service",
	"us",
	"create-user-provided-service",
	"cups",
	"update-user-provided-service",
	"uups",
}

var cfCommandsParamsBashCompletion = map[string][]interface{}{
	"login":  []interface{}{"-a", "-u", "-p", "-o", "-s", "--sso", "--skip-ssl-validation"},
	"logout": []interface{}{},
	"passwd": []interface{}{},
	"target": []interface{}{"-o", "-s"},
	"api":    []interface{}{"--unset", "--skip-ssl-validation"},
	"auth":   []interface{}{},

	"marketplace":                  []interface{}{"-s"},
	"services":                     []interface{}{},
	"service":                      []interface{}{"--guid"},
	"create-service":               []interface{}{"-c", "-t"},
	"update-service":               []interface{}{"-p", "-c", "-t"},
	"delete-service":               []interface{}{"-f"},
	"rename-service":               []interface{}{},
	"create-service-key":           []interface{}{"-c"},
	"service-keys":                 []interface{}{},
	"service-key":                  []interface{}{"--guid"},
	"delete-service-key":           []interface{}{"-f"},
	"bind-service":                 []interface{}{"-c"},
	"unbind-service":               []interface{}{},
	"create-user-provided-service": []interface{}{"-p", "-l"},
	"update-user-provided-service": []interface{}{"-p", "-l"},

	"apps":                 []interface{}{},
	"app":                  []interface{}{"--guid"},
	"push":                 []interface{}{"-b", "-c", "-d", "-f", "-i", "-k", "-m", "-n", "-p", "-s", "-t", "--no-hostname", "--no-manifest", "--no-route", "--no-start", "--random-route"},
	"scale":                []interface{}{"-i", "-k", "-m", "-f"},
	"delete":               []interface{}{"-f", "-r"},
	"rename":               []interface{}{},
	"start":                []interface{}{},
	"stop":                 []interface{}{},
	"restart":              []interface{}{},
	"restage":              []interface{}{},
	"restart-app-instance": []interface{}{},
	"events":               []interface{}{},
	"files":                []interface{}{"-i"},
	"logs":                 []interface{}{"--recent"},
	"env":                  []interface{}{},
	"set-env":              []interface{}{},
	"unset-env":            []interface{}{},
	"stacks":               []interface{}{},
	"stack":                []interface{}{"--guid"},
	"copy-source":          []interface{}{"-o", "-s", "--no-restart"},
	"create-app-manifest":  []interface{}{"-p"},
}

var cfCommandsUnusedParamsBashCompletion = map[string][]interface{}{
	"login -a API_URL":                  []interface{}{"-u", "-p", "-o", "-s", "--sso", "--skip-ssl-validation"},
	"login -a API_URL --sso":            []interface{}{"-u", "-p", "-o", "-s", "--skip-ssl-validation"},
	"target -s SPACE":                   []interface{}{"-o"},
	"api --unset --skip-ssl-validation": []interface{}{},
	"marketplace -s SERVICE":            []interface{}{},
	"service SERVICE_INSTANCE --guid":   []interface{}{},
	"create-service -c 'CONFIG_JSON'":   []interface{}{"-t"},
	"push --random-route -b BUILD_PACK": []interface{}{"-c", "-d", "-f", "-i", "-k", "-m", "-n", "-p", "-s", "-t", "--no-hostname", "--no-manifest", "--no-route", "--no-start"},
	"logs APP_NAME --recent":            []interface{}{},
	"copy-source -o ORG -s SPACE":       []interface{}{"--no-restart"},
}

var cfCommandsOrgsAndSpacesBashCompletion = []string{"login", "target", "copy-source"}

var cfCommandsAppsBashCompletion_NoArgs = []string{"app", "push", "scale", "delete", "rename", "start", "stop",
	"restart", "restage", "restart-app-instance", "events", "files", "logs", "env", "set-env", "unset-env",
	"copy-source", "create-app-manifest", "bind-service", "unbind-service"}

var cfCommandsAppsBashCompletion_WithArgs = [][]string{
	[]string{"copy-source", "DUMMY_ARG"},
}

var cfCommandsServiceInstancesBashCompletion_NoArgs = []string{"service", "update-service", "delete-service",
	"rename-service", "create-service-key", "service-keys", "service-key", "delete-service-key",
	"update-user-provided-service"}

var cfCommandsServiceInstancesBashCompletion_WithArgs = [][]string{
	[]string{"bind-service", "DUMMY_ARG"},
	[]string{"unbind-service", "DUMMY_ARG"},
}
