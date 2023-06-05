package client

type IdGamesClient interface {
	Search(string) string
	Install(string, ArchiveManager) bool
	Set(string, string)
	ValidateCommand(string)
}
