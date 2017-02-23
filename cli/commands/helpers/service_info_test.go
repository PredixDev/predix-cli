package helpers_test

import (
	"fmt"

	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	. "github.build.ge.com/adoption/predix-cli/cli/commands/helpers"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceInfo", func() {
	var (
		ui      *testterm.FakeUI
		oldCurl cf.CurlInterface
		curl    *cffakes.FakeCurlInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		curl = &cffakes.FakeCurlInterface{}

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

		oldCurl = cf.Curl
		cf.Curl = curl
	})

	AfterEach(func() {
		cf.Curl = oldCurl
	})

	Context("fetch", func() {
		Describe("when successfully fetched info", func() {
			It("returns the info", func() {
				curl.PostItemReturns(&cf.Item{
					GUID: "app-guid",
					URL:  "app-url",
				}, nil)
				resource := &cf.CurlResource{}
				resource.Metadata.URL = "binding-url"
				resource.Entity.Credentials = map[string]interface{}{
					"some-key": "some-value",
					"some-other-key": map[string]interface{}{
						"nested-key": "some-other-value",
					},
				}
				curl.PostResourceReturns(resource, nil)

				var creds map[string]interface{}
				testcmd.Run(func() {
					creds = ServiceInfo.FetchFor(&cf.Item{
						Name: "dummy-service",
						GUID: "dummy-service-guid",
					})
				})

				Expect(creds).To(Equal(resource.Entity.Credentials))
				Expect(ui.Outputs).To(BeNil())
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/apps"))
				Expect(calledData).To(MatchRegexp(`{"name":"predix-cli-.*","space_guid":"FakeSpaceGUID","memory": 1,"instances": 1}`))
				Expect(curl.PostResourceCallCount()).To(Equal(1))
				calledPath, calledData = curl.PostResourceArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/service_bindings"))
				Expect(calledData).To(Equal(`{"service_instance_guid":"dummy-service-guid","app_guid":"app-guid"}`))
				Expect(curl.DeleteCallCount()).To(Equal(2))
				Expect(curl.DeleteArgsForCall(0)).To(Equal("binding-url"))
				Expect(curl.DeleteArgsForCall(1)).To(Equal("app-url"))
			})
		})

		Describe("when failed to create app", func() {
			It("shows an error", func() {
				curl.PostItemReturns(nil, fmt.Errorf("Failed to create app"))

				testcmd.Run(func() {
					ServiceInfo.FetchFor(&cf.Item{
						Name: "dummy-service",
						GUID: "dummy-service-guid",
					})
				})

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Unable to get info for service instance dummy-service"))
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/apps"))
				Expect(calledData).To(MatchRegexp(`{"name":"predix-cli-.*","space_guid":"FakeSpaceGUID","memory": 1,"instances": 1}`))
				Expect(curl.PostResourceCallCount()).To(Equal(0))
				Expect(curl.DeleteCallCount()).To(Equal(0))
			})
		})

		Describe("when failed to create binding", func() {
			It("shows an error", func() {
				curl.PostItemReturns(&cf.Item{
					GUID: "app-guid",
					URL:  "app-url",
				}, nil)
				curl.PostResourceReturns(nil, fmt.Errorf("Failed to create binding"))

				testcmd.Run(func() {
					ServiceInfo.FetchFor(&cf.Item{
						Name: "dummy-service",
						GUID: "dummy-service-guid",
					})
				})

				Expect(ui.Outputs).To(ConsistOf("FAILED", "Unable to get info for service instance dummy-service"))
				Expect(curl.PostItemCallCount()).To(Equal(1))
				calledPath, calledData := curl.PostItemArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/apps"))
				Expect(calledData).To(MatchRegexp(`{"name":"predix-cli-.*","space_guid":"FakeSpaceGUID","memory": 1,"instances": 1}`))
				Expect(curl.PostResourceCallCount()).To(Equal(1))
				calledPath, calledData = curl.PostResourceArgsForCall(0)
				Expect(calledPath).To(Equal("/v2/service_bindings"))
				Expect(calledData).To(Equal(`{"service_instance_guid":"dummy-service-guid","app_guid":"app-guid"}`))
				Expect(curl.DeleteCallCount()).To(Equal(1))
				Expect(curl.DeleteArgsForCall(0)).To(Equal("app-url"))
			})
		})
	})

	Context("print", func() {
		It("prints the info fetched", func() {
			curl.PostItemReturns(&cf.Item{
				GUID: "app-guid",
				URL:  "app-url",
			}, nil)
			resource := &cf.CurlResource{}
			resource.Metadata.URL = "binding-url"
			resource.Entity.Credentials = map[string]interface{}{
				"some-key": "some-value",
				"some-other-key": map[string]interface{}{
					"nested-key": "some-other-value",
				},
			}
			curl.PostResourceReturns(resource, nil)

			testcmd.Run(func() {
				ServiceInfo.PrintFor(&cf.Item{
					Name: "dummy-service",
					GUID: "dummy-service-guid",
				})
			})

			Expect(ui.Outputs).To(ConsistOf("Getting info for service instance dummy-service", "{", "  \"some-key\": \"some-value\",",
				"  \"some-other-key\": {", "    \"nested-key\": \"some-other-value\"", "  }", "}",
				"Note: Depending on the service broker implementation this info may change",
				"Use the 'service-info' command to lookup the binding info for an app and service instance"))
		})
	})

	Context("resolve json path", func() {
		It("returns the value", func() {
			obj := map[string]interface{}{
				"some-key": "some-value",
				"some-other-key": map[string]interface{}{
					"nested-key": "some-other-value",
				},
			}

			Expect(ServiceInfo.ResolveJSONPath(obj, "some-key")).To(Equal("some-value"))
			Expect(ServiceInfo.ResolveJSONPath(obj, "some-other-key/nested-key")).To(Equal("some-other-value"))
		})
	})
})
