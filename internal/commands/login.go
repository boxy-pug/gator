package commands

import (
	"context"
	"fmt"

	"github.com/boxy-pug/gator/internal/config"
)

// Create a login handler function:
// This will be the function signature of all command handlers.
func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	userName := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("user does not exist: %v", err)
	}

	err = s.Config.SetUser(userName)
	if err != nil {
		return fmt.Errorf("error setting username: %v", err)
	}

	fmt.Printf("User %s has been set", userName)

	return nil
}

/*If the command's arg's slice is empty, return an error;
the login handler expects a single argument, the username.*/

/*Use the state's access to the config struct to set the user to the given username.
	Remember to return any errors.
    Print a message to the terminal that the user has been set.*/
