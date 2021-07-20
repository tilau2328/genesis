package main

import (
	"github.com/tilau2328/genesis/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
