/*
	Author:
		Nicholas Siow | compilewithstyle@gmail.com

	Description:
		This module manages site configuration options and
		global site tasks
*/

package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//-----------------------------------------------------------------------------
//	struct to hold the JSON config options
//-----------------------------------------------------------------------------

type SiteConfig struct {
	Root      string
	DevPort   int
	LogFile   string
	AdminPort int
}

// global variable for site configuration, should be available
// for all modules to use
var Site SiteConfig

//-----------------------------------------------------------------------------
//	module-level variables
//-----------------------------------------------------------------------------

// filepath to look for config file
var fp string = "/usr/local/etc/siow/config.json"

// module logger
var l *log.Logger

// log file stream
var logFileStream *File

// logger for initial setup errors
var errLog *log.Logger = log.New(os.Stderr, "[ERROR]:", 50)

/*
	looks for the config file and builds the site config according to user parameters
*/
func init() {

	// make sure the config file exists
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		errLog.Fatalln("Could not find config file @", fp)
	}

	// TODO - offer the option to create with default settings?

	// open the config file for reading
	file, err := os.Open(fp)
	if err != nil {
		errLog.Fatalf("Could not open config file @ %s: %v\n", fp, err)
	}

	// read in contents
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		errLog.Fatalf("Could not read config file @ %s: %v\n", fp, err)
	}

	// parse JSON into config struct
	err = json.Unmarshal(contents, &Site)
	if err != nil {
		errLog.Fatalf("Problem parsing JSON file %s: %v\n", fp, err)
	}

	// make sure the configuration is valid
	verify(&Site)

}

/*
	verify that all fields in the configuration are valid
*/
func verify(s *SiteConfig) {

	// make sure all required fields exist
	if s.Root == "" {
		errLog.Fatalf("Missing field <%s> in JSON config file: %s\n", "Root", fp)
	}

	if s.LogFile == "" {
		errLog.Fatalf("Missing field <%s> in JSON config file: %s\n", "LogFile", fp)
	}

	if s.LogMask == 0 {
		errLog.Fatalf("Missing or zero-value field <%s> in JSON config file: %s\n", "LogMask", fp)
	}

	if s.AdminPort == 0 {
		errLog.Fatalf("Missing or zero-value field <%s> in JSON config file: %s\n", "AdminPort", fp)
	}

	if s.DevPort == 0 {
		errLog.Fatalf("Missing or zero-value field <%s> in JSON config file: %s\n", "DevPort", fp)
	}

	// make sure the web root was declared and is a valid directory
	fi, err := os.Stat(s.Root)
	if os.IsNotExist(err) {
		errLog.Fatalln("Specified web root directory does not exist:", s.Root)
	}

	if !fi.IsDir() {
		errLog.Fatalln("Specified web root directory is not a directory:", s.Root)
	}

	// make sure that the given ports are in valid ranges
	if s.AdminPort < 1024 || s.AdminPort > 65536 {
		errLog.Fatalln("Invalid port value for AdminPort", s.AdminPort)
	}

	if s.DevPort < 1024 || s.AdminPort > 65536 {
		errLog.Fatalln("Invalid port value for DevPort", s.AdminPort)
	}

	if s.DevPort == s.AdminPort {
		errLog.Fatalln("AdminPort value is the same as DevPort - should be different")
	}

	// make sure it's possible to write to the specified logfile
	file, err := os.OpenFile(s.LogFile, os.O_APPEND, os.ModeAppend)
	if err != nil {
		errLog.Fatalln("Error opening log file:", err)
	}

	// set the log file stream to the opened file
	logFileStream = file

	// TODO - register closing of the log file with the exit handler

}

/*
	returns a new logger that is pre-configured with the site logging
	options
*/
func NewLogger(prefix string) *log.Logger {
	return &log.New(logFileStream, prefix, 50)
}
