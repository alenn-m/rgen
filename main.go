package main

import (
	"log"

	"github.com/alenn-m/rgen/v2/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd.Execute()
}
