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
		id := r.URL.Query().Get("id")
		task, err := db.GetTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "задача не найдена"})
			return
		}
		writeJson(w, task)
	case http.MethodPut:
		var task db.Task
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// десериализуем JSON
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "ошибка"})
			return
		}

		if checkDate(&task) != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "ошибка"})
			return
		}

		err = db.UpdateTask(&task)
		if err != nil {
			writeJson(w, map[string]string{"error": "задача не найдена"})
			return
		}
		writeJson(w, make(map[string]interface{}))
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		_, err := db.GetTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "задача не найдена"})
			return
		}
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "ошибка"})
			return
		}
		writeJson(w, make(map[string]interface{}))
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var task db.Task
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// десериализуем JSON
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "ошибка"})
			return
		}

		if checkDate(&task) != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "ошибка"})
			return
		}

		id, err := db.AddTask(&task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJson(w, map[string]int64{"id": id})
	}
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
