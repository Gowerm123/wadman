package idGamesClient

import (
	"encoding/json"
	"fmt"
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

func newPackageManager(path string) packageManager {
	pm := packageManager{}

	pm.path = path + ".pkglist"

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
	helpers.HandleFatalErr(pm.Commit())
}

func (pm *packageManager) Contains(filename, url string) bool {
	for _, entry := range pm.entries {
		if entry.Name == filename && entry.Uri == url {
			return true
		}
	}
	return false
}

func (pm *packageManager) Commit() error {
	bytes, err := json.Marshal(pm.entries)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pm.path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
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

func (pm *packageManager) AddAlias(target, alias string) {
	for i, entry := range pm.entries {
		if entry.Name == target {
			if helpers.Contains(entry.Aliases, alias) {
				fmt.Println("Skipping - known alias")
			}
			entry.Aliases = append(entry.Aliases, alias)
			pm.entries[i] = entry
		}
	}

	helpers.HandleFatalErr(pm.Commit())
}

func (pm *packageManager) Remove(target string) {

}
