package api

import (
	"go1f/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка"})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
