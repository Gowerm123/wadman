package idGamesClient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gowerm123/wadman/pkg/helpers"
)

var path string = helpers.WadmanConfigPath()

const (
	launchKey = "LAUNCHER"
	lArgKey   = "LAUNCHARGS"
	iwadKey   = "IWAD"
)

var setKeys = []string{launchKey, iwadKey, lArgKey}

type LaunchConfiguration struct {
	Launcher   string            `json:"launcher"`
	LaunchArgs []string          `json:"launchArgs"`
	IWads      map[string]string `json:"iwads"`
}

func loadConfigs() LaunchConfiguration {
	bytes, err := ioutil.ReadFile(path)
	helpers.HandleFatalErr(err)

	var config LaunchConfiguration
	err = json.Unmarshal(bytes, &config)
	helpers.HandleFatalErr(err)

	return config
}

func (cfg LaunchConfiguration) Update(key, value string) {
	if key == iwadKey {
		spl := strings.Split(value, "=")
		if len(spl) != 2 {
			log.Fatal("error in assigning iwad, please use the format, wadman -s IWAD ${ALIAS}=${PATH_TO_IWAD}")
		}

		cfg.IWads[spl[0]] = spl[1]
	} else {
		switch key {
		case lArgKey:
			cfg.LaunchArgs = strings.Split(value, ",")
		case launchKey:
			cfg.Launcher = value
		}
	}

	cfg.Commit()
}

func (cfg LaunchConfiguration) Commit() {
	bytes, _ := json.MarshalIndent(cfg, "", "	")

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
