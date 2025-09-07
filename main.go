package main

import (
	"fmt"
	"log"

	config "github.com/AdonaIsium/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	err = cfg.SetUser("Lanre")
	if err != nil {
		log.Fatalf("error setting username: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
	fmt.Printf("db_url: %s, current_user_name: %s", cfg.DBURL, cfg.CurrentUserName)
}
