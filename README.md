# gator

A CLI RSS feed aggregator written in Go. Register users, follow feeds, aggregate posts in the background, and browse them from the terminal.

## Prerequisites

You'll need the following installed locally:

- **Go** (1.25+) — to build and install the CLI: https://go.dev/doc/install
- **PostgreSQL** (15+) — gator stores users, feeds, follows, and posts in Postgres: https://www.postgresql.org/download/

Make sure your Postgres server is running and you have a database you can connect to (e.g. `postgres://postgres:postgres@localhost:5432/gator`).

## Install

Install the `gator` binary with `go install`:

```
go install github.com/stolexiy/gator@latest
```

This drops the `gator` executable into `$(go env GOPATH)/bin`. Make sure that directory is on your `PATH`.

## Database setup

Run the migrations against your database with [goose](https://github.com/pressly/goose):

```
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

## Config file

gator reads its config from `~/.gatorconfig.json`. Create it with your database connection string:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

The `current_user_name` field is updated automatically when you `register` or `login`.

## Usage

Run commands as `gator <command> [args]`. A few of the things you can do:

- `gator register <name>` — create a new user and log in as them.
- `gator login <name>` — switch to an existing user.
- `gator users` — list all registered users.
- `gator addfeed <name> <url>` — add a new RSS feed and follow it.
- `gator feeds` — list every feed in the database.
- `gator follow <url>` — follow an existing feed.
- `gator unfollow <url>` — stop following a feed.
- `gator following` — show the feeds the current user follows.
- `gator agg <interval>` — start the background aggregator (e.g. `gator agg 1m` fetches each due feed once a minute). Leave it running in a terminal.
- `gator browse [limit]` — show the latest posts from the feeds you follow.
- `gator reset` — wipe all data (useful while developing).