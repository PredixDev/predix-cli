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

var _ = Describe("Target", func() {
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

	Context("before", func() {
		Describe("when called with 0 args", func() {
			It("returns nil", func() {
				Expect(testcmd.BeforeCLICommand(TargetCommand, nil, ui)).To(BeNil())
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when called with 1 arg", func() {
			It("returns nil", func() {
				Expect(testcmd.BeforeCLICommand(TargetCommand, []string{"1"}, ui)).To(BeNil())
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when called with more than 1 arg", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(TargetCommand, []string{"arg1", "arg2"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})
	})

	Context("action", func() {
		Describe("when called with no args", func() {
			It("prints current target", func() {
				testcmd.RunCLICommand(TargetCommand, nil, ui)
				Expect(targets.PrintCurrentCallCount()).To(Equal(1))
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when called with 1 args", func() {
			It("sets the target and prints it", func() {
				testcmd.RunCLICommand(TargetCommand, []string{"1"}, ui)
				Expect(targets.SetCurrentForIDCallCount()).To(Equal(1))
				Expect(targets.SetCurrentForIDArgsForCall(0)).To(Equal(1))
				Expect(targets.PrintCurrentCallCount()).To(Equal(1))
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when called with an invalid args", func() {
			It("show an error", func() {
				testcmd.RunCLICommand(TargetCommand, []string{"Invalid"}, ui)
				Expect(targets.PrintCurrentCallCount()).To(Equal(0))
				Expect(ui.Outputs).To(ConsistOf("FAILED", "Invalid ID."))
			})
		})
	})
})
