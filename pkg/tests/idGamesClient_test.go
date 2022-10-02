package tests

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gowerm123/wadman/pkg/helpers"
	"github.com/gowerm123/wadman/pkg/idGamesClient"
)

const testQuery = "sunlust"
const sunlustText = "Filename: sunlust.zip"

type TestBuffer []byte

func (tb *TestBuffer) Write(p []byte) (n int, err error) {
	*tb = append(*tb, p...)

	return len(p), nil
}

func TestSearchValidTermReturnsResults(t *testing.T) {
	logBuffer := TestBuffer{}

	client := idGamesClient.New()

	log.SetOutput(&logBuffer)

	client.SearchAndPrint(testQuery)

	if len(logBuffer) == 0 {
		t.Fatal("expected log entries from stdout")
	}

	if !strings.Contains(string(logBuffer), sunlustText) {
		t.Fatal("expected files to be found")
	}
}

func TestInstallAndRemoveWorkProperly(t *testing.T) {
	client := idGamesClient.New()

	client.Install(testQuery)

	if _, err := os.Stat("" + testQuery); err != nil {
		t.Fatal(err.Error())
	}

	if len(client.GetFilePath(testQuery)) == 0 {
		t.Fatal("err expected filepath entry in wadmanifest")
	}

	client.Remove(testQuery)

	if _, err := os.Stat(helpers.GetHome() + "/.config/" + testQuery); err == nil {
		t.Fatal(helpers.GetHome() + "/.config/" + testQuery + "  should not exist after removal")
	}

	if len(client.GetFilePath(testQuery)) > 0 {
		t.Fatal("err pklist entry should have been removed")
	}
}
