package cf_test

import (
	"fmt"

	. "github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/cf/cffakes"
	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/global"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lookup", func() {
	var (
		oldCache             CacheInterface
		oldCurl              CurlInterface
		cache                *cffakes.FakeCacheInterface
		curl                 *cffakes.FakeCurlInterface
		generateGetItemsStub func(answers [][]Item) func(path string) []Item
	)

	BeforeEach(func() {
		cache = &cffakes.FakeCacheInterface{}
		curl = &cffakes.FakeCurlInterface{}

		global.Env.NoCache = true

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

		oldCache = Cache
		Cache = cache

		oldCurl = Curl
		Curl = curl

		generateGetItemsStub = func(answers [][]Item) func(path string) []Item {
			return func(path string) []Item {
				if len(answers) == 0 {
					return nil
				}
				answer := answers[0]
				answers = answers[1:]
				return answer
			}
		}
	})

	AfterEach(func() {
		Curl = oldCurl
		Cache = oldCache
	})

	Context("items", func() {
		It("returns fetched items", func() {
			items := []Item{Item{}, Item{}}
			curl.GetItemsReturns(items)

			returnedItems := Lookup.Items("curl-path", Caches[OrgsCache])

			Expect(returnedItems).To(Equal(items))
		})
	})

	Context("item names", func() {
		It("returns names from fetched items", func() {
			items := []Item{Item{Name: "Item1"}, Item{Name: "Item2"}}
			curl.GetItemsReturns(items)

			names := Lookup.ItemNames("curl-path", Caches[OrgsCache])

			Expect(names).To(ConsistOf("Item1", "Item2"))
		})
	})

	Context("orgs", func() {
		It("returns fetched orgs", func() {
			items := []Item{Item{Name: "Org1"}, Item{Name: "Org2"}}
			curl.GetItemsReturns(items)

			names := Lookup.Orgs()

			Expect(names).To(ConsistOf("Org1", "Org2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("/v2/organizations"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			Expect(cacheUsed).To(Equal(Caches[OrgsCache]))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("orgs", func() {
		It("returns fetched orgs", func() {
			items := []Item{Item{Name: "Org1"}, Item{Name: "Org2"}}
			curl.GetItemsReturns(items)

			names := Lookup.Orgs()

			Expect(names).To(ConsistOf("Org1", "Org2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("/v2/organizations"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			Expect(cacheUsed).To(Equal(Caches[OrgsCache]))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("spaces", func() {
		Describe("when no org specified", func() {
			It("returns spaces fetched for current org", func() {
				items := []Item{Item{Name: "Space1"}, Item{Name: "Space2"}}
				curl.GetItemsReturns(items)

				names := Lookup.Spaces(nil)

				Expect(names).To(ConsistOf("Space1", "Space2"))
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/spaces"))
				cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
				Expect(cacheUsed).To(Equal(Caches[SpacesCache]))
				Expect(orgForCache).To(Equal(""))
			})
		})
		Describe("when org is specified", func() {
			It("returns spaces fetched for specified org", func() {
				curl.GetItemsStub = generateGetItemsStub([][]Item{
					[]Item{Item{Name: "Org1", URL: "FakeOrgURL1"}, Item{Name: "Org2", URL: "FakeOrgURL2"}},
					[]Item{Item{Name: "Space1"}, Item{Name: "Space2"}},
				})

				names := Lookup.Spaces(map[string]string{"-o": "Org2"})

				Expect(names).To(ConsistOf("Space1", "Space2"))
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("/v2/organizations"))
				Expect(curl.GetItemsArgsForCall(1)).To(Equal("FakeOrgURL2/spaces"))
				cacheUsed, _ := cache.InitForOrgArgsForCall(0)
				Expect(cacheUsed).To(Equal(Caches[OrgsCache]))
				cacheUsed, orgForCache := cache.InitForOrgArgsForCall(1)
				Expect(cacheUsed).To(Equal(Caches[SpacesCache]))
				Expect(orgForCache).To(Equal("Org2"))
			})
		})
		Describe("when specified org not found", func() {
			It("returns nothing", func() {
				curl.GetItemsStub = generateGetItemsStub([][]Item{
					[]Item{Item{Name: "Org1", URL: "FakeOrgURL1"}, Item{Name: "Org2", URL: "FakeOrgURL2"}},
				})

				names := Lookup.Spaces(map[string]string{"-o": "SomeOrg"})

				Expect(len(names)).To(Equal(0))
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("/v2/organizations"))
				cacheUsed, _ := cache.InitForOrgArgsForCall(0)
				Expect(cacheUsed).To(Equal(Caches[OrgsCache]))
			})
		})
	})

	Context("stacks", func() {
		It("returns fetched stacks", func() {
			items := []Item{Item{Name: "Stack1"}, Item{Name: "Stack2"}}
			curl.GetItemsReturns(items)

			names := Lookup.Stacks()

			Expect(names).To(ConsistOf("Stack1", "Stack2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("/v2/stacks"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			Expect(cacheUsed).To(Equal(Caches[StacksCache]))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("marketplace services", func() {
		It("returns fetched services", func() {
			items := []Item{Item{Name: "Service1"}, Item{Name: "Service2"}}
			curl.GetItemsReturns(items)

			names := Lookup.MarketplaceServices()

			Expect(names).To(ConsistOf("Service1", "Service2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			Expect(cacheUsed).To(Equal(Caches[ServicesCache]))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("marketplace service plans", func() {
		It("returns fetched service plans", func() {
			curl.GetItemsStub = generateGetItemsStub([][]Item{
				[]Item{Item{Name: "some-service", URL: "SomeServiceURL"}},
				[]Item{Item{Name: "Plan1"}, Item{Name: "Plan2"}},
			})

			names := Lookup.MarketplaceServicePlans("some-service")

			Expect(names).To(ConsistOf("Plan1", "Plan2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:some-service"))
			Expect(curl.GetItemsArgsForCall(1)).To(Equal("SomeServiceURL/service_plans"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			actualCache := Caches[ServicePlansCache]
			actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, "some-service")
			Expect(cacheUsed).To(Equal(actualCache))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("marketplace service plan item", func() {
		Describe("when service plan is found", func() {
			It("returns fetched service plan item", func() {
				curl.GetItemsStub = generateGetItemsStub([][]Item{
					[]Item{Item{Name: "some-service", URL: "SomeServiceURL"}},
					[]Item{Item{Name: "Plan1"}, Item{Name: "Plan2"}},
				})

				item := Lookup.MarketplaceServicePlanItem("some-service", "Plan1")

				Expect(item).ToNot(BeNil())
				Expect(item.Name).To(Equal("Plan1"))
			})
		})

		Describe("when service plan is not found", func() {
			It("returns nil", func() {
				item := Lookup.MarketplaceServicePlanItem("some-service", "Plan1")

				Expect(item).To(BeNil())
			})
		})
	})

	Context("service instances", func() {
		It("returns fetched service instances", func() {
			items := []Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}}
			curl.GetItemsReturns(items)

			names := Lookup.ServiceInstances()

			Expect(names).To(ConsistOf("Instance1", "Instance2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeSpaceURL/service_instances"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			Expect(cacheUsed).To(Equal(Caches[ServiceInstancesCache]))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("apps", func() {
		It("returns fetched apps", func() {
			items := []Item{Item{Name: "App1"}, Item{Name: "App2"}}
			curl.GetItemsReturns(items)

			names := Lookup.Apps()

			Expect(names).To(ConsistOf("App1", "App2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeSpaceURL/apps"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			Expect(cacheUsed).To(Equal(Caches[AppsCache]))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("predix uaa instances", func() {
		It("returns fetched predix uaa instances", func() {
			serviceName := constants.PredixUaa

			curl.GetItemsStub = generateGetItemsStub([][]Item{
				[]Item{Item{Name: serviceName, URL: serviceName + "-url"}},
				[]Item{Item{Name: "Plan1", GUID: "GUID1"}, Item{Name: "Plan2", GUID: "GUID2"}},
				[]Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}},
				[]Item{Item{Name: "Instance3"}, Item{Name: "Instance4"}},
			})

			names := Lookup.PredixUaaInstances()

			Expect(names).To(ConsistOf("Instance1", "Instance2", "Instance3", "Instance4"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:" + serviceName))
			Expect(curl.GetItemsArgsForCall(1)).To(Equal(serviceName + "-url/service_plans"))
			Expect(curl.GetItemsArgsForCall(2)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID1"))
			Expect(curl.GetItemsArgsForCall(3)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID2"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			actualCache := Caches[ServiceInstancesCache]
			actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, serviceName)
			Expect(cacheUsed).To(Equal(actualCache))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("predix asset instances", func() {
		It("returns fetched predix asset instances", func() {
			serviceName := constants.PredixAsset

			curl.GetItemsStub = generateGetItemsStub([][]Item{
				[]Item{Item{Name: serviceName, URL: serviceName + "-url"}},
				[]Item{Item{Name: "Plan1", GUID: "GUID1"}, Item{Name: "Plan2", GUID: "GUID2"}},
				[]Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}},
				[]Item{Item{Name: "Instance3"}},
			})

			names := Lookup.PredixAssetInstances()

			Expect(names).To(ConsistOf("Instance1", "Instance2", "Instance3"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:" + serviceName))
			Expect(curl.GetItemsArgsForCall(1)).To(Equal(serviceName + "-url/service_plans"))
			Expect(curl.GetItemsArgsForCall(2)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID1"))
			Expect(curl.GetItemsArgsForCall(3)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID2"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			actualCache := Caches[ServiceInstancesCache]
			actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, serviceName)
			Expect(cacheUsed).To(Equal(actualCache))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("predix timeseries instances", func() {
		It("returns fetched predix timeseries instances", func() {
			serviceName := constants.PredixTimeseries

			curl.GetItemsStub = generateGetItemsStub([][]Item{
				[]Item{Item{Name: serviceName, URL: serviceName + "-url"}},
				[]Item{Item{Name: "Plan1", GUID: "GUID1"}},
				[]Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}},
			})

			names := Lookup.PredixTimeseriesInstances()

			Expect(names).To(ConsistOf("Instance1", "Instance2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:" + serviceName))
			Expect(curl.GetItemsArgsForCall(1)).To(Equal(serviceName + "-url/service_plans"))
			Expect(curl.GetItemsArgsForCall(2)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID1"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			actualCache := Caches[ServiceInstancesCache]
			actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, serviceName)
			Expect(cacheUsed).To(Equal(actualCache))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("predix analytics catalog instances", func() {
		It("returns fetched predix analytics catalog instances", func() {
			serviceName := constants.PredixAnalyticsCatalog

			curl.GetItemsStub = generateGetItemsStub([][]Item{
				[]Item{Item{Name: serviceName, URL: serviceName + "-url"}},
				[]Item{Item{Name: "Plan1", GUID: "GUID1"}, Item{Name: "Plan2", GUID: "GUID2"}},
				[]Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}},
				nil,
			})

			names := Lookup.PredixAnalyticsCatalogInstances()

			Expect(names).To(ConsistOf("Instance1", "Instance2"))
			Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:" + serviceName))
			Expect(curl.GetItemsArgsForCall(1)).To(Equal(serviceName + "-url/service_plans"))
			Expect(curl.GetItemsArgsForCall(2)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID1"))
			Expect(curl.GetItemsArgsForCall(3)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID2"))
			cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
			actualCache := Caches[ServiceInstancesCache]
			actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, serviceName)
			Expect(cacheUsed).To(Equal(actualCache))
			Expect(orgForCache).To(Equal(""))
		})
	})

	Context("predix uaa instance by name", func() {
		Describe("when instance by that name exists", func() {
			It("returns the instance", func() {
				serviceName := constants.PredixUaa

				curl.GetItemsStub = generateGetItemsStub([][]Item{
					[]Item{Item{Name: serviceName, URL: serviceName + "-url"}},
					[]Item{Item{Name: "Plan1", GUID: "GUID1"}, Item{Name: "Plan2", GUID: "GUID2"}},
					[]Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}},
					[]Item{Item{Name: "Instance3"}, Item{Name: "Instance4"}},
				})

				item := Lookup.PredixUaaInstanceItem("Instance2")

				Expect(item).ToNot(BeNil())
				Expect(item.Name).To(Equal("Instance2"))
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:" + serviceName))
				Expect(curl.GetItemsArgsForCall(1)).To(Equal(serviceName + "-url/service_plans"))
				Expect(curl.GetItemsArgsForCall(2)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID1"))
				Expect(curl.GetItemsArgsForCall(3)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID2"))
				cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
				actualCache := Caches[ServiceInstancesCache]
				actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, serviceName)
				Expect(cacheUsed).To(Equal(actualCache))
				Expect(orgForCache).To(Equal(""))
			})
		})
		Describe("when instance by that name does not exist", func() {
			It("returns nil", func() {
				serviceName := constants.PredixUaa

				curl.GetItemsStub = generateGetItemsStub([][]Item{
					[]Item{Item{Name: serviceName, URL: serviceName + "-url"}},
					[]Item{Item{Name: "Plan1", GUID: "GUID1"}, Item{Name: "Plan2", GUID: "GUID2"}},
					[]Item{Item{Name: "Instance1"}, Item{Name: "Instance2"}},
					[]Item{Item{Name: "Instance3"}, Item{Name: "Instance4"}},
				})

				item := Lookup.PredixUaaInstanceItem("SomeInstance")

				Expect(item).To(BeNil())
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeOrgURL/services?q=label:" + serviceName))
				Expect(curl.GetItemsArgsForCall(1)).To(Equal(serviceName + "-url/service_plans"))
				Expect(curl.GetItemsArgsForCall(2)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID1"))
				Expect(curl.GetItemsArgsForCall(3)).To(Equal("FakeSpaceURL/service_instances?q=service_plan_guid:GUID2"))
				cacheUsed, orgForCache := cache.InitForOrgArgsForCall(0)
				actualCache := Caches[ServiceInstancesCache]
				actualCache.Name = fmt.Sprintf("%s-%s", actualCache.Name, serviceName)
				Expect(cacheUsed).To(Equal(actualCache))
				Expect(orgForCache).To(Equal(""))
			})
		})
	})

	Context("lookup instance for instance name", func() {
		Describe("when only one instance by that name exists", func() {
			It("returns the instance", func() {
				curl.GetItemsReturns([]Item{Item{Name: "service-instance"}})

				item := Lookup.InstanceForInstanceName("service-instance")

				Expect(item).ToNot(BeNil())
				Expect(item.Name).To(Equal("service-instance"))
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeSpaceURL/service_instances?q=name:service-instance"))
			})
		})
		Describe("when instance by that name does not exist", func() {
			It("returns nil", func() {
				curl.GetItemsReturns(nil)

				item := Lookup.InstanceForInstanceName("service-instance")

				Expect(item).To(BeNil())
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeSpaceURL/service_instances?q=name:service-instance"))
			})
		})
		Describe("when multiple instances by that name exist", func() {
			It("returns nil", func() {
				curl.GetItemsReturns([]Item{Item{Name: "service-instance-1"}, Item{Name: "service-instance-2"}})

				item := Lookup.InstanceForInstanceName("service-instance")

				Expect(item).To(BeNil())
				Expect(curl.GetItemsArgsForCall(0)).To(Equal("FakeSpaceURL/service_instances?q=name:service-instance"))
			})
		})
	})
})
