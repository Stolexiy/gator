package main

import (
	"fmt"
	"log"

	"github.com/stolexiy/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	err = cfg.SetUser("oleksiy")
	if err != nil {
		log.Fatalf("Failed to set user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	fmt.Println(cfg)
}
