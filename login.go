package main

import (
	"context"
	"fmt"
	"strings"
)

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
