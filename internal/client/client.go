package client

type IdGamesClient interface {
	SearchAndPrint(string) string
	Install(string, ArchiveManager) bool
	List()
	Set(string, string)
	ValidateCommand(string)
}
