package idGamesClient

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gowerm123/wadman/pkg/helpers"
)

const idGamesBaseURI = "http://www.doomworld.com/idgames/api/api.php"
const idGamesSubstr = "idgames://"

var (
	mirrors       []string = []string{"mirrors.syringanetworks.net", "www.quaddicted.com", "ftpmirror1.infania.net"}
	validCommands []string = []string{"search", "install", "run", "help", "list", "alias", "register", "configure", "remove"}
)

type Client struct {
	httpClient http.Client
	packageManager
	Configuration
}

type Payload struct {
	Body    string
	Headers http.Header
}

type searchResponse struct {
	Files []apiFile `xml:"content>file"`
}

type apiFile struct {
	Author      string  `xml:"author"`
	Email       string  `xml:"email"`
	Title       string  `xml:"title"`
	Dir         string  `xml:"dir"`
	Filename    string  `xml:"filename"`
	Size        int64   `xml:"size"`
	Age         int64   `xml:"age"`
	Date        string  `xml:"date"`
	Description string  `xml:"description"`
	Rating      float32 `xml:"rating"`
	Votes       int64   `xml:"votes"`
	Url         string  `xml:"url"`
	IdGamesUrl  string  `xml:"idgamesurl"`
}

func New() Client {
	var client Client

	client.Configuration = loadConfigs()
	client.httpClient = *http.DefaultClient
	client.packageManager = newPackageManager(client.Configuration.InstallDir)

	return client
}

func (dwc *Client) sendQuery(query, queryType string) searchResponse {
	pl, err := dwc.dial("search", map[string]string{"query": query, "type": queryType})
	helpers.HandleFatalErr(err, "failed to send search query")

	response := searchResponse{}

	xml.Unmarshal([]byte(pl.Body), &response)

	return response
}

func (dwc *Client) search(query, queryType string) []apiFile {
	response := dwc.sendQuery(query, queryType)
	return response.Files
}

func (dwc *Client) SearchAndPrint(query, queryType string) {
	response := dwc.sendQuery(query, queryType)

	for _, entry := range response.Files {
		formatAndPrint(entry)
	}
}

func formatAndPrint(file apiFile) {
	file = sanitizeFile(file)
	log.Printf("File Found\n    Filename: %s\n    Title: %s,\n    Author: %s,\n    Date: %s\n    Url: %s\n", file.Filename, file.Title, file.Author, file.Date, file.IdGamesUrl)
}

