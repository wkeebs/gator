package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wkeebs/rss-blog-aggregator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.cmdMap[cmd.name]
	if !exists {
		return fmt.Errorf("Command '%s' has not been registered.", cmd.name)
	}
	return f(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	// logs a user into the application
	if len(cmd.args) == 0 {
		return fmt.Errorf("Login failed, please provide a username.")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Failed to log in, user does not exist: %s", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to log in: %s", err)
	}

	fmt.Printf("Username has been successfully set to: %s\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	// registers a user in the database
	if len(cmd.args) == 0 {
		return fmt.Errorf("Register failed, please provide a name.")
	}

	name := cmd.args[0]
	now := time.Now().UTC()
	id, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("Failed to generate id : %s", err)
	}

	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to register user '%s': %s", name, err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Failed to register user '%s': %s", name, err)
	}

	fmt.Printf("User '%s' was created successfully.\n", name)
	fmt.Println(user)

	return nil
}

func handlerReset(s *state, _ command) error {
	// resets the state of the users in the database - DEV ONLY
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset: %s", err)
	}
	return nil
}

func handlerUsers(s *state, _ command) error {
	// lists all users
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get users: %s", err)
	}

	for _, u := range users {
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}

	return nil
}


func handleAggregator(s *state, _ command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}