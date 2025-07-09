package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	now, err := time.Parse("20060102", req.FormValue("now"))
	if err != nil {
		now = time.Now()
	}
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	next_date, err := NextDate(now, date, repeat)

	if err != nil {
		res.Write([]byte(err.Error()))
	} else {
		res.Write([]byte(next_date))
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("неверный формат условия повторения")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	first_char := repeat[:1]
	day_0 := time.Date(1, 1, 1, 0, 0, 0, 0, now.Location())

	switch first_char {
	case "":
		return "", errors.New("неверный формат условия повторения")
	case "d":
		splitted := strings.Split(repeat, " ")

		if len(splitted) != 2 {
			return "", errors.New("неверный формат условия повторения")
		}

		interval, err := strconv.Atoi(splitted[1])
		if err != nil {
			return "", errors.New("неверный формат условия повторения")
		}

		if interval > 400 {
			return "", errors.New("неверный формат условия повторения")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}

		return date.Format("20060102"), nil
	case "y":
		if repeat != "y" {
			return "", errors.New("неверный формат условия повторения")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format("20060102"), nil
	case "w":
		splitted := strings.Split(repeat, " ")

		if len(splitted) != 2 {
			return "", errors.New("неверный формат условия повторения")
		}

		date_start := day_0
		if now.After(date) {
			date_start = now
		} else {
			date_start = date
		}

		splitted2 := strings.Split(splitted[1], ",")
		weekday := int(date_start.Weekday())
		min := 7
		for _, i := range splitted2 {
			number, err := strconv.Atoi(i)
			if err != nil {
				return "", errors.New("неверный формат условия повторения")
			}
			if (number < 1) || (number > 7) {
				return "", errors.New("неверный формат условия повторения")
			}

			var interval int
			if number > weekday {
				interval = number - weekday
			} else {
				interval = number + 7 - weekday
			}

			if interval < min {
				min = interval
			}
		}

		date = date_start.AddDate(0, 0, min)
		return date.Format("20060102"), nil
	case "m":
		splitted := strings.Split(repeat, " ")

		date_start := day_0
		if now.After(date) {
			date_start = now
		} else {
			date_start = date
		}

		if len(splitted) == 2 {
			new_date := day_0

			splitted2 := strings.Split(splitted[1], ",")
			for _, i := range splitted2 {
				day, err := strconv.Atoi(i)
				if err != nil {
					return "", errors.New("неверный формат условия повторения")
				}
				if !(((day >= 1) && (day <= 31)) || ((day >= -2) && (day <= -1))) {
					return "", errors.New("неверный формат условия повторения")
				}

				new_date_2 := day_0
				if day > 0 {
					new_date_2 = BeginningOfCurrentMonth(date_start).AddDate(0, 0, day-1)

					if (new_date_2.Day() != day) || date_start.After(new_date_2) || (date_start == new_date_2) {
						new_date_2 = BeginningOfNextMonth(date_start).AddDate(0, 0, day-1)
					}
				} else {
					new_date_2 = EndOfCurrentMonth(date_start).AddDate(0, 0, day+1)

					if date_start.After(new_date_2) || (date_start == new_date_2) {
						new_date_2 = new_date_2.AddDate(0, 1, 0)
					}
				}

				if new_date == day_0 {
					new_date = new_date_2
				} else if new_date.After(new_date_2) {
					new_date = new_date_2
				}
			}
			return new_date.Format("20060102"), nil
		} else if len(splitted) == 3 {
			new_date := day_0

			splitted3 := strings.Split(splitted[2], ",")
			for _, i := range splitted3 {
				month, err := strconv.Atoi(i)
				if err != nil {
					return "", errors.New("неверный формат условия повторения")
				}
				if (month < 1) || (month > 12) {
					return "", errors.New("неверный формат условия повторения")
				}

				splitted2 := strings.Split(splitted[1], ",")
				for _, i := range splitted2 {
					day, err := strconv.Atoi(i)
					if err != nil {
						return "", errors.New("неверный формат условия повторения")
					}
					if !(((day >= 1) && (day <= 31)) || ((day >= -2) && (day <= -1))) {
						return "", errors.New("неверный формат условия повторения")
					}

					new_date_2 := day_0
					if day > 0 {
						if !((month == 2) && (day > 28)) {
							new_date_2 = DayInCurrentYear(date_start, month, day)
							if (new_date_2 != day_0) && (date_start.After(new_date_2) || (date_start == new_date_2)) {
								new_date_2 = DayInNextYear(date_start, month, day)
							}
						} else {
							new_date_2 := DayOfFebruary(date_start, month, day)
							if (new_date_2 != day_0) && (date_start.After(new_date_2) || (date_start == new_date_2)) {
								new_date_2 = new_date_2.AddDate(4, 0, 0)
							}
						}
					} else {
						new_date_2 = EndOfMonthInCurrentYear(date_start, month).AddDate(0, 0, day+1)
						if date_start.After(new_date_2) || (date_start == new_date_2) {
							new_date_2 = new_date_2.AddDate(1, 0, 0)
						}
					}

					if new_date_2 != day_0 {
						if new_date == day_0 {
							new_date = new_date_2
						} else if new_date.After(new_date_2) {
							new_date = new_date_2
						}
					}
				}
			}
			return new_date.Format("20060102"), nil
		}
	default:
		return "", errors.New("неверный формат условия повторения")
	}

	return "", nil
}

func BeginningOfCurrentMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func BeginningOfNextMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, 0)
}

func EndOfCurrentMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, 0).AddDate(0, 0, -1)
}

func DayInCurrentYear(t time.Time, month int, day int) time.Time {
	new_date := time.Date(t.Year(), time.Month(month), day, 0, 0, 0, 0, t.Location())
	if new_date.Day() == day {
		return new_date
	} else {
		return time.Date(1, 1, 1, 0, 0, 0, 0, t.Location())
	}
}

func DayInNextYear(t time.Time, month int, day int) time.Time {
	new_date := time.Date(t.Year()+1, time.Month(month), day, 0, 0, 0, 0, t.Location())
	if new_date.Day() == day {
		return new_date
	} else {
		return time.Date(1, 1, 1, 0, 0, 0, 0, t.Location())
	}
}

func DayOfFebruary(t time.Time, month int, day int) time.Time {
	if day > 29 {
		return time.Date(1, 1, 1, 0, 0, 0, 0, t.Location())
	}

	if t.Year()%4 == 0 {
		return time.Date(t.Year(), time.Month(month), day, 0, 0, 0, 0, t.Location())
	}

	return time.Date(t.Year()+4-t.Year()%4, time.Month(month), day, 0, 0, 0, 0, t.Location())
}

func EndOfMonthInCurrentYear(t time.Time, month int) time.Time {
	return time.Date(t.Year(), time.Month(month), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, 0).AddDate(0, 0, -1)
}
