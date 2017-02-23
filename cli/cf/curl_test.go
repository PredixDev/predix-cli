package cf_test

import (
	"fmt"

	testshell "github.build.ge.com/adoption/cli-lib/testhelpers/shell"
	. "github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Curl", func() {
	var (
		sh          *testshell.FakeShell
		guidCount   int
		resourceStr func(string, string, bool) string
		fakeOutput  func(string, error) testshell.FakeOutput
	)

	BeforeEach(func() {
		sh = &testshell.FakeShell{}

		global.Env.NoCache = true
		global.Sh = sh

		guidCount = 1
		resourceStr = func(name, resourceType string, more bool) string {
			defer func() { guidCount = guidCount + 1 }()
			fmtStr := `{
					"metadata": {
						"guid": "00000000-0000-0000-0000-000000%02d",
						"url": "/v2/%s/00000000-0000-0000-0000-000000%02d"
					},
					"entity": {
						"name": "%s"
					}
				}`
			if more {
				fmtStr = fmtStr + ","
			}
			return fmt.Sprintf(fmtStr, guidCount, resourceType, guidCount, name)
		}
		fakeOutput = func(data string, err error) testshell.FakeOutput {
			return testshell.FakeOutput{
				Data: []byte(data),
				Err:  err,
			}
		}
	})

	Context("get resources", func() {
		Describe("when resources fetched successfully", func() {
			It("returns the resources", func() {
				sh.UnmarshalOutputs = []string{`{"resources":[` + resourceStr("Org1", "orgs", true) + resourceStr("Org2", "orgs", false) + `]}`}

				resources := Curl.GetResources("curl-path")

				Expect(resources).ToNot(BeNil())
				Expect(len(resources)).To(Equal(2))
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path"))
				Expect(sh.UnmarshalTypes).To(ConsistOf("*cf.CurlResponse"))
			})
		})
		Describe("when resources fetch fails", func() {
			It("returns nil", func() {
				sh.UnmarshalOutputs = []string{`Invalid JSON`}

				resources := Curl.GetResources("curl-path")

				Expect(resources).To(BeNil())
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path"))
				Expect(sh.UnmarshalTypes).To(ConsistOf("*cf.CurlResponse"))
			})
		})
		Describe("when resources fetch returns an error", func() {
			It("returns nil", func() {
				sh.UnmarshalOutputs = []string{`{"error_code": "123", "description": "Failed"}`}

				resources := Curl.GetResources("curl-path")

				Expect(resources).To(BeNil())
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path"))
				Expect(sh.UnmarshalTypes).To(ConsistOf("*cf.CurlResponse"))
			})
		})
	})

	Context("get resource", func() {
		Describe("when resource fetched successfully", func() {
			It("returns the resource", func() {
				sh.UnmarshalOutputs = []string{resourceStr("Org1", "orgs", false)}

				resource := Curl.GetResource("curl-path")

				Expect(resource).ToNot(BeNil())
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path"))
				Expect(sh.UnmarshalTypes).To(ConsistOf("*cf.CurlResource"))
			})
		})
		Describe("when resource fetch fails", func() {
			It("returns nil", func() {
				sh.UnmarshalOutputs = []string{`Invalid JSON`}

				resource := Curl.GetResource("curl-path")

				Expect(resource).To(BeNil())
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path"))
				Expect(sh.UnmarshalTypes).To(ConsistOf("*cf.CurlResource"))
			})
		})
		Describe("when resource fetch returns an error", func() {
			It("returns nil", func() {
				sh.UnmarshalOutputs = []string{`{"error_code": "123", "description": "Failed"}`}

				resource := Curl.GetResource("curl-path")

				Expect(resource).To(BeNil())
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path"))
				Expect(sh.UnmarshalTypes).To(ConsistOf("*cf.CurlResource"))
			})
		})
	})

	Context("get item from resource", func() {
		It("uses name if present", func() {
			resource := &CurlResource{}
			resource.Metadata.GUID = "ABCD"
			resource.Metadata.URL = "Some-URL"
			resource.Entity.Name = "Some-Name"

			item := Curl.GetItemFromResource(resource)

			Expect(item).ToNot(BeNil())
			Expect(item.GUID).To(Equal(resource.Metadata.GUID))
			Expect(item.URL).To(Equal(resource.Metadata.URL))
			Expect(item.Name).To(Equal(resource.Entity.Name))
		})
		Describe("when name not present", func() {
			It("uses label if present", func() {
				resource := &CurlResource{}
				resource.Metadata.GUID = "ABCD"
				resource.Metadata.URL = "Some-URL"
				resource.Entity.Label = "Some-Label"

				item := Curl.GetItemFromResource(resource)

				Expect(item).ToNot(BeNil())
				Expect(item.GUID).To(Equal(resource.Metadata.GUID))
				Expect(item.URL).To(Equal(resource.Metadata.URL))
				Expect(item.Name).To(Equal(resource.Entity.Label))
			})
		})
	})

	Context("get items", func() {
		Describe("when resources found", func() {
			It("returns the items", func() {
				sh.UnmarshalOutputs = []string{`{"resources":[` + resourceStr("Org1", "orgs", true) + resourceStr("Org2", "orgs", false) + `]}`}
				items := Curl.GetItems("curl-path")
				Expect(items).ToNot(BeNil())
				Expect(len(items)).To(Equal(2))
			})
		})
		Describe("when resources empty", func() {
			It("returns empty items", func() {
				sh.UnmarshalOutputs = []string{`{"resources":[]}`}
				items := Curl.GetItems("curl-path")
				Expect(items).ToNot(BeNil())
				Expect(len(items)).To(Equal(0))
			})
		})
		Describe("when resources not found", func() {
			It("returns nil", func() {
				sh.UnmarshalOutputs = []string{`Invalid JSON`}
				items := Curl.GetItems("curl-path")
				Expect(items).To(BeNil())
			})
		})
	})

	Context("get item", func() {
		Describe("when resource found", func() {
			It("returns the item", func() {
				sh.UnmarshalOutputs = []string{resourceStr("Org1", "orgs", false)}
				item := Curl.GetItem("curl-path")
				Expect(item).ToNot(BeNil())
			})
		})
		Describe("when resource not found", func() {
			It("returns nil", func() {
				sh.UnmarshalOutputs = []string{`Invalid JSON`}
				item := Curl.GetItem("curl-path")
				Expect(item).To(BeNil())
			})
		})
	})

	Context("post resource", func() {
		Describe("when resource posted successfully", func() {
			It("returns the resource", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput(resourceStr("Org1", "orgs", false), nil)}

				resource, err := Curl.PostResource("curl-path", "post data")

				Expect(resource).ToNot(BeNil())
				Expect(err).To(BeNil())
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path", "-X", "POST", "-d", "'post data'"))
			})
		})
		Describe("when resource post fails", func() {
			It("returns the error", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput("", fmt.Errorf("Failed"))}

				_, err := Curl.PostResource("curl-path", "post data")

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed"))
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path", "-X", "POST", "-d", "'post data'"))
			})
		})
		Describe("when resource post returns an error", func() {
			It("returns the error", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput(`{"error_code": "123", "description": "Failed"}`, nil)}

				_, err := Curl.PostResource("curl-path", "post data")

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed"))
				Expect(sh.EnvProps).To(ConsistOf("CF_TRACE", "false"))
				Expect(sh.Commands).To(ConsistOf("cf", "curl", "curl-path", "-X", "POST", "-d", "'post data'"))
			})
		})
	})

	Context("post item", func() {
		Describe("when resource posted", func() {
			It("returns the item", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput(resourceStr("Org1", "orgs", false), nil)}
				item, err := Curl.PostItem("curl-path", "post data")
				Expect(item).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
		Describe("when resource not posted", func() {
			It("returns nil", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput("", fmt.Errorf("Failed"))}
				item, err := Curl.PostItem("curl-path", "post data")
				Expect(item).To(BeNil())
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed"))
			})
		})
	})

	Context("post item", func() {
		Describe("when delete succeeds", func() {
			It("returns nil", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput("", nil)}
				err := Curl.Delete("curl-path")
				Expect(err).To(BeNil())
			})
		})
		Describe("when delete fails", func() {
			It("returns error", func() {
				sh.Outputs = []testshell.FakeOutput{fakeOutput("", fmt.Errorf("Failed"))}
				err := Curl.Delete("curl-path")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed"))
			})
		})
	})
})
