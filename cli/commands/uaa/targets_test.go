package uaa_test

import (
	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/uaa"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.build.ge.com/adoption/predix-cli/cli/uaac/uaacfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Targets", func() {
	var (
		ui         *testterm.FakeUI
		sh         *testshell.FakeShell
		oldTargets uaac.TargetsInterface
		targets    *uaacfakes.FakeTargetsInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		sh = &testshell.FakeShell{}
		targets = &uaacfakes.FakeTargetsInterface{}

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

		oldTargets = uaac.Targets
		uaac.Targets = targets
	})

	AfterEach(func() {
		uaac.Targets = oldTargets
	})

	Context("action", func() {
		Describe("when called with less than 1 arg", func() {
			It("returns error", func() {
				testcmd.RunCLICommand(TargetsCommand, nil, ui)
				Expect(ui.Outputs).To(BeNil())
				Expect(targets.PrintAllCallCount()).To(Equal(1))
			})
		})
	})
})
