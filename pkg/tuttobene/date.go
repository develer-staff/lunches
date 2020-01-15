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

func findWeek(str string) int {

	weekDays := []string{
		"domenica",
		"luned",
		"marted",
		"mercoled",
		"gioved",
		"venerd",
		"sabato",
	}

	weekDay := -1
	for i, d := range weekDays {
		if strings.Contains(str, d) {
			weekDay = i
			break
		}
	}

	return weekDay
}

func findMonth(str string) int {
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

	month := -1
	for i, m := range months {
		if strings.Contains(str, m) {
			month = i + 1
			break
		}
	}

	return month
}

func findDay(str string) int {
	day := -1

	start := -1
	stop := -1

	for i, c := range str {
		if unicode.IsDigit(c) {
			if start == -1 {
				start = i
			}
		} else {
			if stop == -1 && start != -1 {
				stop = i
				break
			}
		}
	}

	if start != -1 && stop != -1 {
		conv, err := strconv.Atoi(str[start:stop])
		if err == nil && conv >= 1 && conv <= 31 {
			day = conv
		}
	}

	return day
}

func parseDate(content string) (bool, time.Time) {

	content = strings.ToLower(content)

	weekDay := findWeek(content)
	month := findMonth(content)
	day := findDay(content)

	log.Println("DATA: ", weekDay, month, day)
	if weekDay == -1 || month == -1 || day == -1 {
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
