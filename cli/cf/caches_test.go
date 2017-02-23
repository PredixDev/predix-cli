package cf_test

import (
	. "github.build.ge.com/adoption/predix-cli/cli/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Caches", func() {
	It("sets the cache timeout", func() {
		UpdateCacheTimeout("orgs", 100)
		Expect(Caches[OrgsCache].Timeout).To(Equal(100))
		UpdateCacheTimeout("spaces", 130)
		Expect(Caches[SpacesCache].Timeout).To(Equal(130))
		UpdateCacheTimeout("stacks", 200)
		Expect(Caches[StacksCache].Timeout).To(Equal(200))
		UpdateCacheTimeout("services", 500)
		Expect(Caches[ServicesCache].Timeout).To(Equal(500))
		UpdateCacheTimeout("plans", 360)
		Expect(Caches[ServicePlansCache].Timeout).To(Equal(360))
		UpdateCacheTimeout("service-instances", 210)
		Expect(Caches[ServiceInstancesCache].Timeout).To(Equal(210))
		UpdateCacheTimeout("apps", 4000)
		Expect(Caches[AppsCache].Timeout).To(Equal(4000))
	})
})
