package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/boxy-pug/gator/internal/config"
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
	feedURL := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()

	// Fetch the RSS feed
	feed, err := FetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	// Print the feed details
	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Link: %s\n", feed.Channel.Link)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Printf("Item Title: %s\n", item.Title)
		fmt.Printf("Item Link: %s\n", item.Link)
		fmt.Printf("Item Description: %s\n", item.Description)
		fmt.Printf("Item PubDate: %s\n", item.PubDate)
	}

	return nil
}
