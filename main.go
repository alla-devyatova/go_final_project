package main

import (
	"go1f/pkg/db"
	"go1f/pkg/server"
	"log"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("error in init db: %v", err)
	}

	server.Run()
	defer db.CloseDB()
}
