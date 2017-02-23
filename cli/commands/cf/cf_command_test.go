package cf_test

import (
	"errors"
	"strings"

	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cf Command", func() {
	var (
		ui *testterm.FakeUI
		sh *testshell.FakeShell
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		sh = &testshell.FakeShell{}

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
	})

	Context("command action", func() {
		It("passes arguments to cf cli", func() {
			sh.Runs = []error{nil}
			args := strings.Split("ARG1 ARG2 ARG3 --FLAG1 --FLAG2 FLAG2_VALUE -FLAG3 -FLAG4 FLAG4_VALUE", " ")

			testcmd.RunCLICommand(CfCommand, args, ui)

			Expect(sh.Commands).To(ConsistOf(append([]string{"interactive", "cf"}, args...)))
		})

		Describe("when cf cli returns an error", func() {
			It("returns and error", func() {
				sh.Runs = []error{errors.New("CLI Error")}
				args := strings.Split("ARG1 ARG2 ARG3 --FLAG1 --FLAG2 FLAG2_VALUE -FLAG3 -FLAG4 FLAG4_VALUE", " ")

				testcmd.RunCLICommand(CfCommand, args, ui)

				Expect(sh.Commands).To(ConsistOf(append([]string{"interactive", "cf"}, args...)))
				Expect(ui.Outputs).To(ConsistOf("FAILED", ""))
			})
		})
	})
})
