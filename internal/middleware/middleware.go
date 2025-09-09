package internal

import (
	"context"

	"github.com/AdonaIsium/gator/internal/database"
	state "github.com/AdonaIsium/gator/internal/state"
)

func MiddlewareLoggedIn(handler func(s *state.State, cmd state.Command, user database.User) error) func(*state.State, state.Command) error {
	return func(s *state.State, cmd state.Command) error {
		user, err := s.DBQueries.GetUserByName(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
