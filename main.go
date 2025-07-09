package main

import (
	"fmt"
	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Run()
}
