package tuttobene

import (
	"fmt"
	"github.com/juju/errors"
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
		currentType MenuRowType
		menuRows    Menu
	)

	for _, r := range rows {
		content, rowType, isTitle, isDailyProposal := parseRow(standardizeSpaces(r))

		if isTitle {
			currentType = rowType
			continue
		}

		// Skip first empty rows
		if currentType == Unknonwn {
			continue
		}

		// Check if this is the end of the menu
		if currentType == Panino && rowType == Empty {
			return &menuRows, nil
		}
		if rowType == Empty {
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

		menuRows = append(menuRows, MenuRow{
			Content:         strings.TrimSpace(content),
			Type:            currentType,
			IsDailyProposal: isDailyProposal,
		})
	}

	return &menuRows, nil
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
