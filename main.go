package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gowerm/dwpm/pkg/dwclient"
)

func main() {
	args := os.Args

	client := dwclient.New()

	command := strings.ToLower(args[1])

	client.ValidateCommand(command)

	if command == "install" {
		var queryType string = "filename"
		if len(args) < 3 {
			fmt.Println("invalid install command. proper use is\n    dwpm install QUERY (optional)QUERYTYPE\nIf you need help, you can run\n    dwpm help\nfor more information")
			os.Exit(0)
		} else if len(args) == 4 {
			queryType = args[3]
		}
		query := args[2]

		client.Install(query, queryType)
	} else if command == "run" {
		if len(args) < 4 {
			fmt.Println("invalid run command. proper use is\n    dwpm run IWAD TARGET\nIf you need help, you can run\n    dwpm help\nfor more information")
			os.Exit(0)
		}

		iwad := args[2]
		target := args[3]

		if iwad == "doom2" {
			iwad = "/home/matt/Doom/DOOM2.WAD"
		}

		launcher := "gzdoom"
		basePath := "/home/matt/dwpm/"

		command := exec.Command(launcher, "-IWAD", iwad, "-file", basePath+target, "&")
		command.Output()
	} else if command == "search" {
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
	} else if command == "list" {
		client.List()
	} else if command == "alias" {
		if len(args) < 4 {
			fmt.Println("invalid alias command. proper use is\n    dwpm alias target alias\nIf you need help, you can run\n    dwpm help\nfor more information")
			os.Exit(0)
		}
		target := args[2]
		alias := args[3]

		client.AddAlias(target, alias)
	}

}