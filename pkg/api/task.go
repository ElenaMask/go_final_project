package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ElenaMask/go_final_project/pkg/db"
)

type APITask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "Некорректный формат JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка добавления задачи в базу данных", http.StatusInternalServerError)
		return
	}

	writeJSON(w, Response{ID: fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
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
		writeError(w, "Некорректный формат JSON", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(apiTask.ID, 10, 64)
	if err != nil {
		writeError(w, "Некорректный идентификатор задачи", http.StatusBadRequest)
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
		writeError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		log.Println("error on task update in database:", err)
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

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, fmt.Sprintf("Ошибка удаления задачи: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, fmt.Sprintf("Ошибка расчета следующей даты: %v", err), http.StatusInternalServerError)
			return
		}
		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeError(w, fmt.Sprintf("Ошибка обновления даты задачи: %v", err), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, Response{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, fmt.Sprintf("Ошибка удаления задачи: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, Response{})
}
