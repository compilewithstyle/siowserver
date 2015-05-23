package minifier

import ()

/*
	function to perform the proper minification on the file content
	based on the given file extension. will minify the following types:

	(1) HTML
	(2) CSS
	(3) Javascript

	and will simply return the original content if no minifications can be done
*/
func Minify(fileext, content string) {
	if fileext == "html" || fileext == "css" || fileext == "js" {
		return _minify(content)
	} else {
		return content
	}
}

func _minify(content string) {

}
