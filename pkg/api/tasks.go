package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ElenaMask/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*APITask `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	if searchQuery != "" {
		parsedTime, dateErr := time.Parse("02.01.2006", searchQuery)
		if dateErr == nil {
			dateFormatted := parsedTime.Format("20060102")
			tasks, err = db.GetTasksByDate(dateFormatted, 50)
		} else {
			tasks, err = db.SearchTasks(searchQuery, 50)
		}
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeError(w, "Ошибка получения задач из базы данных", http.StatusInternalServerError)
		return
	}

	apiTasks := make([]*APITask, len(tasks))
	for i, t := range tasks {
		apiTasks[i] = &APITask{
			ID:      strconv.FormatInt(t.ID, 10),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		}
	}

	writeJSON(w, TasksResp{
		Tasks: apiTasks,
	})
}
