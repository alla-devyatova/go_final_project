package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"go1f/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		findTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "ошибка, не заполнен заголовок задачи"})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if checkDate(&task) != nil {
		writeJson(w, map[string]string{"error": "ошибка, не удалось определить дату задачи"})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]int64{"id": id})
}

func findTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка, задача не найдена"})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем JSON
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "ошибка, не заполнен заголовок задачи"})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if checkDate(&task) != nil {
		writeJson(w, map[string]string{"error": "ошибка, не удалось определить дату задачи"})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := task.ID
	_, err = db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка, задача не найдена"})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка, не удалось обновить задачу"})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, make(map[string]interface{}))
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка, задача не найдена"})
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка, не удалось удалить задачу"})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, make(map[string]interface{}))
}

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	next := ""
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if db.BeginningOfDay(now).After(t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format("20060102")
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}

	return nil
}
