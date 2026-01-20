package main

import _ "github.com/lib/pq"
import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ryansurv/go-blog-aggregator/internal/config"
	"github.com/ryansurv/go-blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
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
	cmds.register("agg", handlerAggregate)
	cmds.register("feeds", handlerFeeds)

	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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