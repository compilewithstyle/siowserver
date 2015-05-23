/*
	Author:
		Nicholas Siow | compilewithstyle@gmail.com

	Description:
		FIXME
*/

package main

import (
	"github.com/compilewithstyle/siowserver/lib/pages"
	"net/http"
)

/*
	start listening and serving requests
*/
func main() {

	// serve static css files at /css/*
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// set up the main page handler
	http.HandleFunc("/", pages.Handler)
	http.ListenAndServe(":8080", nil)

}
