package idGamesClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gowerm123/wadman/pkg/helpers"
)

var path string = fmt.Sprintf("%s/.config/wadman-config.json", helpers.GetHome())

type Configuration struct {
	Launcher   string            `json:"launcher"`
	LaunchArgs []string          `json:"launchArgs"`
	IWads      map[string]string `json:"iwads"`
	InstallDir string            `json:"installDir"`
}

func loadConfigs() Configuration {
	bytes, err := ioutil.ReadFile(path)
	helpers.HandleFatalErr(err)

	var config Configuration
	fmt.Println(path)
	err = json.Unmarshal(bytes, &config)
	helpers.HandleFatalErr(err)

	return config
}

func UpdateConfigs(launcher, args, iwads, installPath string) {
	var config Configuration
	if launcher != "" {
		config.Launcher = launcher
	} else {
		config.Launcher = "gzdoom"
	}
	if args != "" {
		config.LaunchArgs = strings.Split(args, ",")
	}
	if iwads != "" {
		config.IWads = convertToMap(iwads)
	}
	if installPath != "" {
		config.InstallDir = installPath
	} else {
		config.InstallDir = fmt.Sprintf("%s/.wadman/", helpers.GetHome())
	}

	CommitConfig(config)
}

func CommitConfig(config Configuration) {
	bytes, _ := json.MarshalIndent(config, "", "	")

	helpers.HandleFatalErr(os.WriteFile(path, bytes, 0644))
}

func convertToMap(str string) map[string]string {
	spl := strings.Split(str, ",")

	var mp map[string]string = make(map[string]string)
	for _, pair := range spl {
		spl2 := strings.Split(pair, "=")

		mp[spl2[0]] = spl2[1]
	}

	return mp
}
