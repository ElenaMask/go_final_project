package db

type Task struct {
	ID      int64  `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
