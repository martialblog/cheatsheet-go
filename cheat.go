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
	"github.com/gookit/color"
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

// Checks for the error and prints it to the console
func Log(err error) {
	if err != nil {
		color.Warn.Println(err)
	}
}

// Checks for the error and exists if there is one
func ExitIf(err error) {
	if err != nil {
		color.Error.Println(err)
		os.Exit(1)
	}
}

// Reads the JSON file at the path and returns a Cheatsheet struct
func readCheatsheet(path string) (result Cheatsheet, err error) {
	// Read JSON file
	file, err := ioutil.ReadFile(path)
	ExitIf(err)

	// Parse JSON and return
	if err := json.Unmarshal(file, &result); err != nil {
		Log(err)
	}

	return result, err
}

// Checks for a directory and creates it if necessary
func checkOrCreateDir(path string) {
	// Check if path exists and create directory if not
	if err := os.MkdirAll(path, 0755); err != nil {
		ExitIf(err)
	}
}

// Loads a remote file to a directory
func loadRemoteCheatsheetToCache(url string, filepath string) error {

	// Load from URL
	resp, err := http.Get(url)
	ExitIf(err)

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	ExitIf(err)

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// Prints out the Cheatsheet struct to console
func printCheatsheet(sheet Cheatsheet, colored bool) {

	k := color.FgWhite.Render
	v := color.FgWhite.Render

	if colored {
		k = color.FgCyan.Render
		v = color.FgGreen.Render
	}

	for key, value := range sheet.Cheats {
		fmt.Printf("%s \t %s\n", k(key), v(value))
	}
}

// Where the magic happens
func main() {
	flag.StringVar(&sheet, "sheet", "", "Cheatsheet to display")
	flag.StringVar(&sheetDir, "sheetdir", "examples", "Directory of local Cheatsheets")
	flag.StringVar(&cacheDir, "cachedir", "/tmp/cheat", "Cache directory of remote Cheatsheets")
	flag.BoolVar(&colored, "colored", false, "Display colored output")
	flag.Parse()

	Sheet, err := readCheatsheet(path.Join(sheetDir, sheet + ".json"))
	ExitIf(err)

	if Sheet.Remote {
		cachedFile := path.Join(cacheDir, sheet + ".json")
		url := Sheet.Url

		if _, err := os.Stat(cachedFile); os.IsNotExist(err) {
			// File not in cache, trying to download it...
			checkOrCreateDir(cacheDir)
			loadRemoteCheatsheetToCache(url, cachedFile)
		}
		Sheet, err = readCheatsheet(cachedFile)
		ExitIf(err)
	}

	printCheatsheet(Sheet, colored)

	os.Exit(0)
}
