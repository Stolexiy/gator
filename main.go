package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/stolexiy/gator/internal/config"
	"github.com/stolexiy/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	arg  []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if cmdHandler, exist := c.list[cmd.name]; exist {
		return cmdHandler(s, cmd)
	}

	return fmt.Errorf("%s command doesn't exist", cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if _, exist := c.list[name]; !exist {
		c.list[name] = f
	}
	return fmt.Errorf("%s command is already registered", name)
}

func main() {
	cmds := commands{
		list: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", handleAddfeed)
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", handleFollow)
	cmds.register("following", handleFollowing)

	args := os.Args

	if len(args) < 2 {
		log.Fatalln("not enough arguments provided")
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	dbQueries := database.New(db)

	st := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	err = cmds.run(st, command{
		name: args[1],
		arg:  args[2:],
	})
	if err != nil {
		log.Fatalln(err)
	}
}
