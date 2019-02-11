package tuttobene

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/tealeg/xlsx"
	"strings"
)

func ParseMenuBytes(bs []byte) (*Menu, error) {
	f, err := xlsx.OpenBinary(bs)
	if err != nil {
		return nil, errors.Annotate(err, "while opening binary")
	}

	if len(f.Sheet) == 0 {
		return nil, errors.New("no sheets in file")
	}

	// Menu is expected to be on the first sheet
	return parseSheet(f.Sheets[0])
}

func ParseMenu(path string) (*Menu, error) {
	f, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, errors.Annotatef(err, "while opening file %s", path)
	}

	if len(f.Sheet) == 0 {
		return nil, errors.New(fmt.Sprintf("no sheets in file %s", path))
	}

	// Menu is expected to be on the first sheet
	return parseSheet(f.Sheets[0])
}

func parseSheet(s *xlsx.Sheet) (*Menu, error) {
	// attempt at having a sensible number of rows required in menu
	if len(s.Rows) < 12 {
		return nil, errors.New(fmt.Sprintf("not enough rows: %d", len(s.Rows)))
	}

	var (
		currentType MenuRowType
		menuRows    Menu
	)

	for _, r := range s.Rows {
		if len(r.Cells) < 3 {
			continue
		}

		content, rowType, isTitle, isDailyProposal:= parseRow(r.Cells[1].String())

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

		m := MenuRow{
			Content:         strings.TrimSpace(content),
			Type:            currentType,
			IsDailyProposal: isDailyProposal,
		}

		menuRows = append(menuRows, m)
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
		return content, titleType, true, isDailyProposal
	}

	if strings.HasPrefix(content, "Proposta del giorno: ") {
		content = strings.TrimPrefix(content, "Proposta del giorno: ")
		isDailyProposal = true
	}

	return content, Unknonwn, false, isDailyProposal
}

func parseTitle(content string) (bool, MenuRowType) {
	switch content {
	case "Primi piatti":
		return true, Primo
	case "Secondi piatti":
		return true, Secondo
	case "Piatti vegetariani":
		return true, Vegetariano
	case "Frutta":
		return true, Frutta
	case "I NOSTRI  PANINI  ESPRESSI…":
		return true, Panino
	}

	return false, Empty
}
