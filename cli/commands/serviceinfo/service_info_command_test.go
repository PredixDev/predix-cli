package serviceinfo_test

import (
	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers/helpersfakes"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/serviceinfo"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Info", func() {
	var (
		ui             *testterm.FakeUI
		sh             *testshell.FakeShell
		oldCurl        cf.CurlInterface
		oldLookup      cf.LookupInterface
		oldServiceInfo helpers.ServiceInfoInterface
		curl           *cffakes.FakeCurlInterface
		lookup         *cffakes.FakeLookupInterface
		serviceInfo    *helpersfakes.FakeServiceInfoInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		sh = &testshell.FakeShell{}
		curl = &cffakes.FakeCurlInterface{}
		lookup = &cffakes.FakeLookupInterface{}
		serviceInfo = &helpersfakes.FakeServiceInfoInterface{}

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

		oldServiceInfo = helpers.ServiceInfo
		helpers.ServiceInfo = serviceInfo
	})

	AfterEach(func() {
		helpers.ServiceInfo = oldServiceInfo
		cf.Lookup = oldLookup
		cf.Curl = oldCurl
	})

	Context("before", func() {
		Describe("when called with less than 1 arg", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(ServiceInfoCommand, nil, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})
		Describe("when called with less than 2 arg", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(ServiceInfoCommand, []string{"arg1"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})
		Describe("when called with more than 2 args", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(ServiceInfoCommand, []string{"arg1", "arg2", "arg3"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})
		Describe("when called with 2 args", func() {
			It("returns nil", func() {
				Expect(testcmd.BeforeCLICommand(ServiceInfoCommand, []string{"arg1", "arg2"}, ui)).To(BeNil())
				Expect(ui.Outputs).To(BeNil())
			})
		})
	})

	Context("bash complete", func() {
		Describe("when called with 0 args", func() {
			It("shows apps", func() {
				lookup.AppsReturns([]string{"App1", "App2"})

				testcmd.BashCompleteCLICommand(ServiceInfoCommand, nil, ui)

				Expect(ui.Outputs).To(ConsistOf("App1", "App2"))
			})
		})

		Describe("when called with 1 arg", func() {
			It("shows service instances", func() {
				lookup.ServiceInstancesReturns([]string{"ServiceInstance1", "ServiceInstance2"})

				testcmd.BashCompleteCLICommand(ServiceInfoCommand, []string{"arg1"}, ui)

				Expect(ui.Outputs).To(ConsistOf("ServiceInstance1", "ServiceInstance2"))
			})
		})

		Describe("when called with 2 args", func() {
			It("shows nothing", func() {
				testcmd.BashCompleteCLICommand(ServiceInfoCommand, []string{"app", "service-instance"}, ui)

				Expect(ui.Outputs).To(BeNil())
			})
		})
	})

	Context("action", func() {
		Describe("when successfully looked up info", func() {
			It("shows the info", func() {
				app := &cf.Item{
					Name: "dummy-app",
				}
				lookup.AppForNameReturns(app)
				instance := &cf.Item{
					Name: "dummy-service",
				}
				lookup.InstanceForInstanceNameReturns(instance)

				testcmd.RunCLICommand(ServiceInfoCommand, []string{"dummy-app", "dummy-service"}, ui)

				Expect(ui.Outputs).To(BeNil())
				Expect(lookup.AppForNameCallCount()).To(Equal(1))
				Expect(lookup.AppForNameArgsForCall(0)).To(Equal("dummy-app"))
				Expect(lookup.InstanceForInstanceNameCallCount()).To(Equal(1))
				Expect(lookup.InstanceForInstanceNameArgsForCall(0)).To(Equal("dummy-service"))
				Expect(serviceInfo.PrintForAppAndServiceInstanceCallCount()).To(Equal(1))
				calledApp, calledInstance := serviceInfo.PrintForAppAndServiceInstanceArgsForCall(0)
				Expect(calledApp).To(Equal(app))
				Expect(calledInstance).To(Equal(instance))
			})
		})
		Describe("when app not found", func() {
			It("shows an error", func() {
				instance := &cf.Item{
					Name: "dummy-service",
				}
				lookup.InstanceForInstanceNameReturns(instance)
				testcmd.RunCLICommand(ServiceInfoCommand, []string{"dummy-app", "dummy-service"}, ui)

				Expect(ui.Outputs).To(ConsistOf("FAILED", "App dummy-app not found"))
			})
		})
		Describe("when service instance not found", func() {
			It("shows an error", func() {
				app := &cf.Item{
					Name: "dummy-app",
				}
				lookup.AppForNameReturns(app)
				testcmd.RunCLICommand(ServiceInfoCommand, []string{"dummy-app", "dummy-service"}, ui)

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Service instance dummy-service not found"))
			})
		})
	})
})
