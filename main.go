package main

import _ "github.com/lib/pq"
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/ryansurv/go-blog-aggregator/internal/config"
	"github.com/ryansurv/go-blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
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

// Sets the current user if exists
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) <= 0 {
		return errors.New("Login expects a username")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		os.Exit(1)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

// Registers a new user
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) <= 0 {
		return errors.New("Register expects a name")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return errors.New("User already exists")
		os.Exit(1)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: username,
	})
	if err != nil {
		return err
	}

	fmt.Println("User registered!")
	fmt.Println(user)

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	return nil
}

// Wipes the database of users
func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
	return nil // Go compiler complains if no return here
}

// Returns all users, appends "(current)" to logged-in user
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func main() {
	// Setup - State/commands
	cfg := config.Read()
	appState := state{cfg: cfg}
	cmds := commands{commands: make(map[string]func(*state, command) error)}

	// Setup - Database
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	appState.db = dbQueries
	
	// Register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

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