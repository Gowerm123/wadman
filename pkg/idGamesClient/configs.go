package idGamesClient

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gowerm/dwpm/pkg/helpers"
)

var path string = "/usr/share/.dwpmConfig"

type Configuration struct {
	Launcher   string            `json:"launcher"`
	LaunchArgs []string          `json:"launchArgs"`
	IWads      map[string]string `json:"iwads"`
	InstallDir string            `json:"installDir`
}

func SetConfigPath(customPath string) {
	path = customPath
}

func loadConfigs() Configuration {
	bytes, err := ioutil.ReadFile(path)
	helpers.HandleFatalErr(err)

	var config Configuration
	err = json.Unmarshal(bytes, &config)
	helpers.HandleFatalErr(err)

	return config
}
