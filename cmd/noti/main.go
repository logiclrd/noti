package main

import (
	"log"

	"github.com/logiclrd/noti/internal/command"
)

func main() {
	if err := command.Root.Execute(); err != nil {
		log.Fatal(err)
	}
}
