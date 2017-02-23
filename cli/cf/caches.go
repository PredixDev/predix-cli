package cf

import (
	"strings"

	. "github.build.ge.com/adoption/predix-cli/cli/cf/constants"
)

const (
	OrgsCache             = iota
	SpacesCache           = iota
	StacksCache           = iota
	ServicesCache         = iota
	ServicePlansCache     = iota
	ServiceInstancesCache = iota
	AppsCache             = iota
)

var Caches = []CacheType{
	CacheType{
		Name:    "orgs",
		Timeout: TimeoutThirtyDays,
		Type:    UserCacheType,
	},
	CacheType{
		Name:    "spaces",
		Timeout: TimeoutThirtyDays,
		Type:    OrgCacheType,
	},
	CacheType{
		Name:    "stacks",
		Timeout: TimeoutThirtyDays,
		Type:    UserCacheType,
	},
	CacheType{
		Name:    "services",
		Timeout: TimeoutThirtyDays,
		Type:    OrgCacheType,
	},
	CacheType{
		Name:    "plans",
		Timeout: TimeoutThirtyDays,
		Type:    OrgCacheType,
	},
	CacheType{
		Name:    "service-instances",
		Timeout: TimeoutOneMin,
		Type:    ServiceInstancesCacheType,
	},
	CacheType{
		Name:    "apps",
		Timeout: TimeoutOneMin,
		Type:    SpaceCacheType,
	},
}

func UpdateCacheTimeout(cacheName string, timeout int) {
	for i := 0; i < len(Caches); i++ {
		if strings.Compare(cacheName, Caches[i].Name) == 0 {
			Caches[i].Timeout = timeout
			break
		}
	}
}
