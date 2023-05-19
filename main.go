package main

import (
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
	case "-r", "--remove":
		handleRemoveCommand(arguments)
		return
	case "-q", "--query":
		handleSearchCommand(arguments)
		return
	case "-p", "--play":
		handleRunCommand(arguments)
		return
	case "-l", "--list":
		client.List()
		return
	case "-s", "--set":
		handleSetCommand(arguments)
	default:
		log.Println("Unknown command")
		return
	}
}

func handleInstallCommand(arguments ArrayFlags) {
	enforceRoot("install")
	for _, argument := range arguments {
		if !client.Install(argument) {
			log.Fatalf("failed to install target %s", argument)
		}
	}
}

func handleSetCommand(args ArrayFlags) {
	enforceRoot("set")
	if len(args) != 2 {
		log.Fatal("format for set command is wadman -s KEY VALUE\nPlease see help section for list of available KEYs")
	}

	client.Set(args[0], args[1])
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

	launcher := client.LaunchConfiguration.Launcher

	wadFiles := client.CollectPWads(helpers.GetWadmanHomeDir() + file)

	var command *exec.Cmd
	if wadFiles[1] == "" {
		command = exec.Command(launcher, "-iwad", iwad, "-file", wadFiles[0])
	} else {
		command = exec.Command(launcher, "-iwad", iwad, "-file", wadFiles[0], wadFiles[1])
	}

	if _, err := command.Output(); err != nil {
		helpers.HandleFatalErr(err)
	}
}

func handleSearchCommand(args ArrayFlags) {
	buffer := ""
	for _, arg := range args {
		buffer += client.SearchAndPrint(arg)
	}
	log.Println(buffer)
}

func handleRegisterCommand(args ArrayFlags) {
	enforceRoot("register")

	target := args[0]
	iwad := args[1]

	client.RegisterIwad(target, iwad)
}

func handleRemoveCommand(args ArrayFlags) {
	enforceRoot("remove")

	for _, target := range args {
		client.Remove(target)
	}
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
