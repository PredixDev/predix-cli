package uaac

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.build.ge.com/adoption/cli-lib/terminal"
	"github.build.ge.com/adoption/predix-cli/cli/cf"
	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/PredixDev/go-uaa-lib"
)

type TargetsInterface interface {
	PrintAll()
	PrintCurrent()
	SetCurrentForID(id int)
	SetCurrent(url string, cfInstanceURL string, skipSslVerify bool, caCertFile string, tr *lib.TokenResponse)
	LookupAndSetCurrent(url string, context string) bool
	GetCurrent() (target *Target, context *Context, instance *cf.Item)
	LoadConfig()
}

type targets struct {
	list []*Target
}

var Targets TargetsInterface = &targets{}

type Target struct {
	TargetURL      string              `json:"url,omitempty"`
	CfInstanceURL  string              `json:"cf_url,omitempty"`
	CaCertFilePath string              `json:"ca_cert,omitempty"`
	SkipSsl        bool                `json:"skip_ssl_validation,omitempty"`
	Current        bool                `json:"current,omitempty"`
	Contexts       map[string]*Context `json:"contexts,omitempty"`
}

type Context struct {
	ID       string   `json:"jti,omitempty"`
	UserID   string   `json:"user_id,omitempty"`
	UserName string   `json:"user_name,omitempty"`
	ClientID string   `json:"client_id,omitempty"`
	Access   string   `json:"access_token,omitempty"`
	Refresh  string   `json:"refresh_token,omitempty"`
	Type     string   `json:"token_type,omitempty"`
	Scopes   []string `json:"scope,omitempty"`
	Current  bool     `json:"current,omitempty"`
}

func (t *Target) URL() string {
	return t.TargetURL
}

func (t *Target) CaCertFile() string {
	return t.CaCertFilePath
}

func (t *Target) SkipSslVerify() bool {
	return t.SkipSsl
}

func (c *Context) AccessToken() string {
	return c.Access
}

func (c *Context) TokenType() string {
	return c.Type
}

func (o *targets) PrintAll() {
	if len(o.list) == 0 {
		global.UI.Say("No UAA targets")
		return
	}

	count := 0
	fetchError := false
	table := global.UI.Table([]string{terminal.HeaderColor("ID"), "", terminal.HeaderColor("Target"), terminal.HeaderColor("Context")})

	for _, t := range o.list {
		serviceInstance := cf.Curl.GetItem(t.CfInstanceURL)

		if serviceInstance == nil {
			fetchError = true
			continue
		}

		if len(t.Contexts) == 0 {
			marker := targetMarker(t)

			count++
			table.Add(strconv.Itoa(count), marker, terminal.EntityNameColor(serviceInstance.Name), "No context")
		} else {
			for _, c := range t.Contexts {
				context := readContext(c)
				marker := contextMarker(t, c)

				count++
				table.Add(strconv.Itoa(count), marker, terminal.EntityNameColor(serviceInstance.Name), context)
			}
		}
	}

	if fetchError {
		global.UI.Say("Failed to fetch one or more target UAA's info")
	}

	if count > 0 {
		table.Print()
	}
}

func targetMarker(t *Target) string {
	if t.Current {
		return "*"
	}
	return ""
}

func readContext(c *Context) string {
	if c.UserID != "" {
		return fmt.Sprintf("User: %s, Client: %s", c.UserName, c.ClientID)
	}
	return fmt.Sprintf("Client: %s", c.ClientID)
}

func contextMarker(t *Target, c *Context) string {
	if t.Current && c.Current {
		return "*"
	}
	return ""
}

func (o *targets) PrintCurrent() {
	var t *Target
	for _, v := range o.list {
		if v.Current {
			t = v
			break
		}
	}

	if t == nil {
		global.UI.Say("No UAA target set")
	} else {
		serviceInstance := cf.Curl.GetItem(t.CfInstanceURL)
		table := global.UI.Table([]string{"", ""})

		name := ""
		if serviceInstance != nil {
			name = serviceInstance.Name
		} else {
			global.UI.Say("Failed to fetch target UAA's info")
		}
		table.Add(terminal.EntityNameColor("Target:"), terminal.EntityNameColor(name))
		table.Add(terminal.EntityNameColor("URL:"), terminal.EntityNameColor(t.TargetURL))

		var c *Context
		for _, v := range t.Contexts {
			if v.Current {
				c = v
				break
			}
		}

		if c == nil {
			global.UI.Say("No context set")
		} else {
			if c.UserName != "" {
				table.Add(terminal.EntityNameColor("User:"), terminal.EntityNameColor(c.UserName))
			}
			table.Add(terminal.EntityNameColor("Client:"), terminal.EntityNameColor(c.ClientID))
			table.Add(terminal.EntityNameColor("Access Token:"), c.Access)
		}

		table.Print()
	}
}

