package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boxy-pug/gator/internal/config"
	"github.com/boxy-pug/gator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(*config.State, Command) error
}

// NewCommands creates and returns a new Commands instance.
func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*config.State, Command) error),
	}
}

// This method registers a new handler function for a Command name.
func (c *Commands) Register(name string, handler func(*config.State, Command) error) {
	c.handlers[name] = handler
}

// This method runs a given Command with the provided state if it exists.
func (c *Commands) Run(s *config.State, cmd Command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return errors.New("command not found")
	}
	return handler(s, cmd)
}

// HandlerReset deletes all users from the database.
func HandlerReset(s *config.State, cmd Command) error {
	// Execute the DeleteAllUsers query
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting database: %v", err)
	}

	fmt.Println("Database reset successfully.")
	return nil
}

func HandlerUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching users")
	}

	for _, name := range users {
		if name == s.Config.CurrentUserName {
			fmt.Printf("* %v (current)\n", name)
		} else {
			fmt.Printf("* %v\n", name)
		}
	}
	return nil

}

// HandlerAgg fetches an RSS feed and prints it .
func HandlerAgg(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("expected time between req argument")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

func MiddleWareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not logged in or does not exist: %w", err)
		}
		return handler(s, cmd, user)
	}
}
