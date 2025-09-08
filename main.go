package main

import (
	"database/sql"
	"log"
	"os"

	config "github.com/AdonaIsium/gator/internal/config"
	"github.com/AdonaIsium/gator/internal/database"
	state "github.com/AdonaIsium/gator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	s := state.State{}
	s.Config = &cfg

	db, err := sql.Open("postgres", s.Config.DBURL)
	if err != nil {
		log.Fatalf("error opening sql database: %v", err)
	}

	dbQueries := database.New(db)
	s.DBQueries = dbQueries

	c := state.Commands{Handlers: map[string]func(*state.State, state.Command) error{}}
	c.Register("login", state.HandlerLogin)
	c.Register("register", state.HandlerRegister)
	c.Register("users", state.HandlerUsers)
	c.Register("reset", state.HandlerReset)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("command must be supplied")
	}
	cmd := state.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = c.Run(&s, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
