package client

type ArchiveManager interface {
	Install(ApiFile) bool
	AddAlias(string, string)
	LookupLocalPath(string) string
	CollectPWads(dir string) []string
	LookupWADAlias(string) string
	LookupIwad(string) string
	RegisterIwad(string, string)
	Remove(string)
}
