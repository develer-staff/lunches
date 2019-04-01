package tuttobene

import (
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
