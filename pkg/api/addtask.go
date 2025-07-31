package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ElenaMask/go_final_project/pkg/db"
)

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, Response{Error: "Некорректный формат JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, Response{Error: "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, Response{Error: err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, Response{Error: "Ошибка добавления задачи в базу данных"})
		return
	}

	writeJSON(w, Response{ID: fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	} else if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	return nil
}
