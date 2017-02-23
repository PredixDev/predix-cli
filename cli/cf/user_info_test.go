package cf_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	testterm "github.build.ge.com/adoption/cli-lib/testhelpers/terminal"
	testcmd "github.build.ge.com/adoption/predix-cli/cli/testhelpers/commands"

	. "github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"

	"github.com/PredixDev/go-uaa-lib"
	"github.com/PredixDev/go-uaa-lib/libfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Info", func() {
	var (
		ui                    *testterm.FakeUI
		oldTokenClaimsFactory lib.TokenClaimsFactoryInterface
		tokenClaimsFactory    *libfakes.FakeTokenClaimsFactoryInterface
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}

		global.UI = ui
		global.Env.CfHomeDir, _ = ioutil.TempDir("", "test-cf-home")
		_ = os.MkdirAll(filepath.Join(global.Env.CfHomeDir, ".cf"), os.FileMode(0700))

		oldTokenClaimsFactory = lib.TokenClaimsFactory
		tokenClaimsFactory = &libfakes.FakeTokenClaimsFactoryInterface{}
		lib.TokenClaimsFactory = tokenClaimsFactory

		global.CurrentUserInfo = nil
	})
	AfterEach(func() {
		global.CurrentUserInfo = nil

		lib.TokenClaimsFactory = oldTokenClaimsFactory

		_ = os.RemoveAll(filepath.Join(global.Env.CfHomeDir, ".cf"))
		_ = os.RemoveAll(global.Env.CfHomeDir)
		global.Env.CfHomeDir = ""
	})

	Describe("when the cf config file is valid", func() {
		It("loads the info", func() {
			configData := ConfigData{}
			configData.Target = "CF Target"
			configData.AccessToken = "CF Access Token"
			configData.OrganizationFields.GUID = "Org GUID"
			configData.OrganizationFields.Name = "Org Name"
			configData.SpaceFields.GUID = "Space GUID"
			configData.SpaceFields.Name = "Space Name"

			configDataJSON, _ := json.MarshalIndent(configData, "", "  ")
			err := ioutil.WriteFile(filepath.Join(global.Env.CfHomeDir, ".cf", "config.json"), configDataJSON, os.FileMode(0700))
			Expect(err).To(BeNil())

			tokenClaimsFactory.NewReturns(&lib.TokenClaims{
				UserName: "CF User",
			}, nil)

			userInfo := CurrentUserInfo()

			Expect(userInfo.API).To(Equal("CF Target"))
			Expect(userInfo.Name).To(Equal("CF User"))
			Expect(userInfo.Org).To(Equal("Org Name"))
			Expect(userInfo.OrgGUID).To(Equal("Org GUID"))
			Expect(userInfo.OrgURL).To(Equal("/v2/organizations/Org GUID"))
			Expect(userInfo.Space).To(Equal("Space Name"))
			Expect(userInfo.SpaceGUID).To(Equal("Space GUID"))
			Expect(userInfo.SpaceURL).To(Equal("/v2/spaces/Space GUID"))
		})
	})

	Describe("when the cf config file is invalid", func() {
		It("shows an error", func() {
			err := ioutil.WriteFile(filepath.Join(global.Env.CfHomeDir, ".cf", "config.json"), []byte("Invalid JSON"), os.FileMode(0700))
			Expect(err).To(BeNil())

			testcmd.Run(func() {
				CurrentUserInfo()
			})

			Expect(ui.Outputs).To(ConsistOf("FAILED", MatchRegexp("The CF Config file is invalid:.*")))
		})
	})

	Describe("when the cf token is valid", func() {
		It("shows an error", func() {
			configData := ConfigData{}
			configDataJSON, _ := json.MarshalIndent(configData, "", "  ")
			err := ioutil.WriteFile(filepath.Join(global.Env.CfHomeDir, ".cf", "config.json"), configDataJSON, os.FileMode(0700))
			Expect(err).To(BeNil())

			tokenClaimsFactory.NewReturns(nil, fmt.Errorf("Error"))

			testcmd.Run(func() {
				CurrentUserInfo()
			})

			Expect(ui.Outputs).To(ConsistOf("FAILED", MatchRegexp("The CF Config file has an invalid access token:.*")))
		})
	})
})
