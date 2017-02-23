package cf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/PredixDev/go-uaa-lib"
)

type ConfigData struct {
	Target             string
	AccessToken        string
	OrganizationFields struct {
		GUID string
		Name string
	}
	SpaceFields struct {
		GUID string
		Name string
	}
}

func CurrentUserInfo() *global.UserInfo {
	if global.CurrentUserInfo == nil {
		global.CurrentUserInfo = &global.UserInfo{}
		cfConfigFilePath := filepath.Join(global.Env.CfHomeDir, ".cf", "config.json")
		data, err := ioutil.ReadFile(cfConfigFilePath)
		if err == nil {
			var cfConfigData ConfigData
			err = json.Unmarshal(data, &cfConfigData)
			if err != nil {
				global.UI.Failed(fmt.Sprintf("The CF Config file is invalid: %s", cfConfigFilePath))
			}

			global.CurrentUserInfo.API = cfConfigData.Target

			cfConfigData.AccessToken = strings.Replace(cfConfigData.AccessToken, "bearer ", "", 1)
			tc, err := lib.TokenClaimsFactory.New(cfConfigData.AccessToken)
			if err != nil {
				global.UI.Failed(fmt.Sprintf("The CF Config file has an invalid access token: %s", cfConfigFilePath))
			}
			global.CurrentUserInfo.Name = tc.UserName

			if cfConfigData.OrganizationFields.GUID != "" {
				global.CurrentUserInfo.Org = cfConfigData.OrganizationFields.Name
				global.CurrentUserInfo.OrgGUID = cfConfigData.OrganizationFields.GUID
				global.CurrentUserInfo.OrgURL = fmt.Sprintf("/v2/organizations/%s", cfConfigData.OrganizationFields.GUID)
			}

			if cfConfigData.SpaceFields.GUID != "" {
				global.CurrentUserInfo.Space = cfConfigData.SpaceFields.Name
				global.CurrentUserInfo.SpaceGUID = cfConfigData.SpaceFields.GUID
				global.CurrentUserInfo.SpaceURL = fmt.Sprintf("/v2/spaces/%s", cfConfigData.SpaceFields.GUID)
			}
		}
	}
	return global.CurrentUserInfo
}
