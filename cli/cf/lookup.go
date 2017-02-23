package cf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
)

type LookupInterface interface {
	Orgs() (orgs []string)
	Spaces(params map[string]string) []string
	Stacks() []string
	MarketplaceServices() []string
	MarketplaceServicePlans(serviceName string) []string
	ServiceInstances() []string
	Apps() []string
	MarketplaceServicePlansItems(serviceName string) []Item
	MarketplaceServicePlanItem(serviceName string, planName string) *Item
	InstancesForService(serviceName string) []string
	InstancesItemsForService(serviceName string) []Item
	InstanceForInstanceName(serviceInstance string) (instance *Item)
	AppForName(appName string) (instance *Item)
	PredixUaaInstances() []string
	PredixAssetInstances() []string
	PredixTimeseriesInstances() []string
	PredixAnalyticsCatalogInstances() []string
	PredixUaaInstanceItem(instanceName string) *Item

	ItemNames(path string, cache CacheType) []string
	ItemNamesForOrg(path string, cache CacheType, org string) []string
	Items(path string, cache CacheType) []Item
	ItemsForOrg(path string, cache CacheType, org string) []Item
	NameFromItems(cfItems []Item) (items []string)

	InCache(cache CacheType) (items []Item)
	InCacheForOrg(cache CacheType, org string) (items []Item)
}

type lookup struct{}

var Lookup LookupInterface = lookup{}

func (o lookup) Orgs() (orgs []string) {
	return o.ItemNames("/v2/organizations", Caches[OrgsCache])
}

func (o lookup) Spaces(params map[string]string) []string {
	if params == nil || params["-o"] == "" {
		return o.ItemNames(fmt.Sprintf("%s/spaces", CurrentUserInfo().OrgURL), Caches[SpacesCache])
	}
	org := params["-o"]
	orgItems := o.Items("/v2/organizations", Caches[OrgsCache])
	for _, orgItem := range orgItems {
		if strings.Compare(orgItem.Name, org) == 0 {
			return o.ItemNamesForOrg(fmt.Sprintf("%s/spaces", orgItem.URL), Caches[SpacesCache], org)
		}
	}
	return []string{}
}

func (o lookup) Stacks() []string {
	return o.ItemNames("/v2/stacks", Caches[StacksCache])
}

func (o lookup) MarketplaceServices() []string {
	return o.ItemNames(fmt.Sprintf("%s/services", CurrentUserInfo().OrgURL), Caches[ServicesCache])
}

func (o lookup) MarketplaceServicePlans(serviceName string) []string {
	return o.NameFromItems(o.MarketplaceServicePlansItems(serviceName))
}

func (o lookup) ServiceInstances() []string {
	return o.ItemNames(fmt.Sprintf("%s/service_instances", CurrentUserInfo().SpaceURL), Caches[ServiceInstancesCache])
}

func (o lookup) Apps() []string {
	return o.ItemNames(fmt.Sprintf("%s/apps", CurrentUserInfo().SpaceURL), Caches[AppsCache])
}

func (o lookup) MarketplaceServicePlansItems(serviceName string) []Item {
	cache := Caches[ServicePlansCache]
	cache.Name = fmt.Sprintf("%s-%s", cache.Name, serviceName)

	plans := o.InCache(cache)
	if plans == nil {
		services := o.InCache(Caches[ServicesCache])
		if services == nil {
			services = Curl.GetItems(fmt.Sprintf("%s/services?q=label:%s", CurrentUserInfo().OrgURL, serviceName))
		}
		for _, service := range services {
			if strings.Compare(service.Name, serviceName) == 0 {
				plans = Curl.GetItems(fmt.Sprintf("%s/service_plans", service.URL))
				if plans != nil {
					Cache.Write(cache, plans)
				}
				break
			}
		}
	}
	return plans
}

func (o lookup) MarketplaceServicePlanItem(serviceName string, planName string) *Item {
	plans := o.MarketplaceServicePlansItems(serviceName)
	for _, plan := range plans {
		if strings.Compare(plan.Name, planName) == 0 {
			return &plan
		}
	}
	return nil
}

func (o lookup) InstancesForService(serviceName string) []string {
	return o.NameFromItems(o.InstancesItemsForService(serviceName))
}

