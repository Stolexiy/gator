package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) < 1 {
		return fmt.Errorf("username is required")
	}

	u, err := s.db.GetUser(context.Background(), cmd.arg[0])
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("user %s not exists", cmd.arg[0])
		}
		return err
	}

	err = s.cfg.SetUser(u.Name)
	if err != nil {
		return err
	}

	fmt.Println("current user has been set")
	return nil
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
