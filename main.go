package main

import (
	// "fmt"
	"log"

	"github.com/alla-devyatova/go_final_project/pkg/db"
	"github.com/alla-devyatova/go_final_project/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("error in init db: %v", err)
	}
	defer db.CloseDB()

	// var tasks []*db.Task
	// tasks, _ = db.Tasks(5)
	// for _, task := range tasks {
	// 	fmt.Println("-", task)
	// }

	server.Run()
}
