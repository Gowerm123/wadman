package idGamesClient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gowerm123/wadman/pkg/helpers"
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
	Iwad    string   `json:"iwad"`
}

func newPackageManager() packageManager {
	pm := packageManager{}

	pm.path = helpers.GetWadmanHomeDir() + "wadmanifest.json"

	pm.load()

	return pm
}

func (pm *packageManager) load() {
	body, err := ioutil.ReadFile(pm.path)
	helpers.HandleFatalErr(err)

	json.Unmarshal(body, &pm.entries)
}

func (pm *packageManager) NewEntry(filename, path, url string) {
	if pm.Exists(url) {
		return
	}
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
	bytes, err := json.MarshalIndent(pm.entries, "", "	")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pm.path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (pm *packageManager) GetFilePath(target string) string {
	index := pm.findEntry(target)

	if index == -1 {
		return ""
	}

	return pm.entries[index].Dir
}

func (pm *packageManager) AddAlias(target, alias string) {
	index := pm.findEntry(target)

	if helpers.Contains(pm.entries[index].Aliases, alias) {
		log.Println("Alias already known")
		return
	}

	pm.entries[index].Aliases = append(pm.entries[index].Aliases, alias)

	helpers.HandleFatalErr(pm.Commit())
}

func (pm *packageManager) Remove(target string) {
	index := pm.findEntry(target)

	dir := pm.entries[index].Dir

	helpers.HandleFatalErr(os.RemoveAll(dir))

	pm.entries = append(pm.entries[:index], pm.entries[index+1:]...)
	helpers.HandleFatalErr(pm.Commit())
}

func (pm *packageManager) LookupIwad(target string) string {
	index := pm.findEntry(target)
	if index == -1 {
		index = pm.findByAlias(target)
	}

	return pm.entries[index].Iwad
}

func (pm *packageManager) RegisterIwad(target, iwad string) {
	index := pm.findEntry(target)

	pm.entries[index].Iwad = iwad

	helpers.HandleFatalErr(pm.Commit())
}

func (pm *packageManager) findEntry(target string) int {
	for i, entry := range pm.entries {
		if entry.Name == target || helpers.Contains(entry.Aliases, target) {
			return i
		}
	}
	return -1
}

func (pm *packageManager) findByAlias(target string) int {
	for i, entry := range pm.entries {
		if helpers.Contains(entry.Aliases, target) {
			return i
		}
	}
	return -1
}

// Existence lookups should be performed on idgames url to ensure package distinction
func (pm *packageManager) Exists(url string) bool {
	for _, entry := range pm.entries {
		if entry.Uri == url {
			return true
		}
	}
	return false
}
