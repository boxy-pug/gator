package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boxy-pug/gator/internal/config"
	"github.com/boxy-pug/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the register command expects a username")
	}

	userName := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), userName)
	if err == nil {
		return fmt.Errorf("user already exists: %v", err)
	}

	userId := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	_, err = s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        userId,
		Name:      userName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	// set current user in config:
	err = s.Config.SetUser(userName)
	if err != nil {
		fmt.Errorf("error setting username in config")
	}
	fmt.Printf("User %s has been created successfully\n", userName)
	log.Printf("User created: %v, %v, %v\n", userId, userName, createdAt)
	return nil
}
