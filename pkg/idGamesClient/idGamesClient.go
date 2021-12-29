package idGamesClient

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	fmt.Printf("File Found\n    Filename: %s\n    Title: %s,\n    Author: %s,\n    Date: %s\n    Url: %s\n", file.Filename, file.Title, file.Author, file.Date, file.IdGamesUrl)
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

	if len(files) == 0 {
		fmt.Print("Entry not found for search query, try a different QUERYTYPE?")
		os.Exit(0)
	}

	file := files[0]

	dirName := strings.Replace(file.Filename, ".zip", "", 1)
	if dwc.packageManager.Contains(dirName, file.IdGamesUrl) {
		fmt.Println("skipping " + dirName + " it is already installed")
		return false
	}
	for _, mirror := range mirrors {
		installPath := saveContentToZipFile(file, mirror, dwc)

		unzipped := fmt.Sprintf("%s%s", dwc.Configuration.InstallDir, dirName)
		err := helpers.Unzip(installPath, unzipped)
		helpers.HandleFatalErr(err, "failed to unzip archive", installPath, "-")

		fmt.Println("Removing unnnecessary zip archive")
		err = os.Remove(installPath)
		helpers.HandleFatalErr(err, "failed to delete zip archive", installPath, "-")

		dwc.packageManager.NewEntry(dirName, unzipped, file.IdGamesUrl)

		return true
	}

	return false
}

func (dwc *Client) List() {
	for _, entry := range dwc.packageManager.entries {
		fmt.Printf("Package - Name: %s, Dir: %s, Uri: %s, Aliases: %s\n", entry.Name, entry.Dir, entry.Uri, entry.Aliases)
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

	fmt.Println("Attempting install from ", endpointUri)

	content, err := dwc.httpClient.Get(endpointUri)
	helpers.HandleFatalErr(err, "failed to retrieve file contents from mirror", mirror, "-")

	fmtdLocalPath := fmt.Sprint(dwc.Configuration.InstallDir, file.Filename)

	_, err = os.Create(fmtdLocalPath)
	helpers.HandleFatalErr(err, "failed to create file ", fmt.Sprint(dwc.Configuration.InstallDir, file.Filename), "-")

	bytes, _ := ioutil.ReadAll(content.Body)

	err = os.WriteFile(fmtdLocalPath, bytes, 0644)
	helpers.HandleFatalErr(err, "failed to write file -")

	fmt.Printf("Successfully wrote contents to %s\n", fmtdLocalPath)

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
