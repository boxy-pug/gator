package main

import (
	"fmt"
	"log"

	"github.com/boxy-pug/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	fmt.Println(c.DbUrl, c.CurrentUserName)

	err = c.SetUser("Larry")
	if err != nil {
		log.Fatalf("error setting user: %v", err)
	}

	fmt.Println(c.DbUrl, c.CurrentUserName)
}
