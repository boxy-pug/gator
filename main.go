package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/boxy-pug/gator/internal/commands"
	"github.com/boxy-pug/gator/internal/config"
	"github.com/boxy-pug/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	//load in your database URL to the config struct
	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		log.Fatalf("error loading in database")
	}
	defer db.Close()
	dbQueries := database.New(db)

	// Initialize state
	appState := &config.State{
		Db:     dbQueries,
		Config: &c,
	}

	// Initialize commands and handlers
	cmds := commands.NewCommands()
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("addfeed", commands.HandlerAddFeed)
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("follow", commands.HandlerFollow)
	cmds.Register("following", commands.HandlerFollowing)

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
