package main

import (
	"log"

	"free-hls.go/cmd"
)

func init() {
	log.SetFlags(log.Llongfile | log.Ltime)
	cmd.Execute()
}

func main() {}
