package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	config "github.com/AdonaIsium/gator/internal/config"
	"github.com/AdonaIsium/gator/internal/database"
	rss "github.com/AdonaIsium/gator/internal/rss"
	"github.com/google/uuid"
)

type State struct {
	Config    *config.Config
	DBQueries *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login command requires username to be provided")
	}

	_, err := s.DBQueries.GetUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("username has been set to %s\n", s.Config.CurrentUserName)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("register command requires username to be provided")
	}

	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0]}

	dbUser, err := s.DBQueries.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatalf("error creating user: %v", err)
	}

	s.Config.CurrentUserName = dbUser.Name
	s.Config.SetUser(dbUser.Name)

	fmt.Printf("created user: %s", dbUser.Name)

	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DBQueries.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("unable to get all users: %v", err)
	}

	currentUser, err := s.DBQueries.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		log.Fatalf("current user not in db: %v", err)
	}

	for _, user := range users {
		if currentUser.ID == user.ID {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	rssFeed, err := rss.FetchFeed(context.Background(), "https://wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", rssFeed)

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		log.Fatalf("exactly 2 arguments (name, url) required to create feed.")
	}
	currentUser, err := s.DBQueries.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}

	createFeedParams := database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0], Url: cmd.Args[1], UserID: currentUser.ID}

	feed, err := s.DBQueries.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %s, CreateAt: %v, UpdatedAt: %v, Name: %s, URL: %s, UserID: %s\n", feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)

	return nil
}

func HandlerGetFeeds(s *State, cmd Command) error {
	feeds, err := s.DBQueries.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.DBQueries.GetUser(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s, URL: %s, User: %s", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.DBQueries.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatalf("error resetting database: %v", err)
	}

	fmt.Printf("database reset successfully\n")
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	if thisCom, exists := c.Handlers[cmd.Name]; !exists {
		return fmt.Errorf("command %s does not exist", cmd.Name)
	} else {
		err := thisCom(s, cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}
