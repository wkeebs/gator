package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wkeebs/gator/internal/database"
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
	id := uuid.New()

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
	// resets the state of the database - DEV ONLY
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset users: %s", err)
	}

	err = s.db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset feeds: %s", err)
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

func handlerAggregator(s *state, _ command) error {
	// main aggregator handler
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// adds a feed to the database
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: addfeed <name> <url>")
	}

	// add feed to database
	name := cmd.args[0]
	url := cmd.args[1]

	id := uuid.New()
	now := time.Now().UTC()

	params := database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create feed: %s", err)
	}

	// automatically create a feed follows entry
	new_args := []string{url}
	new_cmd := command{
		name: "follow",
		args: new_args,
	}
	err = handlerFollow(s, new_cmd, user)
	if err != nil {
		return fmt.Errorf("Failed to add feed follows: %s", err)
	}

	fmt.Println("Feed added succesfully.")
	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, _ command) error {
	// lists all feeds
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds: %s", err)
	}

	for _, f := range feeds {
		fmt.Printf("%s:\n", f.Name)
		fmt.Printf("* URL: %s\n", f.Url)
		fmt.Printf("* Created By: %s\n", f.UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	// add a feed to the current user's following
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: follow <url>")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to get feed by url: %s", err)
	}

	now := time.Now().UTC()
	id := uuid.New()
	params := database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		ID_2:      user.ID,
		ID_3:      feed.ID,
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create feed_follow: %s", err)
	}

	fmt.Println(feed_follow)

	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	// gets all feeds followed by a user
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get feed follows for user: %s", err)
	}

	for _, f := range follows {
		fmt.Println(f)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	// unfollow a feed for the current user
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: unfollow <url>")
	}

	url := cmd.args[0]
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	}
	err := s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to unfollow: %s", err)
	}

	fmt.Printf("Succesfully unfollowed %s.\n", url)

	return nil
}
