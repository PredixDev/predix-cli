package cf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.build.ge.com/adoption/predix-cli/cli/cf/constants"
	"github.build.ge.com/adoption/predix-cli/cli/global"
)

type CacheInterface interface {
	Init(cache CacheType) (cacheDir string)
	InitForOrg(cache CacheType, org string) (cacheDir string)
	Write(cache CacheType, items []Item)
	WriteForOrg(cache CacheType, items []Item, org string)
	Invalidate(cache CacheType)
	InvalidateForOrg(cache CacheType, org string)
	InvalidateType(cache CacheType)
	PurgeCurrent()
	PurgeAll()
	InvalidateOrgs()
	InvalidateSpaces()
	InvalidateStacks()
	InvalidateServiceInstances()
	InvalidateApps()
}

type cache struct{}

var Cache CacheInterface = cache{}

type CacheEntry struct {
	Items   []Item
	Expires time.Time
}

type CacheType struct {
	Name    string
	Timeout int
	Type    int
}

func (o cache) Init(cache CacheType) (cacheDir string) {
	return o.InitForOrg(cache, "")
}

func (o cache) InitForOrg(cache CacheType, org string) (cacheDir string) {
	if global.Env.NoCache {
		return ""
	}

	if CurrentUserInfo().IsValid() && global.Env.ConfigDir != "" {
		cacheDir = filepath.Join(global.Env.ConfigDir, "cf_lookup",
			CurrentUserInfo().GetAPIHash(), CurrentUserInfo().GetNameHash())
		if cache.Type >= constants.OrgCacheType {
			if org != "" {
				cacheDir = filepath.Join(cacheDir, global.Md5Hash(org))
			} else {
				cacheDir = filepath.Join(cacheDir, CurrentUserInfo().GetOrgHash())
			}
		}
		if cache.Type >= constants.SpaceCacheType {
			cacheDir = filepath.Join(cacheDir, CurrentUserInfo().GetSpaceHash())
		}
		if cache.Type >= constants.ServiceInstancesCacheType {
			cacheDir = filepath.Join(cacheDir, "service-instances")
		}
		cacheDirErr := os.MkdirAll(cacheDir, os.FileMode(0700))
		if cacheDirErr != nil {
			cacheDir = ""
		}
	}
	return cacheDir
}

func (o cache) Write(cache CacheType, items []Item) {
	o.WriteForOrg(cache, items, "")
}

func (o cache) WriteForOrg(cache CacheType, items []Item, org string) {
	cacheDir := o.InitForOrg(cache, org)
	if cacheDir != "" {
		cacheEntry := CacheEntry{
			Items:   items,
			Expires: time.Now().Add(time.Duration(cache.Timeout) * time.Second),
		}
		cacheEntryJSON, _ := json.Marshal(cacheEntry)
		_ = ioutil.WriteFile(filepath.Join(cacheDir, cache.Name), cacheEntryJSON, os.FileMode(0700))
	}
}

func (o cache) Invalidate(cache CacheType) {
	o.InvalidateForOrg(cache, "")
}

func (o cache) InvalidateForOrg(cache CacheType, org string) {
	cacheDir := o.InitForOrg(cache, org)
	if cacheDir != "" {
		_ = os.RemoveAll(filepath.Join(cacheDir, cache.Name))
	}
}

func (o cache) InvalidateType(cache CacheType) {
	cacheDir := o.Init(cache)
	if cacheDir != "" {
		_ = os.RemoveAll(filepath.Join(cacheDir))
	}
}

func (o cache) PurgeCurrent() {
	if CurrentUserInfo().IsValid() && global.Env.ConfigDir != "" {
		cacheDir := filepath.Join(global.Env.ConfigDir, "cf_lookup",
			CurrentUserInfo().GetAPIHash(), CurrentUserInfo().GetNameHash())
		_ = os.RemoveAll(cacheDir)
	}
}

func (o cache) PurgeAll() {
	if global.Env.ConfigDir != "" {
		cacheDir := filepath.Join(global.Env.ConfigDir, "cf_lookup")
		_ = os.RemoveAll(cacheDir)
	}
}

func (o cache) InvalidateOrgs() {
	o.Invalidate(Caches[OrgsCache])
}

func (o cache) InvalidateSpaces() {
	o.Invalidate(Caches[SpacesCache])
}

func (o cache) InvalidateStacks() {
	o.Invalidate(Caches[StacksCache])
}

func (o cache) InvalidateServiceInstances() {
	o.InvalidateType(Caches[ServiceInstancesCache])
}

func (o cache) InvalidateApps() {
	o.Invalidate(Caches[AppsCache])
}
