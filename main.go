package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/AdonaIsium/gator/internal/config"
	state "github.com/AdonaIsium/gator/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	s := state.State{}

	s.Config = &cfg

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("command must be supplied")
	}

	fmt.Printf("db_url: %s, current_user_name: %s", s.Config.DBURL, s.Config.CurrentUserName)
}
