package api

import (
	"net/http"
	"time"
)

// NextDateHandler handles the /api/nextdate endpoint
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Get parameters
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	// Parse now parameter or use current time
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

	// Call NextDate function
	result, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write result to response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
