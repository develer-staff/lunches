package tuttobene

import (
	"path/filepath"
	"reflect"
	"testing"

	_ "github.com/tealeg/xlsx"
)

func TestParseMenu(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Menu
		wantErr bool
	}{
		{
			"testmenu1",
			args{filepath.Join("test-fixtures", "testmenu1.xlsx")},
			&Menu{
				{"Rigatoni al ragù dell'aia", 2, false},
				{"Ravioli ricotta e spinaci con burro e salvia", 2, false},
				{"Lasagne con cavolo nero e porri", 2, false},
				{"Minestra di pane", 2, false},
				{"Paccheri con calamari e asparagi", 2, false},
				{"Pasta al ragù", 2, false},
				{"Pasta al pesto", 2, false},
				{"Pasta al pomodoro", 2, false},
				{"Lasagne cavolo nero e porri + macedonia", 2, true},
				{"Roastbeef con patate arrosto", 3, false},
				{"Polpette in umido con verdure", 3, false},
				{"Spezzatino di vitella con asparagi", 3, false},
				{"Baccalà alla livornese con fagioli", 3, false},
				{"Filetto di branzino gratinato con fagiolini", 3, false},
				{"Baccalà alla livornese con fagioli + macedonia", 3, true},
				{"Sformatini di riso con verdure al vapore", 4, false},
				{"Fantasia di verdure grigliate", 4, false},
				{"Macedonia di frutta fresca", 5, false},
				{"Macedonia di frutta fresca piccola", 5, false},
				{"Frutta a tocchi", 5, false},
				{"Diametro 12 mortadella", 6, false},
				{"Diametro 12 crudo pecorino e rucola", 6, false},
				{"Diametro 8 bresaola rucola e brie", 6, false},
				{"Diametro 8 vegetariano", 6, false},
				{"Tubo 15 tonno maionese e pomodoro", 6, false},
				{"Tubo 15 praga radicchi e grana", 6, false},
			},
			false,
		},
		{
			"testmenu2",
			args{filepath.Join("test-fixtures", "testmenu2.xlsx")},
			&Menu{
				{"Sedani alla Carloforte", 2, false},
				{"Strigoli con filangè di verdure e speck", 2, false},
				{"Orecchiette alle rape", 2, false},
				{"Zuppa di zucca con pane croccante", 2, false},
				{"Paccheri alla triglia", 2, false},
				{"Pasta al ragù", 2, false},
				{"Pasta al pesto", 2, false},
				{"Pasta al pomodoro", 2, false},
				{"Orecchiette alle rape  + macedonia", 2, true},
				{"Polpette in umido con purè", 3, false},
				{"Ossibuchi alla livornese con fagioli borlotti", 3, false},
				{"Filetto di maiale con panure a i 3 pepi e patate arrosto", 3, false},
				{"Orata all'isolana con spinaci", 3, false},
				{"Seppie con piselli", 3, false},
				{"Polpette in umido con purè + macedonia", 3, true},
				{"Insalata di spinacina, fagioli di soja, feta e mais", 4, false},
				{"Dadolata di verdure al forno", 4, false},
				{"Macedonia di frutta fresca", 5, false},
				{"Macedonia di frutta fresca piccola", 5, false},
				{"Frutta a tocchi", 5, false},
				{"Diametro 12 mortadella", 6, false},
				{"Diametro 12 crudo pecorino e rucola", 6, false},
				{"Diametro 8 bresaola rucola e brie", 6, false},
				{"Diametro 8 vegetariano", 6, false},
				{"Tubo 15 tonno maionese e pomodoro", 6, false},
				{"Tubo 15 praga radicchi e grana", 6, false},
			},
			false,
		},
		{
			"testmenu3",
			args{filepath.Join("test-fixtures", "testmenu3.xlsx")},
			&Menu{
				{"Penne con salsiccia e rape", 2, false},
				{"Pici cacio e pepe", 2, false},
				{"Crespelle alla fiorentina", 2, false},
				{"Minestrone", 2, false},
				{"Paccheri al polpo", 2, false},
				{"Pasta al ragù", 2, false},
				{"Pasta al pesto", 2, false},
				{"Pasta al pomodoro", 2, false},
				{"Penne con salsiccia e rape  + macedonia", 2, true},
				{"Pollo al curry con riso nero", 3, false},
				{"Hamburger con pomodori grigliati", 3, false},
				{"Bianchetto di vitellla con champignon", 3, false},
				{"Moscardini con piselli", 3, false},
				{"Spada alla griglia con belga", 3, false},
				{"Hamburger con pomodori grigliati + macedonia", 3, true},
				{"Insalata di  zucca gialla con pomodori e olive", 4, false},
				{"Fantasia di verdure al vapore", 4, false},
				{"Macedonia di frutta fresca", 5, false},
				{"Macedonia di frutta fresca piccola", 5, false},
				{"Frutta a tocchi", 5, false},
				{"Diametro 12 mortadella", 6, false},
				{"Diametro 12 crudo pecorino e rucola", 6, false},
				{"Diametro 8 bresaola rucola e brie", 6, false},
				{"Diametro 8 vegetariano", 6, false},
				{"Tubo 15 tonno maionese e pomodoro", 6, false},
				{"Tubo 15 praga radicchi e grana", 6, false},
			},
			false,
		},
		{
			"doesnotexist",
			args{"doesnotexist.xlsx"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMenuFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMenuFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMenuFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
