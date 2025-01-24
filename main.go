package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boxy-pug/gator/internal/commands"
	"github.com/boxy-pug/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	// Initialize state
	appState := &config.State{
		Config: &c,
	}

	// Initialize commands and handlers
	cmds := commands.NewCommands()
	cmds.Register("login", commands.HandlerLogin)

	//If there are fewer than 2 arguments, print an error message to the terminal and exit. Why two? The first argument is automatically the program name, which we ignore, and we require a command name.
	if len(os.Args) < 2 {
		log.Fatalf("command name required")
	}

	cmd := commands.Command{Name: os.Args[1], Args: os.Args[2:]}
	err = cmds.Run(appState, cmd)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
