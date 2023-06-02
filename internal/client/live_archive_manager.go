package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gowerm123/wadman/internal/helpers"
)

type LiveArchiveManager struct {
	path string

	entries []ArchiveEntry
}

type ArchiveEntry struct {
	Name    string   `json:"name"`
	Dir     string   `json:"dir"`
	Uri     string   `json:"uri"`
	Aliases []string `json:"aliases"`
	Iwad    string   `json:"iwad"`
}

func NewArchiveManager() LiveArchiveManager {
	pm := LiveArchiveManager{}

	pm.path = helpers.GetWadmanHomeDir() + "wadmanifest.json"

	pm.load()

	return pm
}

func (pm *LiveArchiveManager) load() {
	body, err := ioutil.ReadFile(pm.path)
	helpers.HandleFatalErr(err)

	json.Unmarshal(body, &pm.entries)
}

func (pm *LiveArchiveManager) newEntry(filename, path, url string) {
	if pm.Exists(url) {
		return
	}
	pm.entries = append(pm.entries, ArchiveEntry{Name: filename, Dir: path, Uri: url})
	helpers.HandleFatalErr(pm.Commit())
}

func (pm *LiveArchiveManager) Contains(filename, url string) bool {
	for _, entry := range pm.entries {
		if entry.Name == filename && entry.Uri == url {
			return true
		}
	}
	return false
}

func (pm *LiveArchiveManager) Commit() error {
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

func (pm *LiveArchiveManager) GetFilePath(target string) string {
	index := pm.findEntry(target)

	if index == -1 {
		return ""
	}

	return pm.entries[index].Dir
}

func (pm *LiveArchiveManager) AddAlias(target, alias string) {
	index := pm.findEntry(target)

	if helpers.Contains(pm.entries[index].Aliases, alias) {
		log.Println("Alias already known")
		return
	}

	pm.entries[index].Aliases = append(pm.entries[index].Aliases, alias)

	helpers.HandleFatalErr(pm.Commit())
}

func (pm *LiveArchiveManager) Remove(target string) {
	index := pm.findEntry(target)

	dir := pm.entries[index].Dir

	helpers.HandleFatalErr(os.RemoveAll(dir))

	pm.entries = append(pm.entries[:index], pm.entries[index+1:]...)
	helpers.HandleFatalErr(pm.Commit())
}

func (pm *LiveArchiveManager) LookupIwad(target string) string {
	index := pm.findEntry(target)
	if index == -1 {
		index = pm.findByAlias(target)
	}

	return pm.entries[index].Iwad
}

func (pm *LiveArchiveManager) RegisterIwad(target, iwad string) {
	index := pm.findEntry(target)

	pm.entries[index].Iwad = iwad

	helpers.HandleFatalErr(pm.Commit())
}

func (pm *LiveArchiveManager) findEntry(target string) int {
	for i, entry := range pm.entries {
		if entry.Name == target || helpers.Contains(entry.Aliases, target) {
			return i
		}
	}
	return -1
}

func (pm *LiveArchiveManager) findByAlias(target string) int {
	for i, entry := range pm.entries {
		if helpers.Contains(entry.Aliases, target) {
			return i
		}
	}
	return -1
}

// Existence lookups should be performed on idgames url to ensure package distinction
func (pm *LiveArchiveManager) Exists(url string) bool {
	for _, entry := range pm.entries {
		if entry.Uri == url {
			return true
		}
	}
	return false
}

func (pm LiveArchiveManager) Install(file ApiFile, installPath string) bool {
	dirName := strings.Replace(file.Filename, ".zip", "", 1)
	if pm.Contains(dirName, file.IdGamesUrl) {
		log.Println("skipping " + dirName + " it is already installed")
		return false
	}

	unzipped := helpers.GetWadmanHomeDir() + dirName
	err := helpers.Unzip(installPath, unzipped)
	helpers.HandleFatalErr(err, "failed to unzip archive", installPath, "-")

	log.Println("Removing unnnecessary zip archive")
	err = os.Remove(installPath)
	helpers.HandleFatalErr(err, "failed to delete zip archive", installPath, "-")

	pm.newEntry(dirName, unzipped, file.IdGamesUrl)

	return true
}
