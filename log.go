package main

import "log"

// init set configurations for log information
func init() {
	log.SetPrefix("DEBUG: ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
