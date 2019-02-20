package tuttobene

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/sahilm/fuzzy"
	"github.com/tealeg/xlsx"
	"io"
	"strings"
)

var Titles = map[MenuRowType]string{
	Primo:       "primi piatti",
	Secondo:     "secondi piatti",
	Contorno:    "contorni",
	Vegetariano: "piatti vegetariani",
	Frutta:      "frutta",
	Panino:      "panini espressi",
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
	for _, r := range s.Rows {
		if len(r.Cells) < 2 {
			continue
		}

		rows = append(rows, r.Cells[1].String())
	}

	return ParseMenuRows(rows)
}

// ParseMenuRows takes a slice of strings and returns a populated menu struct.
func ParseMenuRows(rows []string) (*Menu, error) {
	var (
		currentType          MenuRowType
		menuRows             Menu
		menuTitlesRowIndexes = make(map[int]MenuRowType)
	)

	menuTitlesRowIndexes, err  := getMenuTitles(rows)
	if err != nil {
		return nil, errors.Annotatef(err, "while getting menu titles")
	}

	for i, r := range rows {
		rowType, isTitle := menuTitlesRowIndexes[i]
		if isTitle {
			currentType = rowType
			continue
		}

		// Skip first empty rows
		if currentType == Unknonwn {
			continue
		}

		content := standardizeSpaces(r)
		if content == "" {
			if currentType == Panino {
				return &menuRows, nil
			}
			continue
		}

		// Handle "Pasta al ragù, pesto o pomodoro (sono sempre disponibili)"
		if strings.HasSuffix(content, "(sono sempre disponibili)") {

			menuRows = append(menuRows, MenuRow{
				Content:         "Pasta al ragù",
				Type:            currentType,
				IsDailyProposal: false,
			})

			menuRows = append(menuRows, MenuRow{
				Content:         "Pasta al pesto",
				Type:            currentType,
				IsDailyProposal: false,
			})

			menuRows = append(menuRows, MenuRow{
				Content:         "Pasta al pomodoro",
				Type:            currentType,
				IsDailyProposal: false,
			})

			continue
		}

		content, isDailyProposal := parseRow(content)
		
		menuRows = append(menuRows, MenuRow{
			Content:         strings.TrimSpace(content),
			Type:            currentType,
			IsDailyProposal: isDailyProposal,
		})
	}

	return &menuRows, nil
}

func parseRow(content string) (string, bool) {
	if strings.HasPrefix(content, "Proposta del giorno: ") {
		return strings.TrimPrefix(content, "Proposta del giorno: "), true
	}

	return content, false
}

// getMenuTitles returns a map of the row index for each of the sections found in the menu.
// Fuzzy matching is used to find the titles and some basic validation is done:
// - order: the titles are expected to be in the order in which the relative const enumeration is declared (see menu.go)
// - duplicates: if a duplicate title is found, an error is returned
//
// Note: it is not expected for all secitons to always be present i.e. if a section is missing, no error is thrown.
func getMenuTitles(rows []string) (map[int]MenuRowType, error) {
	var (
		menuTitlesRowIndexes = make(map[int]MenuRowType)
		lastTitleType = Unknonwn
		currentIndex int
	)

	for t, title := range Titles {
		results := fuzzy.Find(title, rows)
		if len(results) == 0 {
			continue
		}

		if t < lastTitleType {
			return nil, errors.New(fmt.Sprintf("Unexpected title order (Found: %s after last: %s)", t, lastTitleType))
		}

		currentIndex = results[0].Index
		if _, found := menuTitlesRowIndexes[currentIndex]; found {
			return nil, errors.New(fmt.Sprintf("Unexptected title duplicate: %s", title))
		}


		// First match is always the title of a section (menu items may contain the same text)
		menuTitlesRowIndexes[results[0].Index] = t
	}

	return menuTitlesRowIndexes, nil
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
