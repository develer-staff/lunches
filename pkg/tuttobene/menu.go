package tuttobene

const (
	Unknonwn MenuRowType = iota
	Primo
	Secondo
	Contorno
	Vegetariano
	Frutta
	Panino
)

//go:generate stringer -type=MenuRowType

type MenuRowType int

type MenuRow struct {
	Content         string
	Type            MenuRowType
	IsDailyProposal bool
}

type Menu []MenuRow
