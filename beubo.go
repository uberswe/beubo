package beubo

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var port = 3000

// Init is called to start Beubo, this calls various other functions that initialises some basic settings
func Init() {
	readCLIFlags()
}

// Run runs the main application
func Run() {
	checkFiles()
	settingsInit()
	databaseInit()
	databaseSeed()
	loadPlugins()
	routesInit()
}

// readCLIFlags parses command line flags such as port number
func readCLIFlags() {
	flag.IntVar(&port, "port", port, "The port you would like the application to listen on")
	flag.Parse()
}

// checkFiles checks if there is a theme present and will otherwise prompt to download that theme
func checkFiles() {
	_, err := os.Stat("./themes")
	if os.IsNotExist(err) {
		if ask("There is no themes folder present, would you like to create it?") {
			err := os.Mkdir("./themes", 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if !hasDirectory("./themes/") {
		if ask("There is no theme installed, would you like to download the default theme?") {
			resp, err := http.Get("https://github.com/uberswe/beubo-default/archive/master.zip")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				log.Fatal(err)
			}

			// Read all the files from zip archive
			for _, zipFile := range zipReader.File {
				err := extractAndWriteFile(zipFile)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
	_, err = os.Stat("./themes/install")
	if os.IsNotExist(err) {
		if ask("There is no install theme installed, would you like to download the default install theme?") {
			resp, err := http.Get("https://github.com/uberswe/beubo-install/archive/master.zip")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				log.Fatal(err)
			}

			// Read all the files from zip archive
			for _, zipFile := range zipReader.File {
				err := extractAndWriteFile(zipFile)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

// ask prompts the user for an answer to a question via the command line. Returns true for yes and false for no
func ask(question string) bool {
	prompt := promptui.Select{
		Label: fmt.Sprintf("%s [Yes/No]", question),
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false
	}
	if result == "Yes" {
		return true
	}
	return false
}

// hasDirectory takes a path and checks if a directory exists in the path
func hasDirectory(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}
	for _, f := range files {
		if f.IsDir() {
			return true
		}
	}
	return false
}

func extractAndWriteFile(f *zip.File) error {
	dest := "./themes/"
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	filename := strings.Replace(f.Name, "beubo-default-master/", "default/", 1)
	filename = strings.Replace(filename, "beubo-install-master/", "install/", 1)

	path := filepath.Join(dest, filename)

	// Check for ZipSlip (Directory traversal)
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