func (o lookup) InstancesItemsForService(serviceName string) []Item {
	cache := Caches[ServiceInstancesCache]
	cache.Name = fmt.Sprintf("%s-%s", cache.Name, serviceName)

	serviceInstances := o.InCache(cache)
	if serviceInstances == nil {
		servicePlans := o.MarketplaceServicePlansItems(serviceName)
		serviceInstancesPath := fmt.Sprintf("%s/service_instances", CurrentUserInfo().SpaceURL)
		serviceInstances = []Item{}
		for _, plan := range servicePlans {
			planInstances := Curl.GetItems(fmt.Sprintf("%s?q=service_plan_guid:%s", serviceInstancesPath, plan.GUID))
			if planInstances != nil {
				serviceInstances = append(serviceInstances, planInstances...)
			}
		}
		if len(serviceInstances) > 0 {
			Cache.Write(cache, serviceInstances)
		}
	}
	return serviceInstances
}

func (o lookup) PredixUaaInstances() []string {
	return o.InstancesForService(constants.PredixUaa)
}

func (o lookup) PredixAssetInstances() []string {
	return o.InstancesForService(constants.PredixAsset)
}

func (o lookup) PredixTimeseriesInstances() []string {
	return o.InstancesForService(constants.PredixTimeseries)
}

func (o lookup) PredixAnalyticsCatalogInstances() []string {
	return o.InstancesForService(constants.PredixAnalyticsCatalog)
}

func (o lookup) PredixUaaInstanceItem(instanceName string) *Item {
	instances := o.InstancesItemsForService(constants.PredixUaa)
	for _, instance := range instances {
		if instance.Name == instanceName {
			return &instance
		}
	}
	return nil
}

func (o lookup) AppForName(appName string) (instance *Item) {
	instance = nil
	appInstances := Curl.GetItems(fmt.Sprintf("%s/apps?q=name:%s", CurrentUserInfo().SpaceURL, appName))
	if len(appInstances) == 1 {
		instance = &appInstances[0]
	}
	return instance
}

func (o lookup) InstanceForInstanceName(serviceInstance string) (instance *Item) {
	instance = nil
	serviceInstances := Curl.GetItems(fmt.Sprintf("%s/service_instances?q=name:%s", CurrentUserInfo().SpaceURL, serviceInstance))
	if len(serviceInstances) == 1 {
		instance = &serviceInstances[0]
	}
	return instance
}

func (o lookup) ItemNames(path string, cache CacheType) []string {
	return o.NameFromItems(o.Items(path, cache))
}

func (o lookup) ItemNamesForOrg(path string, cache CacheType, org string) []string {
	return o.NameFromItems(o.ItemsForOrg(path, cache, org))
}

func (o lookup) Items(path string, cache CacheType) []Item {
	return o.ItemsForOrg(path, cache, "")
}

func (o lookup) ItemsForOrg(path string, cache CacheType, org string) []Item {
	items := o.InCacheForOrg(cache, org)
	if items != nil {
		return items
	}
	items = Curl.GetItems(path)
	if items != nil {
		Cache.WriteForOrg(cache, items, org)
	}
	return items
}

func (o lookup) NameFromItems(cfItems []Item) (items []string) {
	if cfItems == nil {
		return []string{}
	}
	items = make([]string, len(cfItems))

	for i := range cfItems {
		items[i] = cfItems[i].Name
	}
	return items
}
func (o lookup) InCache(cache CacheType) (items []Item) {
	return o.InCacheForOrg(cache, "")
}

func (o lookup) InCacheForOrg(cache CacheType, org string) (items []Item) {
	items = nil
	cacheDir := Cache.InitForOrg(cache, org)
	if cacheDir != "" {
		cacheEntryJSON, err := ioutil.ReadFile(filepath.Join(cacheDir, cache.Name))
		if err == nil {
			var cacheEntry CacheEntry
			err := json.Unmarshal(cacheEntryJSON, &cacheEntry)
			if err == nil && time.Now().Before(cacheEntry.Expires) {
				items = cacheEntry.Items
			} else {
				Cache.InvalidateForOrg(cache, org)
			}
		}
	}
	return items
}
