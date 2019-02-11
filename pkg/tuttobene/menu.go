package tuttobene

const (
	Unknonwn MenuRowType = iota
	Empty
	Primo
	Secondo
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

type Menu []MenuRow
