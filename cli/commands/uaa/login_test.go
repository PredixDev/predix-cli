package uaa_test

import (
	"fmt"

	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers/helpersfakes"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/uaa"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.build.ge.com/adoption/predix-cli/cli/uaac/uaacfakes"
	lib "github.com/PredixDev/go-uaa-lib"
	"github.com/PredixDev/go-uaa-lib/libfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {
	var (
		ui        *testterm.FakeUI
		sh        *testshell.FakeShell
		oldLookup cf.LookupInterface
		lookup    *cffakes.FakeLookupInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		sh = &testshell.FakeShell{}
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

		oldLookup = cf.Lookup
		cf.Lookup = lookup
	})
	AfterEach(func() {
		cf.Lookup = oldLookup
	})

	Context("before", func() {
		Describe("when called with less than 2 args", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(LoginCommand, nil, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())

				err = testcmd.BeforeCLICommand(LoginCommand, []string{"arg1"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when called with 2 args", func() {
			It("returns nil", func() {
				err := testcmd.BeforeCLICommand(LoginCommand, []string{"arg1", "arg2"}, ui)
				Expect(err).To(BeNil())
			})
		})

		Describe("when called with 3 args", func() {
			It("returns nil", func() {
				err := testcmd.BeforeCLICommand(LoginCommand, []string{"arg1", "arg2", "arg3"}, ui)
				Expect(err).To(BeNil())
			})
		})

		Describe("when called with more than 3 args", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(LoginCommand, []string{"arg1", "arg2", "arg3", "arg4"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})
	})

	Context("bash complete", func() {
		Describe("when called with 0 args", func() {
			It("shows predix uaa instances", func() {
				lookup.PredixUaaInstancesReturns([]string{"predix-uaa-1", "predix-uaa-2"})
				testcmd.BashCompleteCLICommand(LoginCommand, nil, ui)
				Expect(ui.Outputs).To(ConsistOf("predix-uaa-1", "predix-uaa-2"))
			})
		})

		Describe("when called with 1 arg", func() {
			It("shows suggestion to enter client id", func() {
				testcmd.BashCompleteCLICommand(LoginCommand, []string{"some-uaa"}, ui)
				Expect(ui.Outputs).To(ConsistOf("Enter: CLIENT_ID", "_"))
			})
		})

		Describe("when called with 0 args and no predix uaa instances to lookup", func() {
			It("shows parameters as completion", func() {
				lookup.PredixUaaInstancesReturns(nil)
				testcmd.BashCompleteCLICommand(LoginCommand, nil, ui)
				Expect(ui.Outputs).To(ConsistOf("--ca-cert", "--skip-ssl-validation", "--secret", "-s", "--scope", "--password", "-p"))
			})
		})

		Describe("when called with 2 args", func() {
			It("shows username and parameters as completion", func() {
				testcmd.BashCompleteCLICommand(LoginCommand, []string{"some-uaa", "some-client"}, ui)
				Expect(ui.Outputs).To(ConsistOf("[USERNAME]", "--ca-cert", "--skip-ssl-validation", "--secret", "-s", "--scope", "--password", "-p"))
			})
		})
	})

	Context("action", func() {
		var (
			oldUaa                helpers.UaaInterface
			oldServiceInfo        helpers.ServiceInfoInterface
			oldTargets            uaac.TargetsInterface
			oldInfoClientFactory  lib.InfoClientFactoryInterface
			oldTokenIssuerFactory lib.TokenIssuerFactoryInterface
			uaa                   *helpersfakes.FakeUaaInterface
			serviceInfo           *helpersfakes.FakeServiceInfoInterface
			targets               *uaacfakes.FakeTargetsInterface
			infoClient            *libfakes.FakeInfoClient
			infoClientFactory     *libfakes.FakeInfoClientFactoryInterface
			tokenIssuer           *libfakes.FakeTokenIssuer
			tokenIssuerFactory    *libfakes.FakeTokenIssuerFactoryInterface
		)

		BeforeEach(func() {
			uaa = &helpersfakes.FakeUaaInterface{}
			serviceInfo = &helpersfakes.FakeServiceInfoInterface{}
			targets = &uaacfakes.FakeTargetsInterface{}
			infoClient = &libfakes.FakeInfoClient{}
			infoClientFactory = &libfakes.FakeInfoClientFactoryInterface{}
			tokenIssuer = &libfakes.FakeTokenIssuer{}
			tokenIssuerFactory = &libfakes.FakeTokenIssuerFactoryInterface{}

			infoClientFactory.NewReturns(infoClient)
			tokenIssuerFactory.NewReturns(tokenIssuer)

			oldUaa = helpers.Uaa
			helpers.Uaa = uaa

			oldServiceInfo = helpers.ServiceInfo
			helpers.ServiceInfo = serviceInfo

			oldTargets = uaac.Targets
			uaac.Targets = targets

			oldInfoClientFactory = lib.InfoClientFactory
			lib.InfoClientFactory = infoClientFactory

			oldTokenIssuerFactory = lib.TokenIssuerFactory
			lib.TokenIssuerFactory = tokenIssuerFactory
		})

		AfterEach(func() {
			lib.TokenIssuerFactory = oldTokenIssuerFactory
			lib.InfoClientFactory = oldInfoClientFactory
			uaac.Targets = oldTargets
			helpers.ServiceInfo = oldServiceInfo
			helpers.Uaa = oldUaa
		})

		Describe("when a valid uaa, client id and client secret is provided", func() {
			It("does client credentials grant", func() {
				tr := &lib.TokenResponse{}

				ui.Inputs = []string{"client-secret"}
				uaa.FetchInstanceReturns(&cf.Item{URL: "dummy-instance-url"}, map[string]interface{}{"uri": "dummy-uaa-url"})
				serviceInfo.ResolveJSONPathReturns("dummy-uaa-url")
				infoClient.ServerReturns(nil)
				tokenIssuer.ClientCredentialsGrantReturns(tr, nil)

				testcmd.RunCLICommand(LoginCommand, []string{"--scope", "scope1,scope2", "--skip-ssl-validation", "--ca-cert", "dummy-cert-file", "dummy-uaa", "client-id"}, ui)

				Expect(ui.Outputs).To(BeNil())
				Expect(uaa.FetchInstanceCallCount()).To(Equal(1))
				Expect(serviceInfo.ResolveJSONPathCallCount()).To(Equal(1))
				Expect(infoClientFactory.NewCallCount()).To(Equal(1))
				calledUrl, _, _ := infoClientFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(ui.PasswordPrompts).To(ConsistOf("Client Secret"))
				Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
				calledUrl, calledClientID, calledClientSecret, calledSslVerify, calledCaCert := tokenIssuerFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(calledClientID).To(Equal("client-id"))
				Expect(calledClientSecret).To(Equal("client-secret"))
				Expect(calledSslVerify).To(BeTrue())
				Expect(calledCaCert).To(Equal("dummy-cert-file"))
				Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(1))
				Expect(tokenIssuer.ClientCredentialsGrantArgsForCall(0)).To(ConsistOf("scope1", "scope2"))
				Expect(tokenIssuer.PasswordGrantCallCount()).To(Equal(0))
				Expect(targets.SetCurrentCallCount()).To(Equal(1))
				calledUrl, calledInstanceUrl, calledSslVerify, calledCaCert, calledTokenResponse := targets.SetCurrentArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(calledInstanceUrl).To(Equal("dummy-instance-url"))
				Expect(calledSslVerify).To(BeTrue())
				Expect(calledCaCert).To(Equal("dummy-cert-file"))
				Expect(calledTokenResponse).To(Equal(tr))
				Expect(targets.PrintCurrentCallCount()).To(Equal(1))
			})
		})

		Describe("when a valid uaa, client id, client secret, username and password is provided", func() {
			It("does password grant", func() {
				tr := &lib.TokenResponse{}

				ui.Inputs = []string{"user-password"}
				uaa.FetchInstanceReturns(&cf.Item{URL: "dummy-instance-url"}, map[string]interface{}{"uri": "dummy-uaa-url"})
				serviceInfo.ResolveJSONPathReturns("dummy-uaa-url")
				infoClient.ServerReturns(nil)
				tokenIssuer.PasswordGrantReturns(tr, nil)

				testcmd.RunCLICommand(LoginCommand, []string{"--scope", "scope1,scope2", "--secret", "client-secret", "dummy-uaa", "client-id", "some-user"}, ui)

				Expect(ui.Outputs).To(BeNil())
				Expect(uaa.FetchInstanceCallCount()).To(Equal(1))
				Expect(serviceInfo.ResolveJSONPathCallCount()).To(Equal(1))
				Expect(infoClientFactory.NewCallCount()).To(Equal(1))
				calledUrl, _, _ := infoClientFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(ui.PasswordPrompts).To(ConsistOf("User Password"))
				Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
				calledUrl, calledClientID, calledClientSecret, calledSslVerify, calledCaCert := tokenIssuerFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(calledClientID).To(Equal("client-id"))
				Expect(calledClientSecret).To(Equal("client-secret"))
				Expect(calledSslVerify).To(BeFalse())
				Expect(calledCaCert).To(Equal(""))
				Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(0))
				Expect(tokenIssuer.PasswordGrantCallCount()).To(Equal(1))
				calledUsername, calledPassword, calledScopes := tokenIssuer.PasswordGrantArgsForCall(0)
				Expect(calledUsername).To(Equal("some-user"))
				Expect(calledPassword).To(Equal("user-password"))
				Expect(calledScopes).To(ConsistOf("scope1", "scope2"))
				Expect(targets.SetCurrentCallCount()).To(Equal(1))
				calledUrl, calledInstanceUrl, calledSslVerify, calledCaCert, calledTokenResponse := targets.SetCurrentArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(calledInstanceUrl).To(Equal("dummy-instance-url"))
				Expect(calledSslVerify).To(BeFalse())
				Expect(calledCaCert).To(Equal(""))
				Expect(calledTokenResponse).To(Equal(tr))
				Expect(targets.PrintCurrentCallCount()).To(Equal(1))
			})
		})

		Describe("when an incorrect uaa is provided", func() {
			It("shows an error", func() {
				uaa.FetchInstanceReturns(&cf.Item{URL: "dummy-instance-url"}, map[string]interface{}{"uri": "dummy-uaa-url"})
				serviceInfo.ResolveJSONPathReturns("dummy-uaa-url")
				infoClient.ServerReturns(fmt.Errorf("Invalid UAA"))

				testcmd.RunCLICommand(LoginCommand, []string{"dummy-uaa", "client-id"}, ui)

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Invalid UAA"))
				Expect(uaa.FetchInstanceCallCount()).To(Equal(1))
				Expect(serviceInfo.ResolveJSONPathCallCount()).To(Equal(1))
				Expect(infoClientFactory.NewCallCount()).To(Equal(1))
				calledUrl, _, _ := infoClientFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(ui.PasswordPrompts).To(BeNil())
				Expect(tokenIssuerFactory.NewCallCount()).To(Equal(0))
				Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(0))
				Expect(tokenIssuer.PasswordGrantCallCount()).To(Equal(0))
				Expect(targets.SetCurrentCallCount()).To(Equal(0))
				Expect(targets.PrintCurrentCallCount()).To(Equal(0))
			})
		})

		Describe("when a valid uaa is provided and client credentials grant fails", func() {
			It("shows an error", func() {
				uaa.FetchInstanceReturns(&cf.Item{URL: "dummy-instance-url"}, map[string]interface{}{"uri": "dummy-uaa-url"})
				serviceInfo.ResolveJSONPathReturns("dummy-uaa-url")
				infoClient.ServerReturns(nil)
				tokenIssuer.ClientCredentialsGrantReturns(nil, fmt.Errorf("Unauthorized"))

				testcmd.RunCLICommand(LoginCommand, []string{"--secret", "client-secret", "dummy-uaa", "client-id"}, ui)

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Unauthorized"))
				Expect(uaa.FetchInstanceCallCount()).To(Equal(1))
				Expect(serviceInfo.ResolveJSONPathCallCount()).To(Equal(1))
				Expect(infoClientFactory.NewCallCount()).To(Equal(1))
				calledUrl, _, _ := infoClientFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(ui.PasswordPrompts).To(BeNil())
				Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
				calledUrl, calledClientID, calledClientSecret, calledSslVerify, calledCaCert := tokenIssuerFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(calledClientID).To(Equal("client-id"))
				Expect(calledClientSecret).To(Equal("client-secret"))
				Expect(calledSslVerify).To(BeFalse())
				Expect(calledCaCert).To(Equal(""))
				Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(1))
				Expect(tokenIssuer.ClientCredentialsGrantArgsForCall(0)).To(BeNil())
				Expect(tokenIssuer.PasswordGrantCallCount()).To(Equal(0))
				Expect(targets.SetCurrentCallCount()).To(Equal(0))
				Expect(targets.PrintCurrentCallCount()).To(Equal(0))
			})
		})

		Describe("when a valid uaa is provided and password grant fails", func() {
			It("shows an error", func() {
				uaa.FetchInstanceReturns(&cf.Item{URL: "dummy-instance-url"}, map[string]interface{}{"uri": "dummy-uaa-url"})
				serviceInfo.ResolveJSONPathReturns("dummy-uaa-url")
				infoClient.ServerReturns(nil)
				tokenIssuer.PasswordGrantReturns(nil, fmt.Errorf("Invalid password"))

				testcmd.RunCLICommand(LoginCommand, []string{"--password", "user-password", "--secret", "client-secret", "dummy-uaa", "client-id", "some-user"}, ui)

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Invalid password"))
				Expect(uaa.FetchInstanceCallCount()).To(Equal(1))
				Expect(serviceInfo.ResolveJSONPathCallCount()).To(Equal(1))
				Expect(infoClientFactory.NewCallCount()).To(Equal(1))
				calledUrl, _, _ := infoClientFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(ui.PasswordPrompts).To(BeNil())
				Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
				calledUrl, calledClientID, calledClientSecret, calledSslVerify, calledCaCert := tokenIssuerFactory.NewArgsForCall(0)
				Expect(calledUrl).To(Equal("dummy-uaa-url"))
				Expect(calledClientID).To(Equal("client-id"))
				Expect(calledClientSecret).To(Equal("client-secret"))
				Expect(calledSslVerify).To(BeFalse())
				Expect(calledCaCert).To(Equal(""))
				Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(0))
				Expect(tokenIssuer.PasswordGrantCallCount()).To(Equal(1))
				calledUsername, calledPassword, calledScopes := tokenIssuer.PasswordGrantArgsForCall(0)
				Expect(calledUsername).To(Equal("some-user"))
				Expect(calledPassword).To(Equal("user-password"))
				Expect(calledScopes).To(BeNil())
				Expect(targets.SetCurrentCallCount()).To(Equal(0))
				Expect(targets.PrintCurrentCallCount()).To(Equal(0))
			})
		})
	})
})
