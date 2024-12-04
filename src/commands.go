package main

import (
	"fmt"
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
	if len(cmd.args) == 0 {
		return fmt.Errorf("Login failed, please provide a username.")
	}

	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("Username has been successfully set to: %s\n", username)

	return nil
}
