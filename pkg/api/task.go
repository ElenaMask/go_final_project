package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ElenaMask/go_final_project/pkg/db"
)

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type APITask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func writeError(w http.ResponseWriter, message string) {
	writeJSON(w, Response{Error: message})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Некорректный формат JSON")
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка добавления задачи в базу данных")
		return
	}

	writeJSON(w, Response{ID: fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор")
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}
	apiTask := &APITask{
		ID:      strconv.FormatInt(t.ID, 10),
		Date:    t.Date,
		Title:   t.Title,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	}

	writeJSON(w, apiTask)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var apiTask APITask

	if err := json.NewDecoder(r.Body).Decode(&apiTask); err != nil {
		writeError(w, "Некорректный формат JSON")
		return
	}
	id, err := strconv.ParseInt(apiTask.ID, 10, 64)
	if err != nil {
		writeError(w, "Некорректный идентификатор задачи")
		return
	}

	task := db.Task{
		ID:      id,
		Date:    apiTask.Date,
		Title:   apiTask.Title,
		Comment: apiTask.Comment,
		Repeat:  apiTask.Repeat,
	}

	if task.ID == 0 {
		writeError(w, "Не указан идентификатор задачи")
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, "Задача не найдена")
		return
	}

	writeJSON(w, Response{})
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
