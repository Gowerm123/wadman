package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gowerm123/wadman/pkg/helpers"
	"github.com/gowerm123/wadman/pkg/idGamesClient"
)

var client idGamesClient.Client

func init() {
	log.SetFlags(0)
}

type ArrayFlags []string

func (af *ArrayFlags) String() string {
	str := ""

	for _, flag := range *af {
		str += " "
		str += flag
	}

	str = str[:0]

	return str
}

func (af *ArrayFlags) Set(str string) error {
	arr := []string{str}
	if len(strings.Split(str, " ")) > 0 {
		arr = strings.Split(str, " ")
	}
	for _, word := range arr {
		*af = append(*af, word)
	}

	return nil
}

func main() {
	client = idGamesClient.New()
	command, arguments := parseCli()
	switch command {
	case "-i", "--install":
		handleInstallCommand(arguments)
		return
	case "-r", "--run":
		handleRunCommand(arguments)
		return
	case "-s", "--search":
		handleSearchCommand(arguments)
		return
	case "-u", "--uninstall":
		handleRemoveCommand(arguments)
		return
	case "-a", "--assign":
		handleRegisterCommand(arguments)
		return
	case "-c", "--configure":
		handleConfigureCommand()
		return
	case "-l", "--list":
		client.List()
		return
	default:
		log.Println("Unknown command")
		return
	}
}

func handleInstallCommand(arguments ArrayFlags) {
	enforceRoot("install")
	for _, argument := range arguments {
		if client.Install(argument) {
			log.Println("Success!")
		}
	}
}

func handleRunCommand(args ArrayFlags) {
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

	wadFiles := client.CollectPWads(basePath + file)

	var command *exec.Cmd
	if wadFiles[1] == "" {
		command = exec.Command(launcher, "-iwad", iwad, "-file", wadFiles[0])
	} else {
		command = exec.Command(launcher, "-iwad", iwad, "-file", wadFiles[0], wadFiles[1])
	}

	command.Output()
}

func handleSearchCommand(args ArrayFlags) {
	buffer := ""
	for _, arg := range args {
		buffer += client.SearchAndPrint(arg)
	}
	log.Println(buffer)
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

func handleRegisterCommand(args ArrayFlags) {
	enforceRoot("register")

	target := args[0]
	iwad := args[1]

	client.RegisterIwad(target, iwad)
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

	log.Println("Installation directory for wad archives (default is $HOME/.wadman)")
	fmt.Scanln(&installDir)

	idGamesClient.UpdateConfigs(launcher, launchArgs, iwads, installDir)
}

func handleRemoveCommand(args ArrayFlags) {
	enforceRoot("remove")

	for _, target := range args {
		client.Remove(target)
	}
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
	if !helpers.IsRoot() {
		log.Printf("please execute wadman as root when using the %s command\n", cmd)
		os.Exit(1)
	}
}

func parseCli() (cmd string, arguments []string) {
	cmd = os.Args[1]
	arguments = os.Args[2:]

	return cmd, arguments
}
