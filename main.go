package main

import (
	"fmt"
	"log"
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
	} else if command == "help" {
		log.Println("WadMan supports eight basic commands\n\n  - `search QUERY <QUERYTYPE>` - Searches the IdGames archive for the specified QUERY,QUERYTYPE is optional, and defaults to filename. Possible options are filename, title, author, email, description, credits, editors, textfile.\n  - `install QUERY <QUERYTYPE>` - First performs a `search QUERY <QUERYTYPE>` then installs the first found file. It is recommended that you search based on filename here to narrow down overlapping projects.\b  - `list` - Lists all currently installed wad archives. Information printed is name of archive, installed directory, idGamesUrl, and Aliases.\n  - `remove NAME` - Removes the archive with the given name. If two are found, the first will be deleted.\n  - `run` - There are two ways to call `run`. You can either call `run ALIAS/NAME` or `run IWAD ALIAS/NAME`. Note that you must include the IWAD if you have not registered an IWAD to the given `ALIAS/NAME`.\n  - `register NAME IWAD` - Assigns the IWAD to the archive entry in the `pkglist` associated with NAME. This is used for the `run` command so you do not have to specify IWADs everytime you load a PWAD.\n  - `configure` - Runs you through a prompt to fill out the configuration file. The file is a simple JSON file found at `/usr/share/.wadmanConfig`\n  - `help` - Prints this text")
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
	args := collectArgs(1, 1)
	firstArg := args[0]
	secondArg := getOptional(args, 1, "")
	var iwad, file string
	if secondArg == "" {
		file = firstArg
		iwad = client.LookupIwad(file)
	} else {
		file = secondArg
		iwad = firstArg
	}

	//Check if iwad path exists, if not, assume alias
	if _, err := os.Stat(iwad); err != nil {
		iwad = client.LookupWADAlias(iwad)
	}

	launcher := client.Configuration.Launcher
	basePath := client.Configuration.InstallDir
	fmt.Println(iwad)
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
	enforceRoot("register")
	args := collectArgs(2, 0)

	target := args[0]
	iwad := args[1]

	client.RegisterIwad(target, iwad)
}

func handleConfigureCommand() {
	enforceRoot("configure")

	var launcher, launchArgs, iwads, installDir string
	fmt.Println("Doom Launcher command (default is gzdoom)")
	fmt.Scanln(&launcher)

	fmt.Println("Extra launch arguments (comma seperated, Example \"fast,respawn,nomonsters\")")
	fmt.Scanln(&launchArgs)

	fmt.Println("IWADs (See the README), enter as comma seperated key=value pairs. Example doom2=/path/to/DOOM2.WAD,plutonia=/path/to/PLUTONIA.WAD")
	fmt.Scanln(&iwads)

	fmt.Println("Installation directory for wad archives (default is /usr/share/wadman/)")
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
