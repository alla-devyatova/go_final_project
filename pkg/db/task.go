package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))

	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT * FROM scheduler ORDER BY date")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	var result []*Task

	now := time.Now()

	for rows.Next() {
		if count >= limit {
			break
		}

		var t Task

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		d, err := time.Parse("20060102", t.Date)
		if err != nil {
			return nil, err
		}

		if BeginningOfDay(now).After(d) {
			continue
		}

		count = count + 1
		result = append(result, &t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return result, nil
	} else {
		return []*Task{}, nil
	}
}

func BeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetTask(id string) (*Task, error) {
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var t Task

	row := DB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id_int))

	err = row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`

	res, err := DB.Exec(query, sql.Named("id", task.ID), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`

	_, err := DB.Exec(query, sql.Named("id", id))

	if err != nil {
		return err
	}
	return nil
}
