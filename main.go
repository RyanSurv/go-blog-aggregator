package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ryansurv/go-blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) <= 0 {
		return errors.New("Login expects a username")
	}

	username := cmd.args[0]

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

func main() {
	// Setup
	cfg := config.Read()
	appState := state{cfg: cfg}
	cmds := commands{commands: make(map[string]func(*state, command) error)}
	
	// Register commands
	cmds.register("login", handlerLogin)

	// Get CLI info
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	fn := args[0]
	var fnArgs []string
	if len(args) > 1 {
		fnArgs = args[1:]
	}

	if err := cmds.run(&appState, command{name: fn, args: fnArgs}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}