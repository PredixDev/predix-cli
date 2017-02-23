package cf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.build.ge.com/adoption/predix-cli/cli/global"
)

type Item struct {
	Name string
	GUID string
	URL  string
}

type CurlResource struct {
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"description"`
	Metadata         struct {
		GUID string `json:"guid"`
		URL  string `json:"url"`
	} `json:"metadata"`
	Entity struct {
		Name               string                 `json:"name"`
		Label              string                 `json:"label"`
		Credentials        map[string]interface{} `json:"credentials"`
		ServiceURL         string                 `json:"service_url"`
		ServicePlanURL     string                 `json:"service_plan_url"`
		ServiceInstanceURL string                 `json:"service_instance_url"`
	} `json:"entity"`
}

type CurlResponse struct {
	ErrorCode        string         `json:"error_code"`
	ErrorDescription string         `json:"description"`
	TotalResults     int            `json:"total_results"`
	Resources        []CurlResource `json:"resources"`
}

type CurlInterface interface {
	GetItemFromResource(resource *CurlResource) (item *Item)
	GetItems(path string) (items []Item)
	GetResources(path string) (resources []CurlResource)
	GetResource(path string) (resource *CurlResource)
	GetItem(path string) (item *Item)
	PostItem(path string, data string) (item *Item, err error)
	PostResource(path string, data string) (r *CurlResource, err error)
	Delete(path string) error
}

type curl struct{}

var Curl CurlInterface = curl{}

func (o curl) GetItemFromResource(resource *CurlResource) (item *Item) {
	item = &Item{
		GUID: resource.Metadata.GUID,
		URL:  resource.Metadata.URL,
	}
	if resource.Entity.Name != "" {
		item.Name = resource.Entity.Name
	} else {
		item.Name = resource.Entity.Label
	}
	return item
}

func (o curl) GetItems(path string) (items []Item) {
	resources := o.GetResources(path)
	if resources != nil {
		items = make([]Item, len(resources))
		for i := range resources {
			items[i] = *o.GetItemFromResource(&resources[i])
		}
	} else {
		items = nil
	}
	return items
}

func (o curl) GetResources(path string) (resources []CurlResource) {
	var r CurlResponse
	err := global.Sh.SetEnv("CF_TRACE", "false").Command("cf", "curl", path).Unmarshal(&r)
	if err == nil && r.ErrorCode == "" {
		return r.Resources
	}
	return nil
}

func (o curl) GetResource(path string) (resource *CurlResource) {
	resource = &CurlResource{}
	err := global.Sh.SetEnv("CF_TRACE", "false").Command("cf", "curl", path).Unmarshal(resource)
	if err == nil && resource.ErrorCode == "" {
		return resource
	}
	return nil
}

func (o curl) GetItem(path string) (item *Item) {
	r := o.GetResource(path)
	item = nil
	if r != nil {
		item = o.GetItemFromResource(r)
	}
	return item
}

func (o curl) PostItem(path string, data string) (item *Item, err error) {
	r, err := o.PostResource(path, data)
	item = nil
	if err == nil {
		item = o.GetItemFromResource(r)
	}
	return item, err
}

func (o curl) PostResource(path string, data string) (r *CurlResource, err error) {
	out, err := global.Sh.SetEnv("CF_TRACE", "false").Command("cf", "curl", path, "-X", "POST", "-d", fmt.Sprintf("'%s'", data)).Output()
	r = &CurlResource{}
	if err == nil {
		err = json.Unmarshal(out, r)
		if err == nil && r.ErrorCode != "" {
			err = errors.New(r.ErrorDescription)
		}
	}
	return r, err
}

func (o curl) Delete(path string) error {
	_, err := global.Sh.SetEnv("CF_TRACE", "false").Command("cf", "curl", path, "-X", "DELETE").Output()
	return err
}