func (o *targets) SetCurrentForID(id int) {
	if id <= 0 {
		return
	}

	var target *Target
	var context *Context

	for _, t := range o.list {
		t.Current = false
		if len(t.Contexts) == 0 {
			id--
			if id == 0 {
				target = t
			}
		}
		for _, c := range t.Contexts {
			c.Current = false
			id--
			if id == 0 {
				target = t
				context = c
			}
		}
	}

	if target != nil {
		target.Current = true
	}
	if context != nil {
		context.Current = true
	}
	o.saveConfig()
}

func (o *targets) SetCurrent(url string, cfInstanceURL string, skipSslVerify bool, caCertFile string, tr *lib.TokenResponse) {
	var t *Target
	for _, v := range o.list {
		v.Current = false
		if v.TargetURL == url {
			t = v
		}
	}

	if t == nil {
		t = &Target{
			TargetURL: url,
		}
		o.list = append(o.list, t)
	}
	t.CfInstanceURL = cfInstanceURL
	t.SkipSsl = skipSslVerify
	t.CaCertFilePath = caCertFile
	t.Current = true
	if t.Contexts == nil {
		t.Contexts = map[string]*Context{}
	}

	updateTarget(t, tr)
	o.saveConfig()
}

func (o *targets) LookupAndSetCurrent(url string, context string) bool {
	var t *Target
	for _, v := range o.list {
		v.Current = false
		if v.TargetURL == url {
			t = v
		}
	}

	if t == nil {
		return false
	}
	t.Current = true

	for _, v := range t.Contexts {
		v.Current = false
	}
	c := t.Contexts[context]
	if c == nil {
		return false
	}
	c.Current = true
	t.Contexts[context] = c

	o.saveConfig()
	return true
}

func (o *targets) GetCurrent() (target *Target, context *Context, instance *cf.Item) {
	for _, t := range o.list {
		if t.Current {
			target = t
			for _, c := range t.Contexts {
				if c.Current {
					context = c
					break
				}
			}
			break
		}
	}
	if target == nil {
		global.UI.Failed("No UAA target set. Login to a UAA using the '%s uaa login' command", global.Name)
	} else if context == nil {
		global.UI.Failed("No UAA context set. Login to a UAA using the '%s uaa login' command", global.Name)
	}
	instance = cf.Curl.GetItem(target.CfInstanceURL)
	if instance == nil {
		global.UI.Failed("Failed to fetch target UAA's info")
	}
	return target, context, instance
}

func updateTarget(t *Target, tr *lib.TokenResponse) {
	tc, err := lib.TokenClaimsFactory.New(tr.Access)
	if err == nil {
		for _, v := range t.Contexts {
			v.Current = false
		}

		var currentContext string
		if tc.UserName != "" {
			currentContext = tc.UserName
		} else if tc.ClientID != "" {
			currentContext = tc.ClientID
		}

		if currentContext != "" {
			c := t.Contexts[currentContext]
			if c == nil {
				c = &Context{}
			}

			c.ID = tc.ID
			c.UserID = tc.UserID
			c.UserName = tc.UserName
			c.ClientID = tc.ClientID
			c.Access = tr.Access
			c.Refresh = tr.Refresh
			c.Type = tr.Type
			c.Scopes = tc.Scopes
			c.Current = true

			t.Contexts[currentContext] = c
		}
	}
}

func (o *targets) LoadConfig() {
	o.list = []*Target{}

	if global.Env.ConfigDir != "" {
		configFilePath := filepath.Join(global.Env.ConfigDir, "uaac.json")
		configJSON, err := ioutil.ReadFile(configFilePath)
		if err == nil {
			err := json.Unmarshal(configJSON, &o.list)
			if err != nil {
				_ = os.RemoveAll(configFilePath)
			}
		}
	}
}

func (o *targets) saveConfig() {
	if global.Env.ConfigDir != "" {
		configFilePath := filepath.Join(global.Env.ConfigDir, "uaac.json")
		configJSON, _ := json.MarshalIndent(o.list, "", "  ")
		_ = ioutil.WriteFile(configFilePath, configJSON, os.FileMode(0700))
	}
}
