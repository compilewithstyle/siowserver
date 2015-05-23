/*
	Author:
		Nicholas Siow | compilewithstyle@gmail.com

	Description:
		Handles logic for reading in CSS templates, altering the color scheme, and
		moving the files to a css folder to be served statically
*/

package colors

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// global color map variable, containing mapping of color names to values
var color_map map[string]string

// location of folders for raw/processed CSS
var (
	RawDir  string
	ProcDir string
)

/*
	set up the color map
*/
func init() {

	// populate the color map
	color_map = map[string]string{
		"background": "DEDEDE",
		"mytext":     "red",
	}

	// determine the folders for raw and processed CSS
	BaseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Could not find base directory for executable. %s\n", err)
	}

	RawDir = filepath.Join(BaseDir, "static", "raw_css")
	ProcDir = filepath.Join(BaseDir, "static", "css")

}

/*
   recolors a single .CSS file, to be used as WalkFunc
*/
func _recolor(path string, info os.FileInfo, err error) error {

	// skip unless it's a CSS file
	if !strings.HasSuffix(path, ".css") {
		return nil
	}

	log.Println("Recoloring file: " + path)

	// separate the directory string from the file name
	_, file_name := filepath.Split(path)

	// read the contents of the raw css file into a string
	file_bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	file_contents := string(file_bytes)

	// make all replacements
	for code, color := range color_map {
		file_contents = strings.Replace(file_contents, "{"+code+"}", "#"+color, -1)
	}

	// write out to processed css directory
	new_file := filepath.Join(ProcDir, file_name)
	err = ioutil.WriteFile(new_file, []byte(file_contents), 0644)

	log.Println("Wrote file to: " + new_file)

	return err

}

/*
   uses the color variables to read in css templates and fill in the various colors
*/
func Recolor() {

	log.Println("Starting recolor process on directory " + RawDir)

	// walk through the directory of raw css and process each of them
	err := filepath.Walk(RawDir, _recolor)
	if err != nil {
		panic(err)
	}

	log.Println("Finished recolor process on directory " + ProcDir)

}
