package main

import (
	"log"

	"github.com/alenn-m/rgen/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd.Execute()
}
