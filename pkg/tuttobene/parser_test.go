package tuttobene

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	_ "github.com/tealeg/xlsx"
)

func TestParseMenu(t *testing.T) {
	type args struct {
		path string
	}
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		args    args
		year    int
		want    *Menu
		wantErr bool
	}{
		{
			"testmenu1",
			args{filepath.Join("test-fixtures", "testmenu1.xlsx")},
			2018,
			&Menu{
				[]MenuRow{
					{"Rigatoni al ragù dell'aia", Primo, false},
					{"Ravioli ricotta e spinaci con burro e salvia", Primo, false},
					{"Lasagne con cavolo nero e porri", Primo, false},
					{"Minestra di pane", Primo, false},
					{"Paccheri con calamari e asparagi", Primo, false},
					{"Pasta al ragù", Primo, false},
					{"Pasta al pesto", Primo, false},
					{"Pasta al pomodoro", Primo, false},
					{"Lasagne cavolo nero e porri + macedonia", Primo, true},
					{"Roastbeef con patate arrosto", Secondo, false},
					{"Polpette in umido con verdure", Secondo, false},
					{"Spezzatino di vitella con asparagi", Secondo, false},
					{"Baccalà alla livornese con fagioli", Secondo, false},
					{"Filetto di branzino gratinato con fagiolini", Secondo, false},
					{"Baccalà alla livornese con fagioli + macedonia", Secondo, true},
					{"Sformatini di riso con verdure al vapore", Vegetariano, false},
					{"Fantasia di verdure grigliate", Vegetariano, false},
					{"Macedonia di frutta fresca", Frutta, false},
					{"Macedonia di frutta fresca piccola", Frutta, false},
					{"Frutta a tocchi", Frutta, false},
					{"Diametro 12 mortadella", Panino, false},
					{"Diametro 12 crudo pecorino e rucola", Panino, false},
					{"Diametro 8 bresaola rucola e brie", Panino, false},
					{"Diametro 8 vegetariano", Panino, false},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false},
					{"Tubo 15 praga radicchi e grana", Panino, false},
				},
				time.Date(2018, 12, 10, 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			"testmenu2",
			args{filepath.Join("test-fixtures", "testmenu2.xlsx")},
			2019,
			&Menu{
				[]MenuRow{
					{"Sedani alla Carloforte", Primo, false},
					{"Strigoli con filangè di verdure e speck", Primo, false},
					{"Orecchiette alle rape", Primo, false},
					{"Zuppa di zucca con pane croccante", Primo, false},
					{"Paccheri alla triglia", Primo, false},
					{"Pasta al ragù", Primo, false},
					{"Pasta al pesto", Primo, false},
					{"Pasta al pomodoro", Primo, false},
					{"Orecchiette alle rape + macedonia", Primo, true},
					{"Polpette in umido con purè", Secondo, false},
					{"Ossibuchi alla livornese con fagioli borlotti", Secondo, false},
					{"Filetto di maiale con panure a i 3 pepi e patate arrosto", Secondo, false},
					{"Orata all'isolana con spinaci", Secondo, false},
					{"Seppie con piselli", Secondo, false},
					{"Polpette in umido con purè + macedonia", Secondo, true},
					{"Insalata di spinacina, fagioli di soja, feta e mais", Vegetariano, false},
					{"Dadolata di verdure al forno", Vegetariano, false},
					{"Macedonia di frutta fresca", Frutta, false},
					{"Macedonia di frutta fresca piccola", Frutta, false},
					{"Frutta a tocchi", Frutta, false},
					{"Diametro 12 mortadella", Panino, false},
					{"Diametro 12 crudo pecorino e rucola", Panino, false},
					{"Diametro 8 bresaola rucola e brie", Panino, false},
					{"Diametro 8 vegetariano", Panino, false},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false},
					{"Tubo 15 praga radicchi e grana", Panino, false},
				},
				time.Date(2019, 2, 12, 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			"testmenu3",
			args{filepath.Join("test-fixtures", "testmenu3.xlsx")},
			2019,
			&Menu{
				[]MenuRow{
					{"Penne con salsiccia e rape", Primo, false},
					{"Pici cacio e pepe", Primo, false},
					{"Crespelle alla fiorentina", Primo, false},
					{"Minestrone", Primo, false},
					{"Paccheri al polpo", Primo, false},
					{"Pasta al ragù", Primo, false},
					{"Pasta al pesto", Primo, false},
					{"Pasta al pomodoro", Primo, false},
					{"Penne con salsiccia e rape + macedonia", Primo, true},
					{"Pollo al curry con riso nero", Secondo, false},
					{"Hamburger con pomodori grigliati", Secondo, false},
					{"Bianchetto di vitellla con champignon", Secondo, false},
					{"Moscardini con piselli", Secondo, false},
					{"Spada alla griglia con belga", Secondo, false},
					{"Hamburger con pomodori grigliati + macedonia", Secondo, true},
					{"Insalata di zucca gialla con pomodori e olive", Vegetariano, false},
					{"Fantasia di verdure al vapore", Vegetariano, false},
					{"Macedonia di frutta fresca", Frutta, false},
					{"Macedonia di frutta fresca piccola", Frutta, false},
					{"Frutta a tocchi", Frutta, false},
					{"Diametro 12 mortadella", Panino, false},
					{"Diametro 12 crudo pecorino e rucola", Panino, false},
					{"Diametro 8 bresaola rucola e brie", Panino, false},
					{"Diametro 8 vegetariano", Panino, false},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false},
					{"Tubo 15 praga radicchi e grana", Panino, false},
				},
				time.Date(2019, 2, 13, 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			"Test with Contorni",
			args{filepath.Join("test-fixtures", "testmenuv2.xlsx")},
			2019,
			&Menu{
				[]MenuRow{
					{"Penne con salsiccia e rape", Primo, false},
					{"Pici cacio e pepe", Primo, false},
					{"Crespelle alla fiorentina", Primo, false},
					{"Minestrone", Primo, false},
					{"Paccheri al polpo", Primo, false},
					{"Pasta olio", Primo, false},
					{"Pasta al ragù", Primo, false},
					{"Riso olio", Primo, false},
					{"Pasta al pomodoro", Primo, false},
					{"Pollo al curry", Secondo, false},
					{"Hamburger", Secondo, false},
					{"Bianchetto di vitellla", Secondo, false},
					{"Moscardini con piselli", Secondo, false},
					{"Spada alla griglia", Secondo, false},
					{"Peperoni alla griglia", Contorno, false},
					{"Melanzane alla griglia", Contorno, false},
					{"Belga alla griglia", Contorno, false},
					{"Radicchio alla griglia", Contorno, false},
					{"Broccoli al vapore", Contorno, false},
					{"Cavolfiore al vapore", Contorno, false},
					{"Carote al vapore", Contorno, false},
					{"Fagiolini al vapore", Contorno, false},
					{"Dadolata di verdure al forno", Contorno, false},
					{"Pomodori", Contorno, false},
					{"Insalata", Contorno, false},
					{"Patate arrosto", Contorno, false},
					{"Spinaci saltati", Contorno, false},
					{"Pomodori grigliati", Contorno, false},
					{"Insalata di zucca gialla con pomodori e olive", Vegetariano, false},
					{"Fantasia di verdure al vapore", Vegetariano, false},
					{"Mozzarelle", Vegetariano, false},
					{"Macedonia di frutta fresca", Frutta, false},
					{"Macedonia di frutta fresca piccola", Frutta, false},
					{"Frutta a tocchi", Frutta, false},
				},
				time.Date(2019, 2, 13, 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			"Test date",
			args{filepath.Join("test-fixtures", "testmenu4.xlsx")},
			2019,
			&Menu{
				[]MenuRow{
					{"Penne all'amatriciana", Primo, false},
					{"Sedani salsiccia e olive", Primo, false},
					{"Paccheri zucchine e speck", Primo, false},
					{"Farro alla sorrentina (freddo)", Primo, false},
					{"Spaghetti allo scoglio", Primo, false},
					{"Pasta olio", Primo, false},
					{"Pasta al ragù", Primo, false},
					{"Pasta al pomodoro", Primo, false},
					{"Riso olio", Primo, false},
					{"Spiedini di carne", Secondo, false},
					{"Roastbeef", Secondo, false},
					{"Pollo ripieno", Secondo, false},
					{"Tagliata di tonno", Secondo, false},
					{"Salmone al vapore", Secondo, false},
					{"Tonno sott'olio", Secondo, false},
					{"Bresaola", Secondo, false},
					{"Prociutto crudo", Secondo, false},
					{"Peperoni alla griglia", Contorno, false},
					{"Melanzane alla griglia", Contorno, false},
					{"Belga alla griglia", Contorno, false},
					{"Finocchi alla griglia", Contorno, false},
					{"Radicchio alla griglia", Contorno, false},
					{"Broccoli al vapore", Contorno, false},
					{"Cavolfiore al vapore", Contorno, false},
					{"Carote al vapore", Contorno, false},
					{"Fagiolini al vapore", Contorno, false},
					{"Pomodori", Contorno, false},
					{"Insalata", Contorno, false},
					{"Patate arrosto", Contorno, false},
					{"Piselli", Contorno, false},
					{"Spinaci saltati", Contorno, false},
					{"Taccole al pomodoro", Contorno, false},
					{"Primosale con insalata mista", Vegetariano, false},
					{"Dadolata di verdure al forno", Vegetariano, false},
					{"Mozzarelle", Vegetariano, false},
					{"Macedonia di frutta fresca", Frutta, false},
					{"Macedonia di frutta fresca piccola", Frutta, false},
					{"Frutta a tocchi", Frutta, false},
				},
				time.Date(2019, 4, 1, 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			"doesnotexist",
			args{"doesnotexist.xlsx"},
			-1,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTestYear(tt.year)
			got, err := ParseMenuFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMenuFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMenuFile() = %v, want %v", got, tt.want)

				for i, item := range got.Rows {
					wanted := (*tt.want).Rows[i]
					if item.Type != wanted.Type {
						t.Errorf("ParseMenuFile() menu[%d] (%s) has wrong Type: got %v, want %v", i, item.Content, item.Type, wanted.Type)
					}

					if item.Content != wanted.Content {
						t.Errorf("ParseMenuFile() %d has wrong Content: got %v, want %v", i, item.Content, wanted.Content)
					}

					if item.IsDailyProposal != wanted.IsDailyProposal {
						t.Errorf("ParseMenuFile() %d has wrong IsDailyProposal: got %v, want %v", i, item.IsDailyProposal, wanted.IsDailyProposal)
					}
				}
			}
		})
	}
}