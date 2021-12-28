package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gowerm/dwpm/pkg/idGamesClient"
)

var client idGamesClient.Client

func main() {
	args := os.Args

	client = idGamesClient.New()

	command := strings.ToLower(args[1])

	client.ValidateCommand(command)

	if command == "install" {
		handleInstallCommand()
	} else if command == "run" {
		handleRunCommand()
	} else if command == "search" {
		handleSearchCommand()
	} else if command == "list" {
		client.List()
	} else if command == "alias" {
		handleAliasCommand()
	}

}

func handleInstallCommand() {

	client.Install(query, queryType)
}

func handleRunCommand() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("invalid run command. proper use is\n    dwpm run IWAD TARGET\nIf you need help, you can run\n    dwpm help\nfor more information")
		os.Exit(0)
	}

	iwad := args[2]
	target := args[3]

	launcher := client.Configuration.Launcher
	basePath := client.Configuration.InstallDir

	command := exec.Command(launcher, "-IWAD", iwad, "-file", basePath+target, "&")
	command.Output()
}

func handleSearchCommand() {
	args := os.Args
	var queryType string = "filename"
	if len(args) < 3 {
		fmt.Println("invalid search command. proper use is\n    dwpm search QUERY (optional)QUERYTYPE\nIf you need help, you can run\n    dwpm help\nfor more information")
		os.Exit(0)
	} else if len(args) == 4 {
		queryType = args[3]
	}
	query := args[2]
	fmt.Println("Searching...")
	client.SearchAndPrint(query, queryType)
}

func handleAliasCommand() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("invalid alias command. proper use is\n    dwpm alias target alias\nIf you need help, you can run\n    dwpm help\nfor more information")
		os.Exit(0)
	}
	target := args[2]
	alias := args[3]

	client.AddAlias(target, alias)
}

func collectArgs(required, optional int) []string {
	var args []string
	for i := 1; i < required; i++ {
		args = append(args, os.Args[i])
	}
	for i := required + 1; i < required+optional; i++ {
		if i >= len(os.Args) {
			break
		}
		args = append(args, os.Args[i])
	}

	return args
}
