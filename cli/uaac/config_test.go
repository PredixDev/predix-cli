package uaac_test

import (
	"fmt"
	"io/ioutil"
	"os"

	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"
	"github.build.ge.com/adoption/predix-cli/cli/uaac"

	"github.com/PredixDev/go-uaa-lib"
	"github.com/PredixDev/go-uaa-lib/libfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		ui                    *testterm.FakeUI
		oldCurl               cf.CurlInterface
		oldTokenClaimsFactory lib.TokenClaimsFactoryInterface
		curl                  *cffakes.FakeCurlInterface
		tokenClaimsFactory    *libfakes.FakeTokenClaimsFactoryInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		curl = &cffakes.FakeCurlInterface{}
		tokenClaimsFactory = &libfakes.FakeTokenClaimsFactoryInterface{}

		global.Env.NoCache = true
		global.Env.ConfigDir, _ = ioutil.TempDir("", "test-config")
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

		oldCurl = cf.Curl
		cf.Curl = curl

		oldTokenClaimsFactory = lib.TokenClaimsFactory
		lib.TokenClaimsFactory = tokenClaimsFactory

		uaac.Targets.LoadConfig()
	})

	AfterEach(func() {
		lib.TokenClaimsFactory = oldTokenClaimsFactory
		cf.Curl = oldCurl
		_ = os.RemoveAll(global.Env.ConfigDir)
	})

	Context("get current", func() {
		Describe("when no target set", func() {
			It("shows an error", func() {
				testcmd.Run(func() {
					uaac.Targets.GetCurrent()
				})
				Expect(ui.Outputs).To(ConsistOf("FAILED", "No UAA target set. Login to a UAA using the 'predix uaa login' command"))
			})
		})

		Describe("when no context set", func() {
			It("shows an error", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{}, nil)

				testcmd.Run(func() {
					uaac.Targets.SetCurrent("dummy-uaa-url", "dummy-instance-url", false, "", &lib.TokenResponse{
						Type:    "Token-Type",
						Access:  "Access-Token",
						Refresh: "Refresh-Token",
					})

					uaac.Targets.GetCurrent()
				})
				Expect(ui.Outputs).To(ConsistOf("FAILED", "No UAA context set. Login to a UAA using the 'predix uaa login' command"))
			})
		})

		Describe("when failed to fetch uaa instance", func() {
			It("shows an error", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "Some-ID",
					ClientID: "dummy-client",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)

				testcmd.Run(func() {
					uaac.Targets.SetCurrent("dummy-uaa-url", "dummy-instance-url", false, "", &lib.TokenResponse{
						Type:    "Token-Type",
						Access:  "Access-Token",
						Refresh: "Refresh-Token",
					})

					uaac.Targets.GetCurrent()
				})
				Expect(ui.Outputs).To(ConsistOf("FAILED", "Failed to fetch target UAA's info"))
			})
		})
	})

	Context("set current", func() {
		Describe("when the context is a client", func() {
			It("sets the current target", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "Some-ID",
					ClientID: "dummy-client",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})

				var target *uaac.Target
				var context *uaac.Context
				var instance *cf.Item
				testcmd.Run(func() {
					uaac.Targets.SetCurrent("dummy-uaa-url", "dummy-instance-url", false, "", &lib.TokenResponse{
						Type:    "Token-Type",
						Access:  "Access-Token",
						Refresh: "Refresh-Token",
					})

					target, context, instance = uaac.Targets.GetCurrent()
				})

				Expect(ui.Outputs).To(BeNil())
				Expect(target).ToNot(BeNil())
				Expect(target.TargetURL).To(Equal("dummy-uaa-url"))
				Expect(target.CfInstanceURL).To(Equal("dummy-instance-url"))
				Expect(target.SkipSsl).To(BeFalse())
				Expect(target.CaCertFilePath).To(Equal(""))
				Expect(target.Current).To(BeTrue())
				Expect(context).ToNot(BeNil())
				Expect(context.ClientID).To(Equal("dummy-client"))
				Expect(context.Scopes).To(ConsistOf("scope1", "scope2"))
				Expect(context.Type).To(Equal("Token-Type"))
				Expect(context.Access).To(Equal("Access-Token"))
				Expect(context.Refresh).To(Equal("Refresh-Token"))
				Expect(instance).ToNot(BeNil())
			})
		})

		Describe("when the context is a user", func() {
			It("sets the current target", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "Some-ID",
					UserID:   "User-ID",
					UserName: "dummy-user",
					ClientID: "dummy-client",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})

				var target *uaac.Target
				var context *uaac.Context
				var instance *cf.Item
				testcmd.Run(func() {
					uaac.Targets.SetCurrent("dummy-uaa-url", "dummy-instance-url", false, "", &lib.TokenResponse{
						Type:    "Token-Type",
						Access:  "Access-Token",
						Refresh: "Refresh-Token",
					})
					target, context, instance = uaac.Targets.GetCurrent()
				})

				Expect(ui.Outputs).To(BeNil())
				Expect(target).ToNot(BeNil())
				Expect(target.TargetURL).To(Equal("dummy-uaa-url"))
				Expect(target.CfInstanceURL).To(Equal("dummy-instance-url"))
				Expect(target.SkipSsl).To(BeFalse())
				Expect(target.CaCertFilePath).To(Equal(""))
				Expect(target.Current).To(BeTrue())
				Expect(context).ToNot(BeNil())
				Expect(context.UserID).To(Equal("User-ID"))
				Expect(context.UserName).To(Equal("dummy-user"))
				Expect(context.ClientID).To(Equal("dummy-client"))
				Expect(context.Scopes).To(ConsistOf("scope1", "scope2"))
				Expect(context.Type).To(Equal("Token-Type"))
				Expect(context.Access).To(Equal("Access-Token"))
				Expect(context.Refresh).To(Equal("Refresh-Token"))
				Expect(instance).ToNot(BeNil())
			})
		})
	})

	Context("lookup and set current", func() {
		Describe("when the target and context are found", func() {
			It("returns true", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-2",
					UserID:   "User-ID",
					UserName: "dummy-user",
					ClientID: "client-2",
					Scopes:   []string{"scope3", "scope4"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-2", "instance-url-2", false, "", &lib.TokenResponse{
					Type:    "Token-Type-2",
					Access:  "Access-Token-2",
					Refresh: "Refresh-Token-2",
				})

				var result bool
				testcmd.Run(func() {
					result = uaac.Targets.LookupAndSetCurrent("uaa-url-1", "client-1")
				})

				Expect(ui.Outputs).To(BeNil())
				Expect(result).To(BeTrue())
			})
		})

		Describe("when the target and context are not found", func() {
			It("returns false", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-2",
					UserID:   "User-ID",
					UserName: "dummy-user",
					ClientID: "client-2",
					Scopes:   []string{"scope3", "scope4"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-2", "instance-url-2", false, "", &lib.TokenResponse{
					Type:    "Token-Type-2",
					Access:  "Access-Token-2",
					Refresh: "Refresh-Token-2",
				})

				var result bool
				testcmd.Run(func() {
					result = uaac.Targets.LookupAndSetCurrent("some-url", "client-1")
				})

				Expect(ui.Outputs).To(BeNil())
				Expect(result).To(BeFalse())
			})
		})
	})

	Context("set current for id", func() {
		It("sets the current target", func() {
			tokenClaimsFactory.NewReturns(&lib.TokenClaims{
				ID:       "ID-1",
				ClientID: "client-1",
				Scopes:   []string{"scope1", "scope2"},
			}, nil)
			curl.GetItemReturns(&cf.Item{})
			uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
				Type:    "Token-Type-1",
				Access:  "Access-Token-1",
				Refresh: "Refresh-Token-1",
			})

			tokenClaimsFactory.NewReturns(&lib.TokenClaims{
				ID:       "ID-2",
				UserID:   "User-ID",
				UserName: "dummy-user",
				ClientID: "client-2",
				Scopes:   []string{"scope3", "scope4"},
			}, nil)
			curl.GetItemReturns(&cf.Item{})
			uaac.Targets.SetCurrent("uaa-url-2", "instance-url-2", false, "", &lib.TokenResponse{
				Type:    "Token-Type-2",
				Access:  "Access-Token-2",
				Refresh: "Refresh-Token-2",
			})

			var target *uaac.Target
			var context *uaac.Context
			testcmd.Run(func() {
				target, context, _ = uaac.Targets.GetCurrent()
			})
			Expect(target).ToNot(BeNil())
			Expect(target.TargetURL).To(Equal("uaa-url-2"))
			Expect(target.CfInstanceURL).To(Equal("instance-url-2"))
			Expect(target.SkipSsl).To(BeFalse())
			Expect(target.CaCertFilePath).To(Equal(""))
			Expect(target.Current).To(BeTrue())
			Expect(context).ToNot(BeNil())
			Expect(context.UserID).To(Equal("User-ID"))
			Expect(context.UserName).To(Equal("dummy-user"))
			Expect(context.ClientID).To(Equal("client-2"))
			Expect(context.Scopes).To(ConsistOf("scope3", "scope4"))
			Expect(context.Type).To(Equal("Token-Type-2"))
			Expect(context.Access).To(Equal("Access-Token-2"))
			Expect(context.Refresh).To(Equal("Refresh-Token-2"))

			tokenClaimsFactory.NewReturns(&lib.TokenClaims{}, nil)

			testcmd.Run(func() {
				uaac.Targets.SetCurrent("uaa-url-3", "instance-url-3", false, "", &lib.TokenResponse{
					Type:    "Token-Type-3",
					Access:  "Access-Token-3",
					Refresh: "Refresh-Token-3",
				})
			})

			testcmd.Run(func() {
				uaac.Targets.SetCurrentForID(1)
			})

			testcmd.Run(func() {
				target, context, _ = uaac.Targets.GetCurrent()
			})
			Expect(ui.Outputs).To(BeNil())
			Expect(target).ToNot(BeNil())
			Expect(target.TargetURL).To(Equal("uaa-url-1"))
			Expect(target.CfInstanceURL).To(Equal("instance-url-1"))
			Expect(target.SkipSsl).To(BeFalse())
			Expect(target.CaCertFilePath).To(Equal(""))
			Expect(target.Current).To(BeTrue())
			Expect(context).ToNot(BeNil())
			Expect(context.UserID).To(Equal(""))
			Expect(context.UserName).To(Equal(""))
			Expect(context.ClientID).To(Equal("client-1"))
			Expect(context.Scopes).To(ConsistOf("scope1", "scope2"))
			Expect(context.Type).To(Equal("Token-Type-1"))
			Expect(context.Access).To(Equal("Access-Token-1"))
			Expect(context.Refresh).To(Equal("Refresh-Token-1"))

			testcmd.Run(func() {
				uaac.Targets.SetCurrentForID(3)
			})

			testcmd.Run(func() {
				target, context, _ = uaac.Targets.GetCurrent()
			})
			Expect(ui.Outputs).To(ConsistOf("FAILED", "No UAA context set. Login to a UAA using the 'predix uaa login' command"))
		})
	})

	Context("print all", func() {
		Describe("when loads the saved config data", func() {
			It("prints the targets", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-2",
					UserID:   "User-ID",
					UserName: "dummy-user",
					ClientID: "client-2",
					Scopes:   []string{"scope3", "scope4"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-2", "instance-url-2", false, "", &lib.TokenResponse{
					Type:    "Token-Type-2",
					Access:  "Access-Token-2",
					Refresh: "Refresh-Token-2",
				})

				tokenClaimsFactory.NewReturns(&lib.TokenClaims{}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-3", "instance-url-3", false, "", &lib.TokenResponse{
					Type:    "Token-Type-3",
					Access:  "Access-Token-3",
					Refresh: "Refresh-Token-3",
				})

				testcmd.Run(func() {
					uaac.Targets.LoadConfig()
					uaac.Targets.PrintAll()
				})

				Expect(ui.Outputs).To(ConsistOf(MatchRegexp(`ID\s+Target\s+Context`), MatchRegexp(`1\s+Client: client-1`),
					MatchRegexp(`2\s+User: dummy-user, Client: client-2`), MatchRegexp(`3\s+\*\s+No context`)))
			})
		})

		Describe("when no target set", func() {
			It("prints error", func() {
				testcmd.Run(func() {
					uaac.Targets.PrintAll()
				})
				Expect(ui.Outputs).To(ConsistOf("No UAA targets"))
			})
		})

		Describe("when failed to fetch UAA's info", func() {
			It("prints error", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				curl.GetItemReturns(nil)
				testcmd.Run(func() {
					uaac.Targets.PrintAll()
				})

				Expect(ui.Outputs).To(ConsistOf("Failed to fetch one or more target UAA's info"))
			})
		})
	})

	Context("print current", func() {
		Describe("when target is a client", func() {
			It("prints the target", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{
					Name: "some-uaa",
				})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				testcmd.Run(func() {
					uaac.Targets.PrintCurrent()
				})

				Expect(ui.Outputs).To(ConsistOf(MatchRegexp(`\s*`), MatchRegexp(`Target:\s+some-uaa`), MatchRegexp(`URL:\s+uaa-url-1`),
					MatchRegexp(`Client:\s+client-1`), MatchRegexp(`Access Token:\s+Access-Token-1`)))
			})
		})

		Describe("when target is a user", func() {
			It("prints the target", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					UserID:   "User-ID",
					UserName: "dummy-user",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{
					Name: "some-uaa",
				})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				testcmd.Run(func() {
					uaac.Targets.PrintCurrent()
				})

				Expect(ui.Outputs).To(ConsistOf(MatchRegexp(`\s*`), MatchRegexp(`Target:\s+some-uaa`), MatchRegexp(`URL:\s+uaa-url-1`),
					MatchRegexp(`User:\s+dummy-user`), MatchRegexp(`Client:\s+client-1`), MatchRegexp(`Access Token:\s+Access-Token-1`)))
			})
		})

		Describe("when no UAA target set", func() {
			It("prints error", func() {
				testcmd.Run(func() {
					uaac.Targets.PrintCurrent()
				})

				Expect(ui.Outputs).To(ConsistOf("No UAA target set"))
			})
		})

		Describe("when failed to fetch UAA's info", func() {
			It("prints error", func() {
				tokenClaimsFactory.NewReturns(&lib.TokenClaims{
					ID:       "ID-1",
					ClientID: "client-1",
					Scopes:   []string{"scope1", "scope2"},
				}, nil)
				curl.GetItemReturns(&cf.Item{})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				curl.GetItemReturns(nil)
				testcmd.Run(func() {
					uaac.Targets.PrintCurrent()
				})

				Expect(ui.Outputs).To(ConsistOf("Failed to fetch target UAA's info", MatchRegexp(`\s*`), MatchRegexp(`Target:\s*`),
					MatchRegexp(`URL:\s+uaa-url-1`), MatchRegexp(`Client:\s+client-1`), MatchRegexp(`Access Token:\s+Access-Token-1`)))
			})
		})

		Describe("when no context set", func() {
			It("prints error", func() {
				tokenClaimsFactory.NewReturns(nil, fmt.Errorf("Error"))
				curl.GetItemReturns(&cf.Item{
					Name: "some-uaa",
				})
				uaac.Targets.SetCurrent("uaa-url-1", "instance-url-1", false, "", &lib.TokenResponse{
					Type:    "Token-Type-1",
					Access:  "Access-Token-1",
					Refresh: "Refresh-Token-1",
				})

				testcmd.Run(func() {
					uaac.Targets.PrintCurrent()
				})

				Expect(ui.Outputs).To(ConsistOf("No context set", MatchRegexp(`\s*`),
					MatchRegexp(`Target:\s+some-uaa`), MatchRegexp(`URL:\s+uaa-url-1`)))
			})
		})
	})
})
