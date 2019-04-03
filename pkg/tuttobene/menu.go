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
		out = out + r.Content + "\n"
	}
	return out
}
