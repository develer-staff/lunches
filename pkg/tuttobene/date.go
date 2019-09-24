package tuttobene

import (
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var testYear = -1

func setTestYear(year int) {
	testYear = year
}

func parseDate(content string) (bool, time.Time) {

	content = strings.ToLower(content)

	months := []string{
		"gennaio",
		"febbraio",
		"marzo",
		"aprile",
		"maggio",
		"giugno",
		"luglio",
		"agosto",
		"settembre",
		"ottobre",
		"novembre",
		"dicembre",
	}

	weekDays := []string{
		"dom",
		"lun",
		"mar",
		"mer",
		"gio",
		"ven",
		"sab",
	}

	args := strings.Split(content, " ")
	if len(args) != 3 {
		return false, time.Time{}
	}

	weekDay := -1
	for i, d := range weekDays {
		if strings.HasPrefix(args[0], d) {
			weekDay = i
			break
		}
	}
	if weekDay == -1 {
		return false, time.Time{}
	}

	cutset := ""
	for _, c := range args[1] {
		if !unicode.IsDigit(c) {
			cutset = cutset + string(c)
		}
	}
	dayString := strings.Trim(args[1], cutset)
	day, err := strconv.Atoi(dayString)
	if err != nil {
		return false, time.Time{}
	}

	month := -1
	for i, m := range months {
		if strings.HasPrefix(args[2], m) {
			month = i + 1
			break
		}
	}
	if month == -1 {
		return false, time.Time{}
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Println("LoadLocation error: ", err)
		return false, time.Time{}
	}

	var year int
	if testYear == -1 {
		year = time.Now().In(loc).Year()
	} else {
		year = testYear
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	if date.Weekday() != time.Weekday(weekDay) {
		log.Println("Weekday mismatch!")
		return false, time.Time{}
	}
	return true, date
}
