package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gowerm123/wadman/internal/client"
	"github.com/gowerm123/wadman/internal/helpers"
)

var idGamesClient client.IdGamesClient
var archiveManager client.ArchiveManager

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
	idGamesClient = client.New()
	archiveManager = client.NewArchiveManager()
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
		archiveManager.List()
		return
	case "-s", "--set":
		handleSetCommand(arguments)
	case "-sy", "--sync":
		handleSyncCommand()
	default:
		log.Println("Unknown command")
		return
	}
}

func handleInstallCommand(arguments ArrayFlags) {
	enforceRoot("install")
	for _, argument := range arguments {
		if !archiveManager.CheckExists(argument) {
			if !idGamesClient.Install(argument, archiveManager) {
				log.Println("no archives found using filename, expanding search")
				if !idGamesClient.Install(argument, archiveManager) {
					log.Fatalf("failed to install target %s", argument)
				}
			}
			log.Println("successfully installed archive")
		} else {
			log.Printf("target %s is already installed", argument)
			os.Exit(0)
		}
	}
}

func handleSyncCommand() {
	enforceRoot("sync")
	for _, entry := range archiveManager.Entries() {
		if st, err := os.Stat(helpers.GetWadmanHomeDir() + helpers.ToTargetName(entry.Name)); err == nil && st.IsDir() {
			log.Printf("target %s already installed\n", helpers.ToTargetName(entry.Name))
		} else {
			log.Printf("syncing %s\n", entry.Name)
			ToLiveClient(idGamesClient).InstallByFile(entry.ToFile(), archiveManager)
		}
	}
}

func handleSetCommand(args ArrayFlags) {
	enforceRoot("set")
	if len(args) != 2 {
		log.Fatal("format for set command is wadman -s KEY VALUE\nPlease see help section for list of available KEYs")
	}

	idGamesClient.Set(args[0], args[1])
}

func handleRunCommand(args ArrayFlags) {
	firstArg := args[0]
	secondArg := getOptional(args, 1, "")

	iwad, file := organizeIntoIwadAndFileArguments(firstArg, secondArg)

	//Check if iwad path exists, if not, assume alias
	if _, err := os.Stat(iwad); err != nil {
		iwad = LookupWADAlias(iwad, idGamesClient)
	}

	launcher := ToLiveClient(idGamesClient).Configuration.Launcher
	var wadFiles []string
	if file != "" {
		wadFiles = archiveManager.CollectPWads(helpers.GetWadmanHomeDir() + file)
	}
	var command *exec.Cmd
	if len(wadFiles) == 1 || (len(wadFiles) > 0 && wadFiles[1] == "") {
		command = exec.Command(launcher, "-iwad", iwad, "-file", wadFiles[0])
	} else if len(wadFiles) > 1 || (len(wadFiles) > 0 && wadFiles[1] != "") {
		command = exec.Command(launcher, "-iwad", iwad, "-file", wadFiles[0], wadFiles[1])
	} else {
		command = exec.Command(launcher, "-iwad", iwad)
	}

	if outp, err := command.Output(); err != nil {
		helpers.HandleFatalErr(err)
	} else {
		log.Println(string(outp))
	}

}

func organizeIntoIwadAndFileArguments(arg1, arg2 string) (string, string) {
	var iwad, file string
	if arg2 == "" {
		file = arg1
		iwad = archiveManager.LookupIwad(file)
		if iwad == "" {
			iwad = file
			file = ""
		}
	} else {
		file = arg2
		iwad = arg1
	}

	return iwad, file
}

func handleSearchCommand(args ArrayFlags) {
	buffer := ""
	for _, arg := range args {
		buffer += idGamesClient.Search(arg)
	}
	log.Println(buffer)
}

func handleRegisterCommand(args ArrayFlags) {
	enforceRoot("register")

	target := args[0]
	iwad := args[1]

	archiveManager.RegisterIwad(target, iwad)
}

func handleRemoveCommand(args ArrayFlags) {
	enforceRoot("remove")

	for _, target := range args {
		if archiveManager.Remove(target) {
			log.Printf("successfully removed %s\n", target)
		}
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

func LookupWADAlias(alias string, cli client.IdGamesClient) string {
	return ToLiveClient(cli).Configuration.IWads[alias]
}

func ToLiveClient(cli client.IdGamesClient) client.LiveClient {
	return interface{}(cli).(client.LiveClient)
}
