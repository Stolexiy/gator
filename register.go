package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stolexiy/gator/internal/database"
)

func registerHandler(st *state, cmd command) error {
	if len(cmd.arg) == 0 {
		return fmt.Errorf("missed name argument")
	}

	now := time.Now()
	u, err := st.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.arg[0],
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	st.cfg.SetUser(u.Name)

	fmt.Printf("new user was created: %v", u)

	return nil
}
