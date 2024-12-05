package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // postgres driver: ignore
	"github.com/wkeebs/rss-blog-aggregator/internal/config"
	"github.com/wkeebs/rss-blog-aggregator/internal/database"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// connect to database
	db, err := sql.Open("postgres", configStruct.DbUrl)
	dbQueries := database.New(db)

	// initialise state and commands
	appState := state{
		db:  dbQueries,
		cfg: &configStruct,
	}
	appCommands := commands{cmdMap: make(map[string]func(*state, command) error)}
	appCommands.register("login", handlerLogin)

	// get arguments from user
	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("Not enough arguments, please provide a command.")
		os.Exit(1)
	}

	name, args := userArgs[1], userArgs[2:]

	userCommand := command{name: name, args: args}
	err = appCommands.run(&appState, userCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(appState.cfg)
	os.Exit(0)
}
