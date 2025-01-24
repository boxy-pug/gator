package commands

import (
	"errors"

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
