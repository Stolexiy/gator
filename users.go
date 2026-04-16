package main

import (
	"context"
	"fmt"
)

func usersHandler(st *state, cmd command) error {
	users, err := st.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	if len(users) == 0 {
		fmt.Println("no users")
	}

	for _, u := range users {
		if st.cfg.CurrentUserName == u.Name {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}

	return nil
}
