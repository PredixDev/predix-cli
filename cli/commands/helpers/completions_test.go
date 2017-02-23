package helpers_test

import (
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
)

var _ = Describe("Completions", func() {
	var (
		ui *testterm.FakeUI
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		global.Env.NoCache = true
		global.UI = ui
	})

	It("prints unused flags", func() {
		cmd := cli.Command{
			Name:  "cf",
			Usage: "Run Cloud Foundry CLI commands on the Predix Platform",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "user,u",
				},
				cli.StringFlag{
					Name: "password, p",
				},
				cli.StringFlag{
					Name: "secret,s",
				},
				cli.BoolFlag{
					Name: "v",
				},
				cli.BoolFlag{
					Name: "L",
				},
				cli.IntFlag{
					Name: "timeout",
				},
			},
		}

		Completions.PrintFlags(testcmd.Context(cmd, []string{"--user", "abc", "-s", "xyz", "-L"}), cmd.Flags)

		Expect(ui.Outputs).To(ConsistOf("--password", "-p", "-v", "--timeout"))
	})
})
