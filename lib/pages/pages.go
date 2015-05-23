/*
	Author:
		Nicholas Siow | compilewithstyle@gmail.com

	Description:
		This package handles all the logic for building the page cache
		and serving pages to HTTP requests

		Allows user to turn off caching for development purposes
*/

package pages

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//-----------------------------------------------------------------------------
//	program variables
//-----------------------------------------------------------------------------

// cache which maps URL -> page content -- allows for very fast page delivery
//   with no IO but server has to be restarted for page changes
var PageCache map[string][]byte

// variable to hold the string for the directory of the html pages
//   (mount point for filesystem access)
var HtmlDir string

// variable to hold HTML boilerplate that should be placed around each page
//   when they are read into the cache
var Boilerplate []byte

/*
	initialize filepaths and do the initial cache load
*/
func init() {

	// find the base directory of the executable
	BaseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Could not find base directory for executable. %s\n", err)
	}

	// join the filepath with `html` and make sure the directory containing
	//   html exists
	HtmlDir = filepath.Join(BaseDir, "html")
	if _, err := os.Stat(HtmlDir); err != nil {
		log.Fatalf("Could not find expected 'html/' directory at %s\n", HtmlDir)
	}

	// read in the boilerplate HTML for the site (banner, footer, etc)
	bp_file := filepath.Join(BaseDir, "boilerplate.html")
	Boilerplate, err = ioutil.ReadFile(bp_file)
	if err != nil {
		log.Fatalf("Could not read html boilerplate file at %s\n", bp_file)
	}

	// create an empty page cache
	PageCache = make(map[string][]byte)

	// call refresh_cache() to read the initial cache into memory
	refresh_cache(HtmlDir)

	log.Println("Page cache is loaded and server is ready to go. Hit CTRL^C to stop it")

}

/*
	function to re-read the files anf refresh the page cache
*/
func refresh_cache(html_dir string) {

	log.Println("Refreshing cache...")

	// declare a new temporary cache so that the cache replacement can be done
	//   in a single swap
	tempcache := make(map[string][]byte)

	// recursively find all files under the specified html top directory
	files, err := filepath.Glob(filepath.Join(html_dir, "**"))
	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range files {
		// the path of the file to read, NOT to be used as the key to the PageCache
		file2read := f

		// trim the leading html_dir and ending .html to get the
		//   path as it should appear in a URL
		f = strings.TrimPrefix(f, html_dir)
		f = strings.TrimSuffix(f, ".html")

		// don't map index.html files directly - this should be done
		//   using the directory name
		if strings.HasSuffix(file2read, "index.html") {
			continue
		}

		// determine if the path points to a file or directory
		fi, err := os.Stat(file2read)
		if err != nil {
			log.Println(err)
			tempcache[f] = make([]byte, 0)
			continue
		}

		// if the path is a directory, have it point to that directory's
		//   'index.html'
		if fi.IsDir() {
			file2read = filepath.Join(file2read, "index.html")
		}

		// map the path->filecontents in tempcache
		data, err := ioutil.ReadFile(file2read)
		if err != nil {
			log.Println(err)
			tempcache[f] = make([]byte, 0)
			continue
		}

		// insert the data into the boilerplate and save into the cache
		tempcache[f] = bytes.Replace(Boilerplate, []byte("REPLACEME"), data, 1)
	}

	// replace the old cache with the new, updated cache
	PageCache = tempcache
	log.Println("done refreshing cache")

}

/*
	base handler for serving html pages and directories
*/
func Handler(w http.ResponseWriter, r *http.Request) {

	// TODO extract some more information from the request to do analytics

	// extract the path of the request, after the domain name
	url := r.URL.Path

	// set the root path '/' to redirect to home
	if url == "/" {
		url = "/home"
	}

	// trim ending / if it exists
	url = strings.TrimSuffix(url, "/")

	// if the URL exists in the cache, send it. otherwise,
	//   send a 404
	page, exists := PageCache[url]
	if exists {
		// serve an error page if there is no content (script encountered a problem
		//   when trying to read the page into cache)
		if len(page) == 0 {
			http.Error(w, "Internal server error, sorry!", 500)
		} else {
			w.Write(page)
		}
	} else {
		http.NotFound(w, r)
	}

}
