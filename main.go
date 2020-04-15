package main

import (
	"log"

	"free-hls.go/cmd"
)

var (
	version string
	commit  string
	date    string
)

func init() {
	log.SetFlags(log.Llongfile | log.Ltime)
	cmd.Execute(version, commit, date)
}

func main() {}
