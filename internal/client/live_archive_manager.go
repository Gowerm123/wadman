package client

import (
	"encoding/json"
	"fmt"
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

func NewArchiveManager() ArchiveManager {
	pm := LiveArchiveManager{}

	pm.path = helpers.GetWadmanHomeDir() + "wadmanifest.json"

	pm.load()

	return &pm
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
	helpers.HandleFatalErr(pm.Commit(), "failed to commit to wadmanifest")
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

func (pm *LiveArchiveManager) Remove(target string) bool {
	index := pm.findEntry(target)
	if index == -1 {
		log.Printf("target %s does not exist\n", target)
		return false
	}
	dir := pm.entries[index].Dir

	helpers.HandleFatalErr(os.RemoveAll(dir))

	pm.entries = append(pm.entries[:index], pm.entries[index+1:]...)
	helpers.HandleFatalErr(pm.Commit())

	return true
}

func (pm *LiveArchiveManager) LookupIwad(target string) string {
	index := pm.findEntry(target)
	if index == -1 {
		index = pm.findByAlias(target)
	}

	if index == -1 {
		return ""
	} else {
		return pm.entries[index].Iwad
	}

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

func (pm *LiveArchiveManager) CheckExists(identifier string) bool {
	for _, entry := range pm.entries {
		if entry.Name == identifier || helpers.Contains(entry.Aliases, identifier) || entry.Dir == identifier || entry.Uri == identifier {
			return true
		}
	}
	return false
}

func (pm *LiveArchiveManager) Install(file ApiFile) bool {
	dirName := strings.Replace(file.Filename, ".zip", "", 1)
	if pm.Contains(dirName, file.IdGamesUrl) {
		log.Println("skipping " + dirName + " it is already installed")
		return false
	}

	unzipped := helpers.GetWadmanHomeDir() + dirName
	err := helpers.Unzip(helpers.GetWadmanHomeDir()+file.Filename, unzipped)
	helpers.HandleFatalErr(err, "failed to unzip archive", file.Filename, "-")

	log.Println("Removing unnnecessary zip archive")
	err = os.Remove(helpers.GetWadmanHomeDir() + file.Filename)
	helpers.HandleFatalErr(err, "failed to delete zip archive", file.Filename, "-")

	pm.newEntry(dirName, unzipped, file.IdGamesUrl)

	return true
}

func (pm *LiveArchiveManager) InstallSyncPackage(file ApiFile) bool {
	dirName := strings.Replace(file.Filename, ".zip", "", 1)

	unzipped := helpers.GetWadmanHomeDir() + dirName
	err := helpers.Unzip(helpers.GetWadmanHomeDir()+file.Filename, unzipped)
	helpers.HandleFatalErr(err, "failed to unzip archive", file.Filename, "-")

	log.Println("Removing unnnecessary zip archive")
	err = os.Remove(helpers.GetWadmanHomeDir() + file.Filename)
	helpers.HandleFatalErr(err, "failed to delete zip archive", file.Filename, "-")

	pm.newEntry(dirName, unzipped, file.IdGamesUrl)

	return true
}

func (am *LiveArchiveManager) CollectPWads(dir string) []string {
	pwads := searchForWads(dir)

	return pwads
}

func (am *LiveArchiveManager) Entries() []ArchiveEntry {
	return am.entries
}

func (am *LiveArchiveManager) List() {
	for _, entry := range am.entries {
		log.Printf("Package - Name: %s, Dir: %s, Uri: %s, Aliases: %s\n", entry.Name, entry.Dir, entry.Uri, entry.Aliases)
	}
}

func searchForWads(dir string) []string {
	var buffer []string
	var wads []string

	push(&buffer, dir)

	for len(buffer) > 0 {
		path := pop(&buffer)

		entries, err := os.ReadDir(path)
		helpers.HandleFatalErr(err)

		for _, entry := range entries {
			if entry.IsDir() {
				push(&buffer, fmt.Sprintf("%s/%s", path, entry.Name()))
			} else if isPWad(entry.Name()) {
				push(&wads, fmt.Sprintf("%s/%s", path, entry.Name()))
			}
		}
	}

	return wads
}
