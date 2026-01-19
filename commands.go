package main

import (
	"errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.commands[cmd.name]; ok {
		if err := c.commands[cmd.name](s, cmd); err != nil {
			return err
		}
	} else {
		return errors.New("Command not found")
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.commands[name]; !ok {
		c.commands[name] = f
	}
}