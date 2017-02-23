package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/satori/go.uuid"
	"github.com/yasuyuky/jsonpath"
)

type ServiceInfoInterface interface {
	FetchFor(serviceInstance *cf.Item) map[string]interface{}
	FetchForAppAndServiceInstance(appInstance *cf.Item, serviceInstance *cf.Item) map[string]interface{}
	PrintFor(serviceInstance *cf.Item)
	PrintForAppAndServiceInstance(appInstance *cf.Item, serviceInstance *cf.Item)
	ResolveJSONPath(obj interface{}, path string) interface{}
}

type serviceInfo struct{}

var ServiceInfo ServiceInfoInterface = serviceInfo{}

func (o serviceInfo) FetchFor(serviceInstance *cf.Item) map[string]interface{} {
	appName := fmt.Sprintf("predix-cli-%s", uuid.NewV4())
	appInstance, err := cf.Curl.PostItem("/v2/apps", fmt.Sprintf(`{"name":"%s","space_guid":"%s","memory": 1,"instances": 1}`,
		appName, cf.CurrentUserInfo().SpaceGUID))

	if err != nil {
		global.UI.Failed("Unable to get info for service instance %s", terminal.EntityNameColor(serviceInstance.Name))
	}
	defer cf.Curl.Delete(appInstance.URL)

	binding, err := cf.Curl.PostResource("/v2/service_bindings", fmt.Sprintf(`{"service_instance_guid":"%s","app_guid":"%s"}`,
		serviceInstance.GUID, appInstance.GUID))
	if err != nil {
		global.UI.Failed("Unable to get info for service instance %s", terminal.EntityNameColor(serviceInstance.Name))
	}
	defer cf.Curl.Delete(binding.Metadata.URL)

	return binding.Entity.Credentials
}

func (o serviceInfo) FetchForAppAndServiceInstance(appInstance *cf.Item, serviceInstance *cf.Item) map[string]interface{} {
	bindings := cf.Curl.GetResources(fmt.Sprintf("%s/service_bindings", appInstance.URL))
	if bindings != nil {
		for _, binding := range bindings {
			if binding.Entity.ServiceInstanceURL == serviceInstance.URL {
				return binding.Entity.Credentials
			}
		}
	}

	global.UI.Failed("Unable to find binding between app %s and service instance %s",
		terminal.EntityNameColor(appInstance.Name), terminal.EntityNameColor(serviceInstance.Name))
	return nil
}

func (o serviceInfo) PrintForAppAndServiceInstance(appInstance *cf.Item, serviceInstance *cf.Item) {
	global.UI.Say("Getting info for app %s and service instance %s",
		terminal.EntityNameColor(appInstance.Name), terminal.EntityNameColor(serviceInstance.Name))
	info := o.FetchForAppAndServiceInstance(appInstance, serviceInstance)
	infoJSON, _ := json.MarshalIndent(info, "", "  ")
	global.UI.Say(string(infoJSON))
}

func (o serviceInfo) PrintFor(serviceInstance *cf.Item) {
	global.UI.Say("Getting info for service instance %s", terminal.EntityNameColor(serviceInstance.Name))
	info := o.FetchFor(serviceInstance)
	infoJSON, _ := json.MarshalIndent(info, "", "  ")
	global.UI.Say(string(infoJSON))
	global.UI.Say("Note: Depending on the service broker implementation this info may change")
	global.UI.Say("Use the 'service-info' command to lookup the binding info for an app and service instance")
}

func (o serviceInfo) ResolveJSONPath(obj interface{}, path string) interface{} {
	pathArray := strings.Split(path, "/")
	pathInterfaceArray := make([]interface{}, len(pathArray))
	for i := 0; i < len(pathArray); i++ {
		pathInterfaceArray[i] = pathArray[i]
	}
	value, _ := jsonpath.Get(obj, pathInterfaceArray, nil)
	return value
}