func (dwc *Client) dial(action string, params map[string]string) (Payload, error) {
	actionEndpoint := fmt.Sprintf("%s?action=%s", idGamesBaseURI, action)

	for k, v := range params {
		actionEndpoint = fmt.Sprintf("%s&%s=%s", actionEndpoint, url.QueryEscape(k), url.QueryEscape(v))
	}

	response, err := dwc.httpClient.Get(actionEndpoint)
	if err != nil {
		return Payload{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Payload{}, err
	}

	return Payload{Body: string(body), Headers: response.Header}, nil
}

func (dwc *Client) Install(query, queryType string) bool {
	files := dwc.search(query, queryType)

	var choice int = 0

	if len(files) == 0 {
		log.Print("Entry not found for search query, try a different QUERYTYPE?")
		os.Exit(0)
	} else if len(files) > 1 {
		log.Println("Multiple files found, please choose...")
		for it, file := range files {
			file = sanitizeFile(file)
			fmt.Printf("%d) %s, by %s, file - %s\n", it, file.Title, file.Author, file.Filename)
		}
		log.Printf("Choice (0 - %d): ", len(files))

		var (
			tmpString string
			err       error
		)

		fmt.Scan(&tmpString)

		choice, err = strconv.Atoi(tmpString)
		helpers.HandleFatalErr(err)
	}

	file := files[choice]

	dirName := strings.Replace(file.Filename, ".zip", "", 1)
	if dwc.packageManager.Contains(dirName, file.IdGamesUrl) {
		log.Println("skipping " + dirName + " it is already installed")
		return false
	}
	for _, mirror := range mirrors {
		installPath := saveContentToZipFile(file, mirror, dwc)

		unzipped := fmt.Sprintf("%s%s", dwc.Configuration.InstallDir, dirName)
		err := helpers.Unzip(installPath, unzipped)
		helpers.HandleFatalErr(err, "failed to unzip archive", installPath, "-")

		log.Println("Removing unnnecessary zip archive")
		err = os.Remove(installPath)
		helpers.HandleFatalErr(err, "failed to delete zip archive", installPath, "-")

		dwc.packageManager.NewEntry(dirName, unzipped, file.IdGamesUrl)

		return true
	}

	return false
}

func (dwc *Client) List() {
	for _, entry := range dwc.packageManager.entries {
		log.Printf("Package - Name: %s, Dir: %s, Uri: %s, Aliases: %s\n", entry.Name, entry.Dir, entry.Uri, entry.Aliases)
	}
}

func (dwc *Client) AddAlias(target, alias string) {
	dwc.packageManager.AddAlias(target, alias)
}

func (dwc *Client) ValidateCommand(cmd string) {
	if !helpers.Contains(validCommands, cmd) {
		err := errors.New(fmt.Sprint("invalid command, valid commands are ", validCommands))
		helpers.HandleFatalErr(err)
	}
}

func (dwc *Client) LookupLocalPath(name string) string {
	return dwc.packageManager.GetFilePath(name)
}

func saveContentToZipFile(file apiFile, mirror string, dwc *Client) string {
	filepath := strings.Replace(file.IdGamesUrl, idGamesSubstr, "", 1)
	endpointUri := fmt.Sprintf("http://%s/idgames/%s", mirror, filepath)

	log.Println("Attempting install from ", endpointUri)

	content, err := dwc.httpClient.Get(endpointUri)
	helpers.HandleFatalErr(err, "failed to retrieve file contents from mirror", mirror, "-")

	fmtdLocalPath := fmt.Sprint(dwc.Configuration.InstallDir, file.Filename)

	_, err = os.Create(fmtdLocalPath)
	helpers.HandleFatalErr(err, "failed to create file ", fmt.Sprint(dwc.Configuration.InstallDir, file.Filename), "-")

	bytes, _ := ioutil.ReadAll(content.Body)

	err = os.WriteFile(fmtdLocalPath, bytes, 0644)
	helpers.HandleFatalErr(err, "failed to write file -")

	log.Printf("Successfully wrote contents to %s\n", fmtdLocalPath)

	return fmtdLocalPath
}

func (dwc *Client) LookupIwad(name string) string {
	return dwc.packageManager.LookupIwad(name)
}

func (dwc *Client) RegisterIwad(name, iwad string) {
	dwc.packageManager.RegisterIwad(name, iwad)
}

func (dwc *Client) LookupWADAlias(alias string) string {
	return dwc.Configuration.IWads[alias]
}

func (dwc *Client) CollectPWads(dir string) [2]string {
	pwads := dwc.searchForWads(dir)

	targets := [2]string{}

	if len(pwads) > 0 {
		log.Println("Found some files in that archive:")
		for index, pwad := range pwads {
			log.Printf("%d) %s", index, pwad)
		}

		var selections string
		log.Print("Choose up to 2 (seperated by a comma): ")
		fmt.Scanln(&selections)

		if selections != "" {
			spl := strings.Split(selections, ",")

			if len(spl) > 1 {
				targetIndex, err := strconv.Atoi(spl[1])
				helpers.HandleFatalErr(err)

				targets[1] = pwads[targetIndex]
			}

			targetIndex, err := strconv.Atoi(spl[0])
			helpers.HandleFatalErr(err)

			targets[0] = pwads[targetIndex]
		}
	}

	return targets
}

func (dwc *Client) searchForWads(dir string) []string {
	var buffer []string
	var wads []string

	push(&buffer, dir)

	for len(buffer) > 0 {
		path := pop(&buffer)

		entries, err := os.ReadDir(path)
		helpers.HandleFatalErr(err)

		for _, entry := range entries {
			if entry.IsDir() {
				push(&buffer, fmt.Sprintf("%s/%s", path, entry.Name()))
			} else if isPWad(entry.Name()) {
				push(&wads, fmt.Sprintf("%s/%s", path, entry.Name()))
			}
		}
	}

	return wads
}

func push(buffer *[]string, item string) {
	*buffer = append(*buffer, item)
}

func pop(buffer *[]string) string {
	item := (*buffer)[len(*buffer)-1]

	*buffer = (*buffer)[:len(*buffer)-1]

	return item
}

func isPWad(name string) bool {
	spl := strings.Split(name, ".")

	return strings.ToLower(spl[len(spl)-1]) == "wad"
}

func sanitize(s string) string {
	var placeholder string = s

	for _, rne := range s {
		if !unicode.IsGraphic(rne) {
			placeholder = strings.ReplaceAll(placeholder, string(rne), "")
		}
	}
	return placeholder
}

func sanitizeFile(file apiFile) apiFile {
	var placeholder apiFile

	placeholder.Age = file.Age
	placeholder.Author = sanitize(file.Author)
	placeholder.Date = sanitize(file.Date)
	placeholder.Description = sanitize(file.Description)
	placeholder.Dir = sanitize(file.Dir)
	placeholder.Filename = sanitize(file.Filename)
	placeholder.IdGamesUrl = sanitize(file.IdGamesUrl)
	placeholder.Rating = file.Rating
	placeholder.Size = file.Size
	placeholder.Title = sanitize(file.Title)
	placeholder.Url = sanitize(file.Url)

	return placeholder
}
