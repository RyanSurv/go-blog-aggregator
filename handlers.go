package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/ryansurv/go-blog-aggregator/internal/database"
)

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

// Wipes the database completely
func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		os.Exit(1)
	}
	if err := s.db.DeleteFeeds(context.Background()); err != nil {
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

// Long-running aggregator service (temporarily set to fetch a single feed)
func handlerAggregate(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

// Creates a new feed
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Not enough arguments")
	}

	// Create the feed
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	// Automatically follow the newly created feed
	if _, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

// Lists all feeds
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed{\n  Name: %s,\n  URL: %s,\n  Username: %s,\n}\n", feed.Name, feed.Url, feed.Name_2.String)
	}

	return nil
}

// Creates a new feed follow
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("Not enough arguments")
	}

	// Get the feed from the URL provided
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s followed %s\n", ff.UserName, ff.FeedName)
	return nil
}

// Lists all names of feeds logged in user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
	ff, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	str := "["
	for i, f := range ff {
		str += fmt.Sprintf("%s", f.FeedName)
		if i != len(ff) - 1 {
			str += ", "
		}
	} 
	str += "]"
	fmt.Println(str)

	return nil
}

// Unfollowed a specified feed from logged in user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("Not enough arguments")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}

	if err := s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return err
	}

	fmt.Printf("%s unfollowed %s\n", user.Name, url)
	return nil
}