package tuttobene

import (
	"log"
	"strings"
	"time"
)

const (
	Unknonwn MenuRowType = iota
	Empty
	Primo
	Secondo
	Contorno
	Vegetariano
	Frutta
	Dolce
	Panino
)

type MenuRowType int

type MenuRow struct {
	Content         string
	Type            MenuRowType
	IsDailyProposal bool
}

type Menu struct {
	Rows []MenuRow
	Date time.Time
}

func (m *Menu) IsUpdated() bool {
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Println("LoadLocation error: ", err)
		return false
	}

	now := time.Now().In(loc)
	return (m.Date.Year() == now.Year()) && (m.Date.Month() == now.Month()) && (m.Date.Day() == now.Day())
}

func (m *Menu) String() string {
	menutype := Unknonwn

	out := "Data: *" + m.Date.Format("02/01/2006") + "*\n"
	for _, r := range m.Rows {
		if r.Type != menutype {
			out = out + "\n*" + strings.ToUpper(Titles[r.Type]) + "*\n"
			menutype = r.Type
		}
		if r.IsDailyProposal {
			out += "Proposta del giorno: "
		}

		out = out + r.Content + "\n"
	}
	return out
}

func (m *Menu) Add(mr *MenuRow) {

	//Check and remove duplicate dishes, keep only the last one added
	for i, r := range m.Rows {
		if r.Content == mr.Content {
			m.Rows = append(m.Rows[:i], m.Rows[i+1:]...)
			break
		}
	}

	m.Rows = append(m.Rows, *mr)
}
