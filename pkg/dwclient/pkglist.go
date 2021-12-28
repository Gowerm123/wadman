package dwclient

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gowerm/dwpm/pkg/helpers"
)

type packageManager struct {
	path string

	entries []packageEntry
}

type packageEntry struct {
	Name    string   `json:"name"`
	Dir     string   `json:"dir"`
	Uri     string   `json:"uri"`
	Aliases []string `json:"aliases"`
}

func newPackageManager() packageManager {
	pm := packageManager{}

	pm.path = LOCALPATH + ".pkglist"

	pm.load()

	return pm
}

func (pm *packageManager) load() {
	body, err := ioutil.ReadFile(pm.path)
	helpers.HandleFatalErr(err)

	json.Unmarshal(body, &pm.entries)
}

func (pm *packageManager) NewEntry(filename, path, url string) {
	pm.entries = append(pm.entries, packageEntry{Name: filename, Dir: path, Uri: url})
	pm.Commit()
}

func (pm *packageManager) Contains(filename, url string) bool {
	for _, entry := range pm.entries {
		if entry.Name == filename && entry.Uri == url {
			return true
		}
	}
	return false
}

func (pm *packageManager) Commit() {
	bytes, err := json.Marshal(pm.entries)
	helpers.HandleFatalErr(err)

	err = ioutil.WriteFile(pm.path, bytes, 0644)
	helpers.HandleFatalErr(err)
}

func (pm *packageManager) GetFilePath(filename string) string {
	body, err := ioutil.ReadFile(pm.path)
	helpers.HandleFatalErr(err, "failed to read pkglist")

	entries := strings.Split(string(body), "\n")

	for _, entry := range entries {
		spl := strings.Split(entry, " ")
		if spl[0] == filename {
			return spl[1]
		}
	}
	return ""
}
