package createservice_test

import (
	"fmt"
	"io/ioutil"
	"os"

	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/createservice"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/commands/helpers/helpersfakes"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"
	"github.build.ge.com/adoption/predix-cli/cli/uaac/uaacfakes"
	"github.com/PredixDev/go-uaa-lib"
	"github.com/PredixDev/go-uaa-lib/libfakes"
	"github.com/urfave/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Service", func() {
	var (
		ui        *testterm.FakeUI
		sh        *testshell.FakeShell
		oldCurl   cf.CurlInterface
		oldLookup cf.LookupInterface
		oldUaa    helpers.UaaInterface
		curl      *cffakes.FakeCurlInterface
		lookup    *cffakes.FakeLookupInterface
		uaa       *helpersfakes.FakeUaaInterface
		guidCount int
		itemObj   func(name, resourceType string) *cf.Item
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		sh = &testshell.FakeShell{}
		curl = &cffakes.FakeCurlInterface{}
		lookup = &cffakes.FakeLookupInterface{}
		uaa = &helpersfakes.FakeUaaInterface{}

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

		oldUaa = helpers.Uaa
		helpers.Uaa = uaa

		guidCount = 1
		itemObj = func(name, resourceType string) *cf.Item {
			defer func() { guidCount = guidCount + 1 }()
			return &cf.Item{
				Name: name,
				GUID: fmt.Sprintf("00000000-0000-0000-0000-000000%02d", guidCount),
				URL:  fmt.Sprintf("/v2/%s/00000000-0000-0000-0000-000000%02d", resourceType, guidCount),
			}
		}
	})

	AfterEach(func() {
		helpers.Uaa = oldUaa
		cf.Lookup = oldLookup
		cf.Curl = oldCurl
	})

	Context("before", func() {
		Describe("when called with less than 3 args", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(CreateServiceCommand, nil, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())

				err = testcmd.BeforeCLICommand(CreateServiceCommand, []string{"arg1", "arg2"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when creating something other than uaa with less than 4 args", func() {
			It("returns error", func() {
				err := testcmd.BeforeCLICommand(CreateServiceCommand, []string{"predix-asset", "arg2", "arg3"}, ui)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Incorrect Usage"))
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when creating for an unknown service", func() {
			It("shows an error", func() {
				err := testcmd.BeforeCLICommand(CreateServiceCommand, []string{"dummy-service", "arg2", "arg3", "arg4"}, ui)
				Expect(err).To(BeNil())
				Expect(ui.Outputs).To(ConsistOf("FAILED", "No handler found for service 'dummy-service'"))
			})
		})

		Describe("when creating predix-uaa with 3 args", func() {
			It("returns nil", func() {
				err := testcmd.BeforeCLICommand(CreateServiceCommand, []string{"predix-uaa", "Tiered", "test-uaa"}, ui)
				Expect(err).To(BeNil())
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when creating a registered service with 4 args", func() {
			It("returns nil", func() {
				err := testcmd.BeforeCLICommand(CreateServiceCommand, []string{"predix-asset", "Tiered", "test-asset", "test-uaa"}, ui)
				Expect(err).To(BeNil())
				Expect(ui.Outputs).To(BeNil())
			})
		})

		Describe("when creating a registered service with 5 args", func() {
			It("returns nil", func() {
				err := testcmd.BeforeCLICommand(CreateServiceCommand, []string{"predix-timeseries", "Bronze", "test-timeseries", "test-uaa", "test-asset"}, ui)
				Expect(err).To(BeNil())
				Expect(ui.Outputs).To(BeNil())
			})
		})
	})

	Context("bash complete", func() {
		Describe("when called with 0 args", func() {
			It("shows marketplace services", func() {
				lookup.MarketplaceServicesReturns([]string{"Service1", "Service2"})

				testcmd.BashCompleteCLICommand(CreateServiceCommand, nil, ui)

				Expect(ui.Outputs).To(ConsistOf("Service1", "Service2"))
			})
		})

		Describe("when called with 1 arg", func() {
			It("shows marketplace services", func() {
				lookup.MarketplaceServicePlansReturns([]string{"Plan1", "Plan2"})

				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"SERVICE"}, ui)

				Expect(ui.Outputs).To(ConsistOf("Plan1", "Plan2"))
				Expect(lookup.MarketplaceServicePlansArgsForCall(0)).To(Equal("SERVICE"))
			})
		})

		Describe("when called with 2 args", func() {
			It("shows prompt", func() {
				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"SERVICE", "PLAN"}, ui)

				Expect(ui.Outputs).To(ConsistOf("Enter: NAME_FOR_NEW_SERVICE_INSTANCE", "_"))
			})
		})

		Describe("when called for predix-uaa with 3 args", func() {
			It("shows parameters as completion", func() {
				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"predix-uaa", "Tiered", "test-uaa"}, ui)

				Expect(ui.Outputs).To(ConsistOf("--admin-secret", "-a", "--client-id", "-c", "--client-secret", "-s", "--skip-ssl-validation", "--ca-cert"))
			})
		})

		Describe("when called for other services with 3 args", func() {
			It("shows predix-uaa instances as completion", func() {
				lookup.PredixUaaInstancesReturns([]string{"test-uaa-1", "test-uaa-2"})

				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"dummy-service", "Tiered", "test-service"}, ui)

				Expect(ui.Outputs).To(ConsistOf("test-uaa-1", "test-uaa-2"))
			})
		})

		Describe("when called for predix-analytics-runtime with 4 args", func() {
			It("shows predix-asset instances as completion", func() {
				lookup.PredixAssetInstancesReturns([]string{"test-asset-1", "test-asset-2"})

				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"predix-analytics-runtime", "Tiered", "test-runtime", "test-uaa"}, ui)

				Expect(ui.Outputs).To(ConsistOf("test-asset-1", "test-asset-2"))
			})
		})

		Describe("when called for predix-analytics-runtime with 5 args", func() {
			It("shows predix-timeseries instances as completion", func() {
				lookup.PredixTimeseriesInstancesReturns([]string{"test-timeseries-1", "test-timeseries-2"})

				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"predix-analytics-runtime", "Tiered", "test-runtime", "test-uaa", "test-asset"}, ui)

				Expect(ui.Outputs).To(ConsistOf("test-timeseries-1", "test-timeseries-2"))
			})
		})

		Describe("when called for predix-analytics-runtime with 6 args", func() {
			It("shows predix-analytics-catalog instances as completion", func() {
				lookup.PredixAnalyticsCatalogInstancesReturns([]string{"test-analytics-catalog-1", "test-analytics-catalog-2"})

				testcmd.BashCompleteCLICommand(CreateServiceCommand, []string{"predix-analytics-runtime", "Tiered", "test-runtime", "test-uaa", "test-asset", "test-timeseries"}, ui)

				Expect(ui.Outputs).To(ConsistOf("test-analytics-catalog-1", "test-analytics-catalog-2"))
			})
		})
	})

	Context("action", func() {
		Describe("when called for unknown service", func() {
			It("shows an error", func() {
				testcmd.RunCLICommand(CreateServiceCommand, []string{"dummy-service"}, ui)
				Expect(ui.Outputs).To(ConsistOf("", "FAILED", "No handler found for service 'dummy-service'"))
			})
		})

		Describe("when called for known service", func() {
			It("calls the handler", func() {
				handlerCalled := false
				CreateServiceHandlers["known-dummy-service"] = func(c *cli.Context) {
					Expect(c).ToNot(BeNil())
					handlerCalled = true
				}
				testcmd.RunCLICommand(CreateServiceCommand, []string{"known-dummy-service"}, ui)
				Expect(ui.Outputs).To(ConsistOf(""))
				Expect(handlerCalled).To(BeTrue())
			})
		})
	})

	Context("ask for client id", func() {
		Describe("when client id not provided", func() {
			It("asks", func() {
				ui.Inputs = []string{"test-client"}
				Expect(AskForClientID(testcmd.Context(CreateServiceCommand, nil))).To(Equal("test-client"))
				Expect(ui.Prompts).To(ConsistOf("Client ID"))
			})
		})

		Describe("when client id is provided", func() {
			It("uses it", func() {
				Expect(AskForClientID(testcmd.Context(CreateServiceCommand, []string{"--client-id", "test-client"}))).To(Equal("test-client"))
			})
		})
	})

	Context("ask for client secret", func() {
		Describe("when client secret not provided", func() {
			It("asks", func() {
				ui.Inputs = []string{"test-secret", "test-secret"}
				Expect(AskForClientSecret(testcmd.Context(CreateServiceCommand, nil))).To(Equal("test-secret"))
				Expect(ui.PasswordPrompts).To(ConsistOf("Client Secret", "Verify Client Secret"))
			})
		})

		Describe("when client secret is provided", func() {
			It("uses it", func() {
				Expect(AskForClientSecret(testcmd.Context(CreateServiceCommand, []string{"--client-secret", "test-secret"}))).To(Equal("test-secret"))
			})
		})
	})

	Context("verify service plan", func() {
		Describe("when service plan does not exist", func() {
			It("shows an error", func() {
				testcmd.Run(func() {
					VerifyServicePlan(testcmd.Context(CreateServiceCommand, []string{"dummy-service", "dummy-plan"}))
				})

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Service plan not found"))
			})
		})

		Describe("when service plan exists", func() {
			It("returns the plan", func() {
				lookup.MarketplaceServicePlanItemReturns(itemObj("dummy-plan", "service_plans"))

				_, plan := VerifyServicePlan(testcmd.Context(CreateServiceCommand, []string{"dummy-service", "dummy-plan"}))

				Expect(ui.Outputs).To(BeNil())
				Expect(plan).ToNot(BeNil())
				Expect(plan.Name).To(Equal("dummy-plan"))

				calledService, calledPlan := lookup.MarketplaceServicePlanItemArgsForCall(0)
				Expect(calledService).To(Equal("dummy-service"))
				Expect(calledPlan).To(Equal("dummy-plan"))
			})
		})
	})

	Context("verify predix uaa instance", func() {
		It("calls helpers uaa fetch instance", func() {
			uaa.FetchInstanceReturns(nil, nil)

			testcmd.Run(func() {
				VerifyPredixUaaInstance([]string{"arg1", "arg2", "arg3", "dummy-uaa"})
			})

			Expect(ui.Outputs).To(BeNil())
			Expect(uaa.FetchInstanceCallCount()).To(Equal(1))
		})
	})

	Context("fake uaa target", func() {
		var (
			oldTokenIssuerFactory lib.TokenIssuerFactoryInterface
			oldTokenClaimsFactory lib.TokenClaimsFactoryInterface
			oldTargets            uaac.TargetsInterface
			tokenIssuer           *libfakes.FakeTokenIssuer
			tokenIssuerFactory    *libfakes.FakeTokenIssuerFactoryInterface
			tokenClaimsFactory    *libfakes.FakeTokenClaimsFactoryInterface
			targets               *uaacfakes.FakeTargetsInterface
		)
		BeforeEach(func() {
			oldTokenIssuerFactory = lib.TokenIssuerFactory
			oldTokenClaimsFactory = lib.TokenClaimsFactory
			oldTargets = uaac.Targets

			tokenIssuer = &libfakes.FakeTokenIssuer{}
			tokenIssuerFactory = &libfakes.FakeTokenIssuerFactoryInterface{}
			tokenClaimsFactory = &libfakes.FakeTokenClaimsFactoryInterface{}
			targets = &uaacfakes.FakeTargetsInterface{}

			global.Env.ConfigDir, _ = ioutil.TempDir("", "test-config")
			tokenIssuerFactory.NewReturns(tokenIssuer)
			lib.TokenIssuerFactory = tokenIssuerFactory
			lib.TokenClaimsFactory = tokenClaimsFactory
			uaac.Targets = targets
		})
		AfterEach(func() {
			uaac.Targets = oldTargets
			lib.TokenClaimsFactory = oldTokenClaimsFactory
			lib.TokenIssuerFactory = oldTokenIssuerFactory
			_ = os.RemoveAll(global.Env.ConfigDir)
		})

		Context("set current target", func() {
			Describe("when no current target set", func() {
				It("issues token and sets the current target", func() {
					targets.LookupAndSetCurrentReturns(false)
					tokenIssuer.ClientCredentialsGrantReturns(&lib.TokenResponse{
						Type:   "Token-Type",
						Access: "Access-Token",
					}, nil)

					uaa.AskForAdminClientSecretReturns("secret")
					testcmd.Run(func() {
						SetCurrentTarget(testcmd.Context(CreateServiceCommand, []string{"--admin-secret", "secret"}),
							&cf.Item{
								URL: "dummy-instance-url",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(targets.LookupAndSetCurrentCallCount()).To(Equal(1))
					calledUrl, calledClientID := targets.LookupAndSetCurrentArgsForCall(0)
					Expect(calledUrl).To(Equal("dummy-uaa-url"))
					Expect(calledClientID).To(Equal("admin"))
					Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
					Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(1))
					calledTarget, calledClientID, calledClientSecret, calledVerifySSL, calledCaCertFile := tokenIssuerFactory.NewArgsForCall(0)
					Expect(calledTarget).To(Equal("dummy-uaa-url"))
					Expect(calledClientID).To(Equal("admin"))
					Expect(calledClientSecret).To(Equal("secret"))
					Expect(calledVerifySSL).To(BeFalse())
					Expect(calledCaCertFile).To(Equal(""))
					Expect(targets.SetCurrentCallCount()).To(Equal(1))
					calledUrl, calledInstanceUrl, calledVerifySSL, calledCaCertFile, calledTokenResponse := targets.SetCurrentArgsForCall(0)
					Expect(calledUrl).To(Equal("dummy-uaa-url"))
					Expect(calledInstanceUrl).To(Equal("dummy-instance-url"))
					Expect(calledVerifySSL).To(BeFalse())
					Expect(calledCaCertFile).To(Equal(""))
					Expect(calledTokenResponse).ToNot(BeNil())
					Expect(targets.PrintCurrentCallCount()).To(Equal(1))
				})
			})

			Describe("when the current target is set to something else and new target is unknown", func() {
				It("issues token and sets the current target", func() {
					targets.LookupAndSetCurrentReturns(false)
					tokenIssuer.ClientCredentialsGrantReturns(&lib.TokenResponse{
						Type:   "Token-Type",
						Access: "Access-Token",
					}, nil)

					uaa.AskForAdminClientSecretReturns("secret")
					testcmd.Run(func() {
						SetCurrentTarget(testcmd.Context(CreateServiceCommand, []string{"--admin-secret", "secret"}),
							&cf.Item{
								URL: "dummy-instance-url",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(targets.LookupAndSetCurrentCallCount()).To(Equal(1))
					calledUrl, calledClientID := targets.LookupAndSetCurrentArgsForCall(0)
					Expect(calledUrl).To(Equal("dummy-uaa-url"))
					Expect(calledClientID).To(Equal("admin"))
					Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
					Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(1))
					calledTarget, calledClientID, calledClientSecret, calledVerifySSL, calledCaCertFile := tokenIssuerFactory.NewArgsForCall(0)
					Expect(calledTarget).To(Equal("dummy-uaa-url"))
					Expect(calledClientID).To(Equal("admin"))
					Expect(calledClientSecret).To(Equal("secret"))
					Expect(calledVerifySSL).To(BeFalse())
					Expect(calledCaCertFile).To(Equal(""))
					Expect(targets.SetCurrentCallCount()).To(Equal(1))
					calledUrl, calledInstanceUrl, calledVerifySSL, calledCaCertFile, calledTokenResponse := targets.SetCurrentArgsForCall(0)
					Expect(calledUrl).To(Equal("dummy-uaa-url"))
					Expect(calledInstanceUrl).To(Equal("dummy-instance-url"))
					Expect(calledVerifySSL).To(BeFalse())
					Expect(calledCaCertFile).To(Equal(""))
					Expect(calledTokenResponse).ToNot(BeNil())
					Expect(targets.PrintCurrentCallCount()).To(Equal(1))
				})
			})

			Describe("when the current target is set to something else and new target is known", func() {
				It("does not issue token and sets the current target", func() {
					targets.LookupAndSetCurrentReturns(true)
					testcmd.Run(func() {
						SetCurrentTarget(testcmd.Context(CreateServiceCommand, []string{"--admin-secret", "secret"}),
							&cf.Item{
								URL: "dummy-instance-url",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(targets.LookupAndSetCurrentCallCount()).To(Equal(1))
					calledUrl, calledClientID := targets.LookupAndSetCurrentArgsForCall(0)
					Expect(calledUrl).To(Equal("dummy-uaa-url"))
					Expect(calledClientID).To(Equal("admin"))
					Expect(tokenIssuerFactory.NewCallCount()).To(Equal(0))
					Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(0))
					Expect(targets.SetCurrentCallCount()).To(Equal(0))
					Expect(targets.PrintCurrentCallCount()).To(Equal(1))
				})
			})

			Describe("when client credentials grant returns and error", func() {
				It("shows an error", func() {
					tokenIssuer.ClientCredentialsGrantReturns(nil, fmt.Errorf("Error"))

					testcmd.Run(func() {
						SetCurrentTarget(testcmd.Context(CreateServiceCommand, []string{"--admin-secret", "secret"}),
							&cf.Item{
								URL: "dummy-instance-url",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(tokenIssuerFactory.NewCallCount()).To(Equal(1))
					Expect(tokenIssuer.ClientCredentialsGrantCallCount()).To(Equal(1))
					Expect(ui.Outputs).To(ConsistOf("FAILED", "Error"))
				})
			})
		})

		Context("verify client id and secret", func() {
			var (
				oldScimFactory lib.ScimFactoryInterface
				scim           libfakes.FakeScim
				scimFactory    libfakes.FakeScimFactoryInterface
			)
			BeforeEach(func() {
				oldScimFactory = lib.ScimFactory
				scim = libfakes.FakeScim{}
				scimFactory = libfakes.FakeScimFactoryInterface{}

				scimFactory.NewReturns(&scim)

				lib.ScimFactory = &scimFactory
			})
			AfterEach(func() {
				lib.ScimFactory = oldScimFactory
			})

			Describe("when the target is set and client does not exist", func() {
				It("says it will create client and add authorities", func() {
					tokenClaimsFactory.NewReturns(&lib.TokenClaims{
						ClientID: "admin",
					}, nil)
					targets.GetCurrentReturns(&uaac.Target{}, &uaac.Context{}, &cf.Item{
						Name: "dummy-uaa",
						URL:  "/v2/service_instances/00000000-0000-0000-0000-00000001",
					})
					scim.GetClientReturns(nil, nil)

					var clientID, clientSecret string
					var client *lib.Client
					var returnedScim lib.Scim
					testcmd.Run(func() {
						clientID, clientSecret, client, returnedScim = VerifyClientIDAndSecret(testcmd.Context(CreateServiceCommand, []string{"--client-id", "dummy-client", "--client-secret", "dummy-secret"}),
							&cf.Item{
								Name: "dummy-uaa",
								URL:  "/v2/service_instances/00000000-0000-0000-0000-00000001",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(clientID).To(Equal("dummy-client"))
					Expect(clientSecret).To(Equal("dummy-secret"))
					Expect(client).To(BeNil())
					Expect(returnedScim).ToNot(BeNil())

					Expect(ui.Outputs).To(ConsistOf("Checking if client dummy-client exists on service instance dummy-uaa",
						"Client dummy-client does not exist. It will be created with the required authorities."))
				})
			})

			Describe("when the target is set and client already exists", func() {
				It("says it will add authorities", func() {
					tokenClaimsFactory.NewReturns(&lib.TokenClaims{
						ClientID: "admin",
					}, nil)
					targets.GetCurrentReturns(&uaac.Target{}, &uaac.Context{}, &cf.Item{
						Name: "dummy-uaa",
						URL:  "/v2/service_instances/00000000-0000-0000-0000-00000001",
					})
					scim.GetClientReturns(&lib.Client{}, nil)

					var clientID, clientSecret string
					var client *lib.Client
					var returnedScim lib.Scim
					testcmd.Run(func() {
						clientID, clientSecret, client, returnedScim = VerifyClientIDAndSecret(testcmd.Context(CreateServiceCommand, []string{"--client-id", "dummy-client"}),
							&cf.Item{
								Name: "dummy-uaa",
								URL:  "/v2/service_instances/00000000-0000-0000-0000-00000001",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(clientID).To(Equal("dummy-client"))
					Expect(clientSecret).To(Equal(""))
					Expect(client).ToNot(BeNil())
					Expect(returnedScim).ToNot(BeNil())

					Expect(ui.Outputs).To(ConsistOf("Checking if client dummy-client exists on service instance dummy-uaa",
						"Client dummy-client exists. The required authorities will be added to it."))
				})
			})

			Describe("when some other UAA is targeted", func() {
				It("shows an error", func() {
					tokenClaimsFactory.NewReturns(&lib.TokenClaims{
						ClientID: "admin",
					}, nil)
					targets.GetCurrentReturns(&uaac.Target{}, &uaac.Context{}, &cf.Item{
						Name: "dummy-uaa",
						URL:  "dummy-instance-url",
					})

					testcmd.Run(func() {
						VerifyClientIDAndSecret(testcmd.Context(CreateServiceCommand, []string{"--client-id", "dummy-client"}),
							&cf.Item{
								Name: "some-other-uaa",
								URL:  "some-instance-url",
							}, map[string]interface{}{"uri": "some-uaa-url"})
					})
					Expect(ui.Outputs).To(ConsistOf("FAILED", "Incorrect target UAA, should be some-other-uaa"))
				})
			})

			Describe("when failed to fetch client", func() {
				It("shows an error", func() {
					tokenClaimsFactory.NewReturns(&lib.TokenClaims{
						ClientID: "admin",
					}, nil)
					targets.GetCurrentReturns(&uaac.Target{}, &uaac.Context{}, &cf.Item{
						Name: "dummy-uaa",
						URL:  "/v2/service_instances/00000000-0000-0000-0000-00000001",
					})
					scim.GetClientReturns(nil, fmt.Errorf("Failed to fetch client"))

					testcmd.Run(func() {
						VerifyClientIDAndSecret(testcmd.Context(CreateServiceCommand, []string{"--client-id", "dummy-client"}),
							&cf.Item{
								Name: "dummy-uaa",
								URL:  "/v2/service_instances/00000000-0000-0000-0000-00000001",
							}, map[string]interface{}{"uri": "dummy-uaa-url"})
					})

					Expect(ui.Outputs).To(ConsistOf("Checking if client dummy-client exists on service instance dummy-uaa",
						"FAILED", "Failed to fetch client"))
				})
			})
		})
	})

	Context("create instance with trusted issuer ids", func() {
		Describe("when the instance is created successfully", func() {
			It("says ok", func() {
				curl.PostItemReturns(&cf.Item{}, nil)

				var instance *cf.Item
				testcmd.Run(func() {
					instance = InstanceWithTrustedIssuerIDs([]string{"service-name", "service-plan", "service-instance-name"},
						&cf.Item{
							GUID: "plan-guid",
						}, map[string]interface{}{"issuerId": "uaa-issuer-url"})
				})

				Expect(instance).ToNot(BeNil())
				Expect(ui.Outputs).To(ConsistOf("Creating service instance service-instance-name in org FakeOrg / space FakeSpace as FakeName", "OK"))
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/service_instances?accepts_incomplete=true"))
				Expect(calledData).To(Equal(`{"name":"service-instance-name","space_guid":"FakeSpaceGUID","service_plan_guid":"plan-guid","parameters":{"trustedIssuerIds":["uaa-issuer-url"]}}`))
			})
		})

		Describe("when the instance creation fails", func() {
			It("shows an error", func() {
				curl.PostItemReturns(nil, fmt.Errorf("Failed to create instance"))

				testcmd.Run(func() {
					InstanceWithTrustedIssuerIDs([]string{"service-name", "service-plan", "service-instance-name"},
						&cf.Item{
							GUID: "plan-guid",
						}, map[string]interface{}{"issuerId": "uaa-issuer-url"})
				})

				Expect(ui.Outputs).To(ConsistOf("Creating service instance service-instance-name in org FakeOrg / space FakeSpace as FakeName", "FAILED", "Failed to create instance"))
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/service_instances?accepts_incomplete=true"))
				Expect(calledData).To(Equal(`{"name":"service-instance-name","space_guid":"FakeSpaceGUID","service_plan_guid":"plan-guid","parameters":{"trustedIssuerIds":["uaa-issuer-url"]}}`))
			})
		})
	})

	Context("create instance with parameters", func() {
		Describe("when the instance is created successfully", func() {
			It("says ok", func() {
				curl.PostItemReturns(&cf.Item{}, nil)

				var instance *cf.Item
				testcmd.Run(func() {
					instance = InstanceWithParameters([]string{"service-name", "service-plan", "service-instance-name"},
						&cf.Item{
							GUID: "plan-guid",
						}, map[string]interface{}{"issuerId": "uaa-issuer-url"})
				})

				Expect(instance).ToNot(BeNil())
				Expect(ui.Outputs).To(ConsistOf("Creating service instance service-instance-name in org FakeOrg / space FakeSpace as FakeName", "OK"))
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/service_instances?accepts_incomplete=true"))
				Expect(calledData).To(Equal(`{"name":"service-instance-name","space_guid":"FakeSpaceGUID","service_plan_guid":"plan-guid","parameters":{"issuerId":"uaa-issuer-url"}}`))
			})
		})

		Describe("when the instance creation fails", func() {
			It("shows an error", func() {
				curl.PostItemReturns(nil, fmt.Errorf("Failed to create instance"))

				testcmd.Run(func() {
					parameters := FakeParameters{}
					parameters.TrustedIssuers = []string{"uaa-issuer-url"}
					parameters.TrustedClient.ID = "client-id"
					parameters.TrustedClient.Secret = "client-secret"
					InstanceWithParameters([]string{"service-name", "service-plan", "service-instance-name"},
						&cf.Item{
							GUID: "plan-guid",
						}, parameters)
				})

				Expect(ui.Outputs).To(ConsistOf("Creating service instance service-instance-name in org FakeOrg / space FakeSpace as FakeName", "FAILED", "Failed to create instance"))
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/service_instances?accepts_incomplete=true"))
				Expect(calledData).To(Equal(`{"name":"service-instance-name","space_guid":"FakeSpaceGUID","service_plan_guid":"plan-guid","parameters":{"trustedIssuerIds":["uaa-issuer-url"],"trustedClientCredential":{"clientId":"client-id","clientSecret":"client-secret"}}}`))
			})
		})
	})

	Context("create or update client", func() {
		Describe("when client does not exist", func() {
			It("creates client", func() {
				scim := &libfakes.FakeScim{}
				testcmd.Run(func() {
					CreateOrUpdateClient("some-client", "some-secret", nil, scim, &cf.Item{
						Name: "some-uaa",
					}, []string{"scope1"}, []string{"authority1"})
				})

				Expect(ui.Outputs).To(ConsistOf("", "Creating client some-client on Predix UAA instance some-uaa", "OK"))
				Expect(scim.CreateClientCallCount()).To(Equal(1))
				Expect(scim.PutClientCallCount()).To(Equal(0))
				Expect(scim.CreateClientArgsForCall(0)).ToNot(BeNil())
				Expect(scim.CreateClientArgsForCall(0).ID).To(Equal("some-client"))
				Expect(scim.CreateClientArgsForCall(0).Secret).To(Equal("some-secret"))
				Expect(scim.CreateClientArgsForCall(0).Scopes).To(ConsistOf("uaa.none", "openid", "scope1"))
				Expect(scim.CreateClientArgsForCall(0).GrantTypes).To(ConsistOf("authorization_code", "client_credentials", "refresh_token", "password"))
				Expect(scim.CreateClientArgsForCall(0).Authorities).To(ConsistOf("openid", "uaa.none", "uaa.resource", "authority1"))
				Expect(scim.CreateClientArgsForCall(0).AutoApprove).To(ConsistOf("openid"))
			})

			It("tries to create client and shows error on failure", func() {
				scim := &libfakes.FakeScim{}
				scim.CreateClientReturns(fmt.Errorf("Failed to create client"))
				testcmd.Run(func() {
					CreateOrUpdateClient("some-client", "some-secret", nil, scim, &cf.Item{
						Name: "some-uaa",
					}, []string{"scope1"}, []string{"authority1"})
				})

				Expect(ui.Outputs).To(ConsistOf("", "Creating client some-client on Predix UAA instance some-uaa", "FAILED", "Failed to create client"))
				Expect(scim.CreateClientCallCount()).To(Equal(1))
				Expect(scim.PutClientCallCount()).To(Equal(0))
				Expect(scim.CreateClientArgsForCall(0)).ToNot(BeNil())
				Expect(scim.CreateClientArgsForCall(0).ID).To(Equal("some-client"))
				Expect(scim.CreateClientArgsForCall(0).Secret).To(Equal("some-secret"))
				Expect(scim.CreateClientArgsForCall(0).Scopes).To(ConsistOf("uaa.none", "openid", "scope1"))
				Expect(scim.CreateClientArgsForCall(0).GrantTypes).To(ConsistOf("authorization_code", "client_credentials", "refresh_token", "password"))
				Expect(scim.CreateClientArgsForCall(0).Authorities).To(ConsistOf("openid", "uaa.none", "uaa.resource", "authority1"))
				Expect(scim.CreateClientArgsForCall(0).AutoApprove).To(ConsistOf("openid"))
			})
		})

		Describe("when client already exist", func() {
			It("updates client", func() {
				scim := &libfakes.FakeScim{}
				scim.PutClientReturns(fmt.Errorf("Failed to update client"))
				testcmd.Run(func() {
					CreateOrUpdateClient("some-client", "", &lib.Client{
						ID:          "some-client",
						Scopes:      []string{"previous-scope"},
						GrantTypes:  []string{"authorization_code", "client_credentials"},
						Authorities: []string{"previous-authorities"},
						AutoApprove: []string{"previous-autoapprove"},
					}, scim, &cf.Item{
						Name: "some-uaa",
					}, []string{"scope1"}, []string{"authority1"})
				})

				Expect(ui.Outputs).To(ConsistOf("", "Updating client some-client on Predix UAA instance some-uaa", "FAILED", "Failed to update client"))
				Expect(scim.CreateClientCallCount()).To(Equal(0))
				Expect(scim.PutClientCallCount()).To(Equal(1))
				Expect(scim.PutClientArgsForCall(0)).ToNot(BeNil())
				Expect(scim.PutClientArgsForCall(0).ID).To(Equal("some-client"))
				Expect(scim.PutClientArgsForCall(0).Secret).To(Equal(""))
				Expect(scim.PutClientArgsForCall(0).Scopes).To(ConsistOf("previous-scope", "scope1"))
				Expect(scim.PutClientArgsForCall(0).GrantTypes).To(ConsistOf("authorization_code", "client_credentials"))
				Expect(scim.PutClientArgsForCall(0).Authorities).To(ConsistOf("previous-authorities", "authority1"))
				Expect(scim.PutClientArgsForCall(0).AutoApprove).To(ConsistOf("previous-autoapprove"))
			})

			It("tries to update client and shows error on failure", func() {
				scim := &libfakes.FakeScim{}
				testcmd.Run(func() {
					CreateOrUpdateClient("some-client", "", &lib.Client{
						ID:          "some-client",
						Scopes:      []string{"previous-scope"},
						GrantTypes:  []string{"authorization_code", "client_credentials"},
						Authorities: []string{"previous-authorities"},
						AutoApprove: []string{"previous-autoapprove"},
					}, scim, &cf.Item{
						Name: "some-uaa",
					}, []string{"scope1"}, []string{"authority1"})
				})

				Expect(ui.Outputs).To(ConsistOf("", "Updating client some-client on Predix UAA instance some-uaa", "OK"))
				Expect(scim.CreateClientCallCount()).To(Equal(0))
				Expect(scim.PutClientCallCount()).To(Equal(1))
				Expect(scim.PutClientArgsForCall(0)).ToNot(BeNil())
				Expect(scim.PutClientArgsForCall(0).ID).To(Equal("some-client"))
				Expect(scim.PutClientArgsForCall(0).Secret).To(Equal(""))
				Expect(scim.PutClientArgsForCall(0).Scopes).To(ConsistOf("previous-scope", "scope1"))
				Expect(scim.PutClientArgsForCall(0).GrantTypes).To(ConsistOf("authorization_code", "client_credentials"))
				Expect(scim.PutClientArgsForCall(0).Authorities).To(ConsistOf("previous-authorities", "authority1"))
				Expect(scim.PutClientArgsForCall(0).AutoApprove).To(ConsistOf("previous-autoapprove"))
			})
		})
	})
})

type FakeParameters struct {
	TrustedIssuers []string `json:"trustedIssuerIds"`
	TrustedClient  struct {
		ID     string `json:"clientId,omitempty"`
		Secret string `json:"clientSecret,omitempty"`
	} `json:"trustedClientCredential,omitempty"`
}
