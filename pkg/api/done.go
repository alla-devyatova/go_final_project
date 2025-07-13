package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка, задача не найдена"}, http.StatusNotFound)
		return
	}
	if task.Repeat != "" {
		new_date, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "ошибка, не удалось определить дату задачи"}, http.StatusNotFound)
			return
		}
		var t db.Task
		t.Comment = task.Comment
		t.Date = new_date
		t.Repeat = task.Repeat
		t.Title = task.Title
		t.ID = id

		err = db.UpdateTask(&t)
		if err != nil {
			writeJson(w, map[string]string{"error": "ошибка, не удалось обновить задачу"}, http.StatusInternalServerError)
			return
		}
	} else {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "ошибка, не удалось удалить задачу"}, http.StatusInternalServerError)
			return
		}
	}
	writeJson(w, make(map[string]interface{}), http.StatusOK)
}
