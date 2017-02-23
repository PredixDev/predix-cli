package helpers_test

import (
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers/helpersfakes"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
)

var _ = Describe("Uaa", func() {
	var (
		ui             *testterm.FakeUI
		oldLookup      cf.LookupInterface
		oldServiceInfo ServiceInfoInterface
		lookup         *cffakes.FakeLookupInterface
		serviceInfo    *helpersfakes.FakeServiceInfoInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		lookup = &cffakes.FakeLookupInterface{}
		serviceInfo = &helpersfakes.FakeServiceInfoInterface{}

		global.Env.NoCache = true
		global.UI = ui

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

		oldLookup = cf.Lookup
		cf.Lookup = lookup

		oldServiceInfo = ServiceInfo
		ServiceInfo = serviceInfo
	})

	AfterEach(func() {
		ServiceInfo = oldServiceInfo
		cf.Lookup = oldLookup
	})

	Context("fetch instance", func() {
		Describe("when instance found", func() {
			It("returns the instance and service info", func() {
				instance := &cf.Item{
					Name: "dummy-uaa",
				}
				lookup.PredixUaaInstanceItemReturns(instance)
				info := map[string]interface{}{
					"some-key": "some-value",
				}
				serviceInfo.FetchForReturns(info)
				var returnedInstance *cf.Item
				var returnedServiceInfo map[string]interface{}
				testcmd.Run(func() {
					returnedInstance, returnedServiceInfo = Uaa.FetchInstance("dummy-uaa")
				})

				Expect(ui.Outputs).To(BeNil())
				Expect(returnedInstance).To(Equal(instance))
				Expect(returnedServiceInfo).To(Equal(info))
				Expect(lookup.PredixUaaInstanceItemCallCount()).To(Equal(1))
				Expect(lookup.PredixUaaInstanceItemArgsForCall(0)).To(Equal("dummy-uaa"))
				Expect(serviceInfo.FetchForCallCount()).To(Equal(1))
				Expect(serviceInfo.FetchForArgsForCall(0)).To(Equal(instance))
			})
		})
		Describe("when instance not found", func() {
			It("shows an error", func() {
				lookup.PredixUaaInstanceItemReturns(nil)

				testcmd.Run(func() {
					Uaa.FetchInstance("dummy-uaa")
				})

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Predix UAA service instance dummy-uaa not found"))
				Expect(lookup.PredixUaaInstanceItemCallCount()).To(Equal(1))
				Expect(lookup.PredixUaaInstanceItemArgsForCall(0)).To(Equal("dummy-uaa"))
				Expect(serviceInfo.FetchForCallCount()).To(Equal(0))
			})
		})
	})

	Context("ask for admin client", func() {
		Describe("when given in args", func() {
			It("returns the value", func() {
				var adminSecret string
				testcmd.Run(func() {
					adminSecret = Uaa.AskForAdminClientSecret(testcmd.Context(cli.Command{
						Flags: []cli.Flag{
							cli.StringFlag{
								Name: "admin-secret",
							},
						},
					}, []string{"--admin-secret", "secret"}))
				})
				Expect(adminSecret).To(Equal("secret"))
			})
		})
		Describe("when not given in args", func() {
			It("prompts user and returns the verified value", func() {
				ui.Inputs = []string{"secret", "secret"}

				var adminSecret string
				testcmd.Run(func() {
					adminSecret = Uaa.AskForAdminClientSecret(testcmd.Context(cli.Command{
						Flags: []cli.Flag{
							cli.StringFlag{
								Name: "admin-secret",
							},
						},
					}, nil))
				})

				Expect(adminSecret).To(Equal("secret"))
				Expect(ui.PasswordPrompts).To(ConsistOf("Admin Client Secret", "Verify Admin Client Secret"))
			})
		})
	})
})
