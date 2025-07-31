package api

import (
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	var now time.Time
	if nowParam != "" {
		var err error
		now, err = time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Invalid now parameter format", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}

	result, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
