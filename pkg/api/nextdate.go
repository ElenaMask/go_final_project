package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

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
			log.Println("error when parse date param:", err)
			return
		}
	} else {
		now = time.Now()
	}

	result, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("error when calculate next task date:", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid format for yearly repetition")
		}
		return handleYearlyRule(now, startDate)
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid format for daily repetition")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid number of days")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("days must be between 1 and 400")
		}
		return handleDailyRule(now, startDate, days)
	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid format for weekly repetition")
		}
		return handleWeeklyRule(now, startDate, parts[1])
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid format for monthly repetition")
		}
		daysPart := parts[1]
		monthsPart := ""
		if len(parts) == 3 {
			monthsPart = parts[2]
		}
		return handleMonthlyRule(now, startDate, daysPart, monthsPart)
	default:
		return "", errors.New("unsupported repetition format")
	}
}

func handleYearlyRule(now, startDate time.Time) (string, error) {
	date := startDate

	date = date.AddDate(1, 0, 0)
	for !afterNow(date, now) {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(DateFormat), nil
}

func handleDailyRule(now, startDate time.Time, days int) (string, error) {
	date := startDate

	date = date.AddDate(0, 0, days)
	for !afterNow(date, now) {
		date = date.AddDate(0, 0, days)
	}

	return date.Format(DateFormat), nil
}

func handleWeeklyRule(now, startDate time.Time, daysStr string) (string, error) {
	dayNumbers := strings.Split(daysStr, ",")
	if len(dayNumbers) == 0 {
		return "", errors.New("no days specified for weekly repetition")
	}

	var validDays [8]bool

	for _, dayStr := range dayNumbers {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", errors.New("invalid day number")
		}
		if day < 1 || day > 7 {
			return "", errors.New("day must be between 1 and 7")
		}
		validDays[day] = true
	}

	date := startDate
	date = date.AddDate(0, 0, 1)

	for {
		if afterNow(date, now) {
			weekday := int(date.Weekday())
			// Convert Go's weekday (Sunday=0) to our format (Monday=1, Sunday=7)
			if weekday == 0 {
				weekday = 7
			}

			if validDays[weekday] {
				return date.Format(DateFormat), nil
			}
		}

		date = date.AddDate(0, 0, 1)

		if date.Year() > now.Year()+10 {
			return "", errors.New("could not find next valid date")
		}
	}
}

func handleMonthlyRule(now, startDate time.Time, daysStr, monthsStr string) (string, error) {
	dayNumbers := strings.Split(daysStr, ",")
	if len(dayNumbers) == 0 {
		return "", errors.New("no days specified for monthly repetition")
	}

	var validDays [34]bool // 33 = -1, 32 = -2
	hasNegativeDays := false

	for _, dayStr := range dayNumbers {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", errors.New("invalid day number")
		}

		switch {
		case day >= 1 && day <= 31:
			validDays[day] = true
		case day == -1:
			hasNegativeDays = true
			validDays[33] = true
		case day == -2:
			hasNegativeDays = true
			validDays[32] = true
		default:
			return "", errors.New("day must be between 1 and 31, or -1, -2")
		}
	}

	var validMonths [13]bool
	monthsSpecified := monthsStr != ""

	if monthsSpecified {
		monthNumbers := strings.Split(monthsStr, ",")
		if len(monthNumbers) == 0 {
			return "", errors.New("no months specified for monthly repetition")
		}

		for _, monthStr := range monthNumbers {
			month, err := strconv.Atoi(strings.TrimSpace(monthStr))
			if err != nil {
				return "", errors.New("invalid month number")
			}
			if month < 1 || month > 12 {
				return "", errors.New("month must be between 1 and 12")
			}
			validMonths[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			validMonths[i] = true
		}
	}

	date := startDate
	date = date.AddDate(0, 0, 1)

	for {
		if afterNow(date, now) {
			day := date.Day()
			month := int(date.Month())

			if !validMonths[month] {
				date = date.AddDate(0, 0, 1)
				continue
			}

			if validDays[day] {
				return date.Format(DateFormat), nil
			}

			if hasNegativeDays {
				lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

				if day == lastDay && validDays[33] {
					return date.Format(DateFormat), nil
				}

				if day == lastDay-1 && validDays[32] {
					return date.Format(DateFormat), nil
				}
			}
		}

		date = date.AddDate(0, 0, 1)

		if date.Year() > now.Year()+10 {
			return "", errors.New("could not find next valid date")
		}
	}
}

func afterNow(date, now time.Time) bool {
	dateStr := date.Format(DateFormat)
	nowStr := now.Format(DateFormat)

	dateOnly, _ := time.Parse(DateFormat, dateStr)
	nowOnly, _ := time.Parse(DateFormat, nowStr)

	return dateOnly.After(nowOnly)
}
