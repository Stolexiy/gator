package main

import (
	"context"
	"fmt"
)

func resetHandler(st *state, cmd command) error {
	err := st.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset the users table: %v", err)
	}
	fmt.Println("the users table was successfully cleaned")
	return nil
}
