package tests

import (
	"log"
	"os"
	"strings"
	"testing"

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

	client.SearchAndPrint(testQuery, "title")

	if len(logBuffer) == 0 {
		t.Fatal("expected log entries from stdout")
	}

	if !strings.Contains(string(logBuffer), sunlustText) {
		t.Fatal("expected files to be found")
	}
}

func TestInstallAndRemoveWorkProperly(t *testing.T) {
	client := idGamesClient.New()

	client.Install(testQuery, "title")

	if _, err := os.Stat("/usr/share/wadman/" + testQuery); err != nil {
		t.Fatal(err.Error())
	}

	if len(client.GetFilePath(testQuery)) == 0 {
		t.Fatal("err expected filepath entry in pkglist")
	}

	client.Remove(testQuery)

	if _, err := os.Stat("/usr/share/wadman/" + testQuery); err == nil {
		t.Fatal("/usr/share/wadman/" + testQuery + "  should not exist after removal")
	}

	if len(client.GetFilePath(testQuery)) > 0 {
		t.Fatal("err pklist entry should have been removed")
	}
}
