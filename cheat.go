package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"net/http"
)

var (
	sheet string
	colored bool
	sheetDir string
	cacheDir string
)

type Cheatsheet struct {
	Name string
	Remote bool
	Url string
	Cheats map[string] string
}

func readCheatsheet(path string) (result Cheatsheet, err error) {
	// Read JSON file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	// Parse JSON and return
	if err := json.Unmarshal(file, &result); err != nil {
		fmt.Println(err)
	}

	return result, err
}

func checkOrCreateDir(path string) error {
	var err error

	// Check if path exists and create directory if not
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Println(err)
	}

	return err
}

func loadRemoteCheatsheetToCache(url string, filepath string) error {

	// Load from URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func printCheatsheet(sheet Cheatsheet) {
	fmt.Println(sheet.Name)
	fmt.Println(sheet.Remote)
	fmt.Println(sheet.Cheats)
}

func main() {
	flag.StringVar(&sheet, "sheet", "", "Cheatsheet to display")
	flag.StringVar(&sheetDir, "sheetdir", "examples", "Directory of local Cheatsheets")
	flag.StringVar(&cacheDir, "cachedir", "/tmp/cheat", "Cache directory of remote Cheatsheets")
	flag.BoolVar(&colored, "colored", true, "Display colored output")
	flag.Parse()

	Sheet, err := readCheatsheet(path.Join(sheetDir, sheet + ".json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if Sheet.Remote {
		cachedFile := path.Join(cacheDir, sheet + ".json")
		url := Sheet.Url

		Sheet, err = readCheatsheet(cachedFile)
		// File probably not in cache, downloading it...
		if err != nil {
			checkOrCreateDir(cacheDir)
			if err != nil {
				// Chould not create cache
				fmt.Println(err)
				os.Exit(1)
			}

			loadRemoteCheatsheetToCache(url, cachedFile)
			Sheet, err = readCheatsheet(cachedFile)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	printCheatsheet(Sheet)

	os.Exit(0)
}
