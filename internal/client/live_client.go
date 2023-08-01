package client

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

	"github.com/gowerm123/wadman/internal/helpers"
)

const (
	idGamesBaseURI = "http://www.doomworld.com/idgames/api/api.php"
	idGamesSubstr  = "idgames://"
)

var queryTypes = []string{"filename", "title", "author"}

var (
	validCommands []string = []string{"search", "install", "run", "help", "list", "alias", "register", "configure", "remove"}
)

type LiveClient struct {
	httpClient http.Client
	Configuration
}

type Payload struct {
	Body    string
	Headers http.Header
}

type searchResponse struct {
	Files []*ApiFile `xml:"content>file"`
}

type ApiFile struct {
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

func New() IdGamesClient {
	var client LiveClient

	client.Configuration = loadConfigs()
	client.httpClient = *http.DefaultClient

	return client
}

func (dwc LiveClient) sendQuery(query string) searchResponse {
	response := searchResponse{}
	for _, queryType := range queryTypes {
		pl, err := dwc.dial("search", map[string]string{"query": query, "type": queryType})
		helpers.HandleFatalErr(err, "failed to send search query")
		tempResponse := searchResponse{}
		xml.Unmarshal([]byte(pl.Body), &tempResponse)

		response.Files = append(response.Files, tempResponse.Files...)
	}
	return response
}

func (dwc LiveClient) search(query string) []*ApiFile {
	response := dwc.sendQuery(query)
	return response.Files
}

func (dwc LiveClient) Search(query string) string {
	response := dwc.sendQuery(query)
	output := ""
	for _, entry := range response.Files {
		output += formatAndPrint(entry)
	}
	return output
}

func formatAndPrint(file *ApiFile) string {
	file = sanitizeFile(file)
	return fmt.Sprintf("[%s] - %s\n\t%s\n\n", helpers.ToTargetName(file.Filename), file.Title, file.Description)
}

func (dwc LiveClient) dial(action string, params map[string]string) (Payload, error) {
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

func collectChoice(dwc LiveClient, query string) *ApiFile {
	if len(query) < 3 && len(query) > 0 {
		query = helpers.ToZipFileName(query)
	}
	files := dwc.search(query)
	filtered := []*ApiFile{}
	for _, file := range files {
		if file.Filename == query {
			filtered = append(filtered, file)
		}
	}

	files = dedup(files)

	var choice int = 0

	if len(files) == 0 {
		log.Printf("no archives found for search %s\n", query)
		os.Exit(0)
	} else if len(files) > 1 {
		if len(filtered) == 1 {
			return filtered[0]
		}
		log.Println("Multiple files found, please choose...")
		for it, file := range files {
			file = sanitizeFile(file)
			log.Printf("%d) %s, by %s, file - %s\n", it+1, file.Title, file.Author, file.Filename)
		}
		log.Printf("Choice (1 - %d): ", len(files))

		var (
			tmpString string
			err       error
		)

		fmt.Scan(&tmpString)

		choice, err = strconv.Atoi(tmpString)
		helpers.HandleFatalErr(err)
		choice -= 1
	}

	return files[choice]
}

func dedup(ls []*ApiFile) []*ApiFile {
	keys := make(map[string]bool)
	outLs := []*ApiFile{}

	for _, val := range ls {
		if _, exists := keys[val.IdGamesUrl]; !exists {
			outLs = append(outLs, val)
			keys[val.IdGamesUrl] = true
		}
	}
	return outLs
}

func (dwc LiveClient) Install(query string, am ArchiveManager) bool {
	file := collectChoice(dwc, query)

	dwc.saveContentToZipFile(file)
	am.Install(*file)

	return true
}

func (dwc LiveClient) Set(key, value string) {
	dwc.Configuration.Update(key, value)
}

func (dwc LiveClient) ValidateCommand(cmd string) {
	if !helpers.Contains(validCommands, cmd) {
		err := errors.New(fmt.Sprint("invalid command, valid commands are ", validCommands))
		helpers.HandleFatalErr(err)
	}
}

func (dwc LiveClient) InstallByFile(file *ApiFile, am ArchiveManager) {
	dwc.saveContentToZipFile(file)
	am.Install(*file)
}

func (dwc LiveClient) saveContentToZipFile(file *ApiFile) {
	for _, mirror := range dwc.Configuration.Mirrors {
		filepath := strings.Replace(file.IdGamesUrl, idGamesSubstr, "", 1)
		endpointUri := fmt.Sprintf("http://%s/idgames/%s", mirror, filepath)
		log.Println("Attempting install from ", endpointUri)

		content, err := dwc.httpClient.Get(endpointUri)
		helpers.HandleFatalErr(err, "failed to retrieve file contents from mirror", mirror, "-")

		fmtdLocalPath := helpers.GetWadmanHomeDir() + file.Filename

		_, err = os.Create(fmtdLocalPath)
		helpers.HandleFatalErr(err, "failed to create file ", helpers.GetWadmanHomeDir()+file.Filename, "-")

		bytes, _ := ioutil.ReadAll(content.Body)

		err = os.WriteFile(fmtdLocalPath, bytes, 0644)
		helpers.HandleFatalErr(err, "failed to write file -")

		log.Printf("Successfully wrote contents to %s\n", fmtdLocalPath)
		return
	}
	panic(fmt.Errorf("failed to download from available mirrors, mirrors can be updated at $HOME/.config/wadman-config.json"))
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

func sanitizeFile(file *ApiFile) *ApiFile {
	var placeholder *ApiFile = &ApiFile{}

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
