package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/gowerm123/wadman/pkg/helpers"
	"github.com/gowerm123/wadman/pkg/idGamesClient"
)

var client idGamesClient.Client

func init() {
	log.SetFlags(0)
}

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
	} else if command == "help" {
		handleHelpCommand()
	} else if command == "configure" {
		handleConfigureCommand()
	} else if command == "remove" {
		handleRemoveCommand()
	}

}

func handleInstallCommand() {
	enforceRoot("install")
	args := collectArgs(1, 1)
	query := args[0]
	queryType := getOptional(args, 1, "filename")
	if client.Install(query, queryType) {
		log.Println("Success!")
	}
}

func handleRunCommand() {
	basePath := "/usr/share/wadman"
	args := collectArgs(1, 1)
	firstArg := args[0]
	secondArg := getOptional(args, 1, "")
	var iwad, file string

	if firstArg == "iwad" {
		iwad = secondArg
		if _, err := os.Stat(iwad); err != nil {
			iwad = client.LookupWADAlias(iwad)
		}

		executeCommand(iwad, []string{}, client.Configuration)
		return
	}

	iwad = firstArg
	file = secondArg
	if _, err := os.Stat(iwad); err != nil {
		iwad = client.LookupWADAlias(iwad)
	}

	wadFiles := client.CollectPWads(fmt.Sprintf("%s/%s", basePath, file))
	executeCommand(iwad, wadFiles, client.Configuration)
}

func executeCommand(iwad string, wadFiles []string, config idGamesClient.Configuration) {
	launcher := config.Launcher
	launchArgs := config.LaunchArgs
	log.Println(launchArgs)

	if len(wadFiles) == 0 {
		launchArgs = append([]string{"-iwad ", iwad}, launchArgs...)
	} else {
		launchArgs = append([]string{"-iwad", iwad, "-file", wadFiles[0], wadFiles[1]}, launchArgs...)
	}

	command := exec.Command(launcher, launchArgs...)

	go command.Output()
}

func handleSearchCommand() {
	args := os.Args
	var queryType string = "filename"
	if len(args) < 3 {
		log.Println("invalid search command. proper use is\n    wadman search QUERY (optional)QUERYTYPE\nIf you need help, you can run\n    wadman help\nfor more information")
		os.Exit(0)
	} else if len(args) == 4 {
		queryType = args[3]
	}
	query := args[2]
	log.Println("Searching...")
	client.SearchAndPrint(query, queryType)
}

func handleAliasCommand() {
	enforceRoot("alias")
	args := os.Args
	if len(args) < 4 {
		log.Println("invalid alias command. proper use is\n    wadman alias target alias\nIf you need help, you can run\n    wadman help\nfor more information")
		os.Exit(0)
	}
	target := args[2]
	alias := args[3]

	client.AddAlias(target, alias)
}

func handleRegisterCommand() {
	enforceRoot("register")
	args := collectArgs(2, 0)

	target := args[0]
	iwad := args[1]

	client.RegisterIwad(target, iwad)
}

func handleHelpCommand() {
	printReadme()
}

func handleConfigureCommand() {
	enforceRoot("configure")

	var launcher, launchArgs, iwads, installDir string
	log.Println("Doom Launcher command (default is gzdoom)")
	fmt.Scanln(&launcher)

	log.Println("Extra launch arguments (comma seperated, Example \"fast,respawn,nomonsters\")")
	fmt.Scanln(&launchArgs)

	log.Println("IWADs (See the README), enter as comma seperated key=value pairs. Example doom2=/path/to/DOOM2.WAD,plutonia=/path/to/PLUTONIA.WAD")
	fmt.Scanln(&iwads)

	log.Println("Installation directory for wad archives (default is /usr/share/wadman/)")
	fmt.Scanln(&installDir)

	idGamesClient.UpdateConfigs(launcher, launchArgs, iwads, installDir)
}

func handleRemoveCommand() {
	enforceRoot("remove")
	target := collectArgs(1, 0)[0]

	client.Remove(target)
}

func collectArgs(required, optional int) []string {
	var args []string
	for i := 2; i < required+2; i++ {
		args = append(args, os.Args[i])
	}
	for i := required + 2; i <= required+optional+2 && i < len(os.Args); i++ {
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

func printReadme() {
	content, _ := ioutil.ReadFile("/usr/share/wadman/README.md")

	log.Println(string(content))
}
