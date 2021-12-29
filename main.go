package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/gowerm123/wadman/pkg/helpers"
	"github.com/gowerm123/wadman/pkg/idGamesClient"
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
	} else if command == "register" {
		handleRegisterCommand()
	}

}

func handleInstallCommand() {
	enforceRoot("install")
	args := collectArgs(2, 1)
	query := args[1]
	queryType := getOptional(args, 2, "filename")
	client.Install(query, queryType)
}

func handleRunCommand() {
	args := collectArgs(2, 1)
	firstArg := args[1]
	secondArg := getOptional(args, 2, "")

	var iwad, file string
	if secondArg == "" {
		file = firstArg
		iwad = client.LookupIwad(file)
	} else {
		file = secondArg
		iwad = firstArg
	}

	launcher := client.Configuration.Launcher
	basePath := client.Configuration.InstallDir

	command := exec.Command(launcher, "-IWAD", iwad, "-file", basePath+file, "&")
	command.Output()
}

func handleSearchCommand() {
	args := os.Args
	var queryType string = "filename"
	if len(args) < 3 {
		fmt.Println("invalid search command. proper use is\n    wadman search QUERY (optional)QUERYTYPE\nIf you need help, you can run\n    wadman help\nfor more information")
		os.Exit(0)
	} else if len(args) == 4 {
		queryType = args[3]
	}
	query := args[2]
	fmt.Println("Searching...")
	client.SearchAndPrint(query, queryType)
}

func handleAliasCommand() {
	enforceRoot("alias")
	args := os.Args
	if len(args) < 4 {
		fmt.Println("invalid alias command. proper use is\n    wadman alias target alias\nIf you need help, you can run\n    wadman help\nfor more information")
		os.Exit(0)
	}
	target := args[2]
	alias := args[3]

	client.AddAlias(target, alias)
}

func handleRegisterCommand() {
	args := collectArgs(2, 0)

	target := args[0]
	iwad := args[1]

	client.RegisterIwad(target, iwad)
}

func collectArgs(required, optional int) []string {
	var args []string
	for i := 1; i <= required+1; i++ {
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

func getOptional(args []string, index int, defaultVal string) string {
	if len(args) <= index {
		return defaultVal
	}
	return args[index]
}

func enforceRoot(cmd string) {
	if !isRoot() {
		fmt.Printf("please execute wadman as root when using the %s command\n", cmd)
		os.Exit(1)
	}
}

func isRoot() bool {
	currentUser, err := user.Current()
	helpers.HandleFatalErr(err)

	return currentUser.Username == "root"
}
