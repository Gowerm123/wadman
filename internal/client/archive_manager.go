package client

import "github.com/gowerm123/wadman/internal/helpers"

type ArchiveManager interface {
	Install(ApiFile) bool
	AddAlias(string, string)
	GetFilePath(string) string
	CollectPWads(dir string) []string
	LookupIwad(string) string
	RegisterIwad(string, string)
	CheckExists(string) bool
	Remove(string) bool
	Entries() []ArchiveEntry
	List()
	InstallSyncPackage(ApiFile) bool
}

type ArchiveEntry struct {
	Name    string   `json:"name"`
	Dir     string   `json:"dir"`
	Uri     string   `json:"uri"`
	Aliases []string `json:"aliases"`
	Iwad    string   `json:"iwad"`
}

func (ae ArchiveEntry) ToFile() *ApiFile {
	return &ApiFile{
		Filename:   helpers.ToZipFileName(ae.Name),
		IdGamesUrl: ae.Uri,
	}
}
