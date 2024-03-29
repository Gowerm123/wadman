package helpers

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func HandleFatalErr(err error, msgs ...string) {
	if err != nil {
		log.Println(msgs, err.Error())
		os.Exit(1)
	}
}

func Unzip(src, dest string) error {
	log.Printf("unzipping %s to %s \n", src, dest)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0755)

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func Contains(data []string, tgt string) bool {
	for _, str := range data {
		if str == tgt {
			return true
		}
	}
	return false
}

func Split(data string) []string {
	return strings.Split(data, ",")
}

func IsRoot() bool {
	currentUser, err := user.Current()
	HandleFatalErr(err)

	return currentUser.Username == "root"
}

func GetWadmanHomeDir() string {
	const wadmanPath = "/.wadman/"
	if IsRoot() {
		username := os.Getenv("SUDO_USER")
		u, err := user.Lookup(username)
		HandleFatalErr(err)

		return u.HomeDir + wadmanPath
	}
	dir, err := os.UserHomeDir()
	HandleFatalErr(err)

	return dir + wadmanPath
}

func WadmanConfigPath() string {
	const configPath = "/.config/wadman-config.json"
	if IsRoot() {
		username := os.Getenv("SUDO_USER")
		u, err := user.Lookup(username)
		HandleFatalErr(err)

		return u.HomeDir + configPath
	}
	dir, err := os.UserHomeDir()
	HandleFatalErr(err)

	return dir + configPath
}

func ToZipFileName(filename string) string {
	if len(filename) < 4 || filename[len(filename)-4:] != ".zip" {
		return filename + ".zip"
	}
	return filename
}

func ToTargetName(filename string) string {
	if len(filename) >= 4 && filename[len(filename)-4:] == ".zip" {
		return filename[:len(filename)-4]
	}
	return filename
}
