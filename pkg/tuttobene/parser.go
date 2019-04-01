package tuttobene

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/juju/errors"
	"github.com/tealeg/xlsx"
)

var Titles = map[MenuRowType]string{
	Primo:       "primi piatti",
	Secondo:     "secondi piatti",
	Contorno:    "contorni",
	Vegetariano: "piatti vegetariani",
	Frutta:      "frutta",
	Panino:      "i nostri panini espressi",
}

// ParseMenuReaderAt takes io.ReaderAt of an XLSX file and returns a populated
// menu struct.
func ParseMenuReaderAt(r io.ReaderAt, size int64) (*Menu, error) {
	f, err := xlsx.OpenReaderAt(r, size)
	if err != nil {
		return nil, errors.Annotate(err, "while opening readerAt")
	}

	if len(f.Sheet) == 0 {
		return nil, errors.New("no sheets in file")
	}

	// Menu is expected to be on the first sheet
	return ParseSheet(f.Sheets[0])
}

// ParseMenuBytes takes io.ReaderAt of an XLSX file and returns a populated
// menu struct.
func ParseMenuBytes(bs []byte) (*Menu, error) {
	f, err := xlsx.OpenBinary(bs)
	if err != nil {
		return nil, errors.Annotate(err, "while opening binary")
	}

	if len(f.Sheet) == 0 {
		return nil, errors.New("no sheets in file")
	}

	// Menu is expected to be on the first sheet
	return ParseSheet(f.Sheets[0])
}

// ParseMenuFile takes the path to an XLSX file and returns a populated
// menu struct.
func ParseMenuFile(path string) (*Menu, error) {
	f, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, errors.Annotatef(err, "while opening file %s", path)
	}

	if len(f.Sheet) == 0 {
		return nil, errors.New(fmt.Sprintf("no sheets in file %s", path))
	}

	// Menu is expected to be on the first sheet
	return ParseSheet(f.Sheets[0])
}

// ParseSheet takes an xlsx.Sheet and returns a populated menu struct.
func ParseSheet(s *xlsx.Sheet) (*Menu, error) {
	// attempt at having a sensible number of rows required in menu
	if len(s.Rows) < 12 {
		return nil, errors.New(fmt.Sprintf("not enough rows: %d", len(s.Rows)))
	}

	var rows = make([]string, 0)

	// Check tuttobene menu format (dishes in column 0 or 1)
	col := 0
	if len(s.Rows[0].Cells) >= 2 {
		sheetTitle := s.Rows[0].Cells[1].String()
		sheetTitle = strings.TrimSpace(sheetTitle)
		sheetTitle = strings.ToLower(sheetTitle)
		if sheetTitle == "tuttobene" {
			col = 1
		}
	}

	for _, r := range s.Rows {
		if len(r.Cells) >= col+1 {
			rows = append(rows, r.Cells[col].String())
		}
	}

	return ParseMenuRows(rows)
}

// ParseMenuRows takes a slice of strings and returns a populated menu struct.
func ParseMenuRows(rows []string) (*Menu, error) {
	var (
		currentType MenuRowType
		menuRows    Menu
	)

	for _, r := range rows {
		content, rowType, isTitle, isDailyProposal := parseRow(standardizeSpaces(r))

		if isTitle {
			currentType = rowType
			continue
		}

		// Skip first empty rows/check menu date
		if currentType == Unknonwn {
			isDate, date := parseDate(standardizeSpaces(r))
			if isDate {
				menuRows.Date = date
			}
			continue
		}

		// Check if this is the end of the menu
		if currentType == Panino && rowType == Empty {
			break
		}
		if rowType == Empty {
			continue
		}

		// Handle "Pasta al ragù, pesto o pomodoro (sono sempre disponibili)"
		if strings.HasSuffix(content, "(sono sempre disponibili)") {

			menuRows.Rows = append(menuRows.Rows, MenuRow{
				Content:         "Pasta al ragù",
				Type:            currentType,
				IsDailyProposal: false,
			})

			menuRows.Rows = append(menuRows.Rows, MenuRow{
				Content:         "Pasta al pesto",
				Type:            currentType,
				IsDailyProposal: false,
			})

			menuRows.Rows = append(menuRows.Rows, MenuRow{
				Content:         "Pasta al pomodoro",
				Type:            currentType,
				IsDailyProposal: false,
			})

			continue
		}

		menuRows.Rows = append(menuRows.Rows, MenuRow{
			Content:         strings.TrimSpace(content),
			Type:            currentType,
			IsDailyProposal: isDailyProposal,
		})
	}

	if (menuRows.Date == time.Time{}) {
		loc, err := time.LoadLocation("Europe/Rome")
		if err != nil {
			log.Println("LoadLocation error: ", err)
			return nil, err
		}
		menuRows.Date = time.Now().In(loc)
	}

	return &menuRows, nil
}

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

func parseRow(content string) (string, MenuRowType, bool, bool) {
	var isDailyProposal bool

	if content == "" {
		return content, Empty, false, false
	}

	isTitle, titleType := parseTitle(content)
	if isTitle {
		return content, titleType, isTitle, isDailyProposal
	}

	if strings.HasPrefix(content, "Proposta del giorno: ") {
		content = strings.TrimPrefix(content, "Proposta del giorno: ")
		isDailyProposal = true
	}

	return content, Unknonwn, isTitle, isDailyProposal
}

func parseTitle(content string) (bool, MenuRowType) {
	content = cleanTitleString(standardizeSpaces(content))
	for k, title := range Titles {
		if strings.EqualFold(title, content) || strings.EqualFold(strings.Replace(title, " ", "", -1), content) {
			return true, k
		}
	}

	return false, Unknonwn
}

var titleDirt = []string{
	"…",
	".",
	",",
	";",
}

func cleanTitleString(s string) string {
	for _, p := range titleDirt {
		s = strings.Replace(s, p, "", -1)
	}

	return strings.ToLower(s)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
