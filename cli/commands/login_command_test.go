package commands_test

import (
	"errors"

	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	. "github.build.ge.com/adoption/predix-cli/cli/commands"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {
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
	})

	Context("interactive usage", func() {
		Context("user chooses us-west", func() {
			It("sets api for us-west PoP and runs login", func() {
				ui.Inputs = []string{"1"}
				sh.Runs = []error{nil}

				testcmd.RunCLICommand(LoginCommand, []string{}, ui)

				Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK"))
				Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
				Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-usw02-pr.ice.predix.io"))
			})

			Describe("when logging in fails", func() {
				It("shows an error", func() {
					ui.Inputs = []string{"1"}
					sh.Runs = []error{errors.New("Login error")}

					testcmd.RunCLICommand(LoginCommand, []string{}, ui)

					Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK", "FAILED", "Error while logging user in"))
					Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
					Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-usw02-pr.ice.predix.io"))
				})
			})

			Context("user provides login flags", func() {
				It("passes the login flags to login", func() {
					ui.Inputs = []string{"1"}
					sh.Runs = []error{nil}

					testcmd.RunCLICommand(LoginCommand, []string{"--sso", "-u", "testUser", "-p", "testPass"}, ui)

					Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK"))
					Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
					Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-usw02-pr.ice.predix.io", "--sso", "-u", "testUser", "-p", "testPass"))
				})
			})
		})

		Context("user chooses us-east", func() {
			It("sets api for us-east PoP and runs login", func() {
				ui.Inputs = []string{"2"}
				sh.Runs = []error{nil}

				testcmd.RunCLICommand(LoginCommand, []string{}, ui)

				Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK"))
				Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
				Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.asv-pr.ice.predix.io"))
			})
		})

		Context("user chooses japan", func() {
			It("sets api for us-west PoP and runs login", func() {
				ui.Inputs = []string{"3"}
				sh.Runs = []error{nil}

				testcmd.RunCLICommand(LoginCommand, []string{}, ui)

				Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK"))
				Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
				Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-jp01-pr.ice.predix.io"))
			})
		})

		Context("user chooses uk", func() {
			It("sets api for us-west PoP and runs login", func() {
				ui.Inputs = []string{"4"}
				sh.Runs = []error{nil}

				testcmd.RunCLICommand(LoginCommand, []string{}, ui)

				Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK"))
				Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
				Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.dc-uk01-pr.ice.predix.io"))
			})
		})

		Context("user gives incorrect choice", func() {
			It("shows error and exits", func() {
				ui.Inputs = []string{"10"}
				sh.Runs = []error{nil}

				testcmd.RunCLICommand(LoginCommand, []string{}, ui)

				Expect(ui.Outputs).To(ConsistOf("1. US-West", "2. US-East", "3. Japan", "4. UK", "FAILED", "Invalid choice!"))
				Expect(ui.Prompts).To(ConsistOf("Choose the PoP to set"))
				Expect(sh.Commands).To(BeEmpty())
			})
		})
	})

	Context("when the user provides the --us-west flag", func() {
		It("does not prompt and sets the api", func() {
			sh.Runs = []error{nil}

			testcmd.RunCLICommand(LoginCommand, []string{"--us-west"}, ui)

			Expect(ui.Outputs).To(BeEmpty())
			Expect(ui.Prompts).To(BeEmpty())
			Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-usw02-pr.ice.predix.io"))
		})

		Context("user provides login flags", func() {
			It("passes the login flags to login", func() {
				sh.Runs = []error{nil}

				testcmd.RunCLICommand(LoginCommand, []string{"--us-west", "--skip-ssl-validation", "-o", "someOrg", "-s", "someSpace"}, ui)

				Expect(ui.Outputs).To(BeEmpty())
				Expect(ui.Prompts).To(BeEmpty())
				Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-usw02-pr.ice.predix.io", "-o", "someOrg", "-s", "someSpace", "--skip-ssl-validation"))
			})
		})

		Describe("when logging in fails", func() {
			It("shows an error", func() {
				sh.Runs = []error{errors.New("Login error")}

				testcmd.RunCLICommand(LoginCommand, []string{"--us-west", "-u", "testUser", "-p", "testPass"}, ui)

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Error while logging user in"))
				Expect(ui.Prompts).To(BeEmpty())
				Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-usw02-pr.ice.predix.io", "-u", "testUser", "-p", "testPass"))
			})
		})
	})

	Context("when the user provides the --us-east flag", func() {
		It("does not prompt and sets the api", func() {
			sh.Runs = []error{nil}

			testcmd.RunCLICommand(LoginCommand, []string{"--us-east"}, ui)

			Expect(ui.Outputs).To(BeEmpty())
			Expect(ui.Prompts).To(BeEmpty())
			Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.asv-pr.ice.predix.io"))
		})
	})

	Context("when the user provides the --japan flag", func() {
		It("does not prompt and sets the api", func() {
			sh.Runs = []error{nil}

			testcmd.RunCLICommand(LoginCommand, []string{"--japan"}, ui)

			Expect(ui.Outputs).To(BeEmpty())
			Expect(ui.Prompts).To(BeEmpty())
			Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.aws-jp01-pr.ice.predix.io"))
		})
	})

	Context("when the user provides the --uk flag", func() {
		It("does not prompt and sets the api", func() {
			sh.Runs = []error{nil}

			testcmd.RunCLICommand(LoginCommand, []string{"--uk"}, ui)

			Expect(ui.Outputs).To(BeEmpty())
			Expect(ui.Prompts).To(BeEmpty())
			Expect(sh.Commands).To(ConsistOf("interactive", "cf", "login", "-a", "https://api.system.dc-uk01-pr.ice.predix.io"))
		})
	})
})
