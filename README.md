# Gator

Gator is a blog aggregator written in Go!

## Requirements

You will need Go installed to build the app. You will also need Postgres to run the app.

## Installation

In the root directory run `go install`

## Setup & Running

1. Create `.gatorconfig.json` in your home directory, this should have the following format;

```
{
  "db_url": "your_postgres_url_here"
}
```

2. Run the program with `go-blog-aggregator [command] [arg]`

## Commands
For all commands, see the file `main.go` in the source. Here are some common commands;

- **register**: Registers a user, requires one argument which should be the username.
	- gator register kvothe
- **login**: Logs in a user, requires one argument which should be the username of a registered user.
	- gator login kvothe
- **addfeed**: Adds a feed to the database and subscribes the currently logged in user, requires two arguments the first being the name of the feed and the second being the URL.
	- gator addfeed "techcrunch" "https://techcrunch.com/feed/"
- **browse**: Returns posts from the currently logged-in users feed, optional argument "limit" (int) default to 2.
	- gator browse 10