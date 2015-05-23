/*
	Author:
		Nicholas Siow | compilewithstyle@gmail.com

	Description:
		Replacement for pages.go that doesn't cache pages
		and offers helpful debugging features - to be used
		for development
*/

package devpages

import (
	"strings"
)

/*
	given a requested URL, returns the system filepath where the file
	would exist (does NOT check for existence though)
*/
func urlToFp(url string) string {

}

/*
	inserts the content of a page between the header and footer boilerplate
*/
func insertContents(page []byte) []byte {

}

/*
 */
func Handler(w http.ResponseWriter, r *http.Request) {
}
