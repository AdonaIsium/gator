package internal

import (
	"fmt"

	config "github.com/AdonaIsium/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*State, command) error
}

var cmds = commands{
	handlers: map[string]func(*State, command) error{
		"login": handlerLogin,
	},
}

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login command requires username to be provided")
	}

	err := s.Config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("username has been set to %s\n", s.Config.CurrentUserName)

	return nil
}

func (c *commands) run(s *State, cmd command) error {
	if thisCom, exists := cmds.handlers[cmd.name]; !exists {
		return fmt.Errorf("command %s does not exist", cmd.name)
	} else {
		err := thisCom(s, cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *commands) register(name string, f func(*State, command) error) {
	cmds.handlers[name] = f
}
