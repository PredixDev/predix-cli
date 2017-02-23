package global

import (
	"os"
	"strconv"
)

type Environment struct {
	CfHomeDir string
	ConfigDir string
	NoCache   bool
}

var Env = Environment{}

func init() {
	Env.NoCache = false
	noCache := os.Getenv("PREDIX_NO_CACHE")
	b, err := strconv.ParseBool(noCache)
	if err == nil {
		Env.NoCache = b
	}
}
