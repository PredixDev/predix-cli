package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.build.ge.com/adoption/predix-cli/cli/global"
)

type Config struct {
	Timeouts Timeouts `json:"timeouts"`
}

type Timeouts map[string]int

func ReadConfigFile() (config Config) {
	cacheDir := global.Env.ConfigDir

	if cacheDir != "" {
		configFilePath := filepath.Join(cacheDir, "cache_config")
		configJSON, err := ioutil.ReadFile(configFilePath)

		if err == nil {
			err := json.Unmarshal(configJSON, &config)

			if err != nil {
				_ = os.RemoveAll(configFilePath)
				config = Config{}
			}
		}
	}

	return config
}

func SaveTimeouts(timeouts Timeouts) {
	config := ReadConfigFile()
	cacheDir := global.Env.ConfigDir

	if cacheDir != "" {
		configFilePath := filepath.Join(cacheDir, "cache_config")

		if config.Timeouts == nil {
			config.Timeouts = Timeouts{}
		}

		for k, v := range timeouts {
			config.Timeouts[k] = v
		}

		configJSON, _ := json.Marshal(config)
		_ = ioutil.WriteFile(configFilePath, configJSON, os.FileMode(0700))
	}
}
