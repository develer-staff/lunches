package tuttobene

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/shopspring/decimal"
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
					{"Rigatoni al ragù dell'aia", Primo, false, decimal.NewFromFloat32(7)},
					{"Ravioli ricotta e spinaci con burro e salvia", Primo, false, decimal.NewFromFloat32(7.5)},
					{"Lasagne con cavolo nero e porri", Primo, false, decimal.NewFromFloat32(7)},
					{"Minestra di pane", Primo, false, decimal.NewFromFloat32(7)},
					{"Paccheri con calamari e asparagi", Primo, false, decimal.NewFromFloat32(8.5)},
					{"Pasta al ragù", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pesto", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pomodoro", Primo, false, decimal.NewFromFloat32(7)},
					{"Lasagne cavolo nero e porri + macedonia", Primo, true, decimal.NewFromFloat32(8.9)},
					{"Roastbeef con patate arrosto", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Polpette in umido con verdure", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Spezzatino di vitella con asparagi", Secondo, false, decimal.NewFromFloat32(11)},
					{"Baccalà alla livornese con fagioli", Secondo, false, decimal.NewFromFloat32(12)},
					{"Filetto di branzino gratinato con fagiolini", Secondo, false, decimal.NewFromFloat32(12)},
					{"Baccalà alla livornese con fagioli + macedonia", Secondo, true, decimal.NewFromFloat32(10.90)},
					{"Sformatini di riso con verdure al vapore", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Fantasia di verdure grigliate", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Macedonia di frutta fresca", Frutta, false, decimal.NewFromFloat32(4)},
					{"Macedonia di frutta fresca piccola", Frutta, false, decimal.NewFromFloat32(2)},
					{"Frutta a tocchi", Frutta, false, decimal.NewFromFloat32(4)},
					{"Diametro 12 mortadella", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 12 crudo pecorino e rucola", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Diametro 8 bresaola rucola e brie", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 8 vegetariano", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Tubo 15 praga radicchi e grana", Panino, false, decimal.NewFromFloat32(3.8)},
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
					{"Sedani alla Carloforte", Primo, false, decimal.NewFromFloat32(7.5)},
					{"Strigoli con filangè di verdure e speck", Primo, false, decimal.NewFromFloat32(7)},
					{"Orecchiette alle rape", Primo, false, decimal.NewFromFloat32(7)},
					{"Zuppa di zucca con pane croccante", Primo, false, decimal.NewFromFloat32(7)},
					{"Paccheri alla triglia", Primo, false, decimal.NewFromFloat32(8.5)},
					{"Pasta al ragù", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pesto", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pomodoro", Primo, false, decimal.NewFromFloat32(7)},
					{"Orecchiette alle rape + macedonia", Primo, true, decimal.NewFromFloat32(8.9)},
					{"Polpette in umido con purè", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Ossibuchi alla livornese con fagioli borlotti", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Filetto di maiale con panure a i 3 pepi e patate arrosto", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Orata all'isolana con spinaci", Secondo, false, decimal.NewFromFloat32(12)},
					{"Seppie con piselli", Secondo, false, decimal.NewFromFloat32(12)},
					{"Polpette in umido con purè + macedonia", Secondo, true, decimal.NewFromFloat32(10.9)},
					{"Insalata di spinacina, fagioli di soja, feta e mais", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Dadolata di verdure al forno", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Macedonia di frutta fresca", Frutta, false, decimal.NewFromFloat32(4)},
					{"Macedonia di frutta fresca piccola", Frutta, false, decimal.NewFromFloat32(2)},
					{"Frutta a tocchi", Frutta, false, decimal.NewFromFloat32(4)},
					{"Diametro 12 mortadella", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 12 crudo pecorino e rucola", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Diametro 8 bresaola rucola e brie", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 8 vegetariano", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Tubo 15 praga radicchi e grana", Panino, false, decimal.NewFromFloat32(3.8)},
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
					{"Penne con salsiccia e rape", Primo, false, decimal.NewFromFloat32(7)},
					{"Pici cacio e pepe", Primo, false, decimal.NewFromFloat32(7)},
					{"Crespelle alla fiorentina", Primo, false, decimal.NewFromFloat32(7.5)},
					{"Minestrone", Primo, false, decimal.NewFromFloat32(7)},
					{"Paccheri al polpo", Primo, false, decimal.NewFromFloat32(8.5)},
					{"Pasta al ragù", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pesto", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pomodoro", Primo, false, decimal.NewFromFloat32(7)},
					{"Penne con salsiccia e rape + macedonia", Primo, true, decimal.NewFromFloat32(8.9)},
					{"Pollo al curry con riso nero", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Hamburger con pomodori grigliati", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Bianchetto di vitellla con champignon", Secondo, false, decimal.NewFromFloat32(11)},
					{"Moscardini con piselli", Secondo, false, decimal.NewFromFloat32(12)},
					{"Spada alla griglia con belga", Secondo, false, decimal.NewFromFloat32(12)},
					{"Hamburger con pomodori grigliati + macedonia", Secondo, true, decimal.NewFromFloat32(10.9)},
					{"Insalata di zucca gialla con pomodori e olive", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Fantasia di verdure al vapore", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Macedonia di frutta fresca", Frutta, false, decimal.NewFromFloat32(4)},
					{"Macedonia di frutta fresca piccola", Frutta, false, decimal.NewFromFloat32(2)},
					{"Frutta a tocchi", Frutta, false, decimal.NewFromFloat32(4)},
					{"Diametro 12 mortadella", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 12 crudo pecorino e rucola", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Diametro 8 bresaola rucola e brie", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 8 vegetariano", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Tubo 15 praga radicchi e grana", Panino, false, decimal.NewFromFloat32(3.8)},
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
					{"Penne con salsiccia e rape", Primo, false, decimal.NewFromFloat32(0)},
					{"Pici cacio e pepe", Primo, false, decimal.NewFromFloat32(0)},
					{"Crespelle alla fiorentina", Primo, false, decimal.NewFromFloat32(0)},
					{"Minestrone", Primo, false, decimal.NewFromFloat32(0)},
					{"Paccheri al polpo", Primo, false, decimal.NewFromFloat32(0)},
					{"Pasta olio", Primo, false, decimal.NewFromFloat32(0)},
					{"Pasta al ragù", Primo, false, decimal.NewFromFloat32(0)},
					{"Riso olio", Primo, false, decimal.NewFromFloat32(0)},
					{"Pasta al pomodoro", Primo, false, decimal.NewFromFloat32(0)},
					{"Pollo al curry", Secondo, false, decimal.NewFromFloat32(0)},
					{"Hamburger", Secondo, false, decimal.NewFromFloat32(0)},
					{"Bianchetto di vitellla", Secondo, false, decimal.NewFromFloat32(0)},
					{"Moscardini con piselli", Secondo, false, decimal.NewFromFloat32(0)},
					{"Spada alla griglia", Secondo, false, decimal.NewFromFloat32(0)},
					{"Peperoni alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Melanzane alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Belga alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Radicchio alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Broccoli al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Cavolfiore al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Carote al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Fagiolini al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Dadolata di verdure al forno", Contorno, false, decimal.NewFromFloat32(0)},
					{"Pomodori", Contorno, false, decimal.NewFromFloat32(0)},
					{"Insalata", Contorno, false, decimal.NewFromFloat32(0)},
					{"Patate arrosto", Contorno, false, decimal.NewFromFloat32(0)},
					{"Spinaci saltati", Contorno, false, decimal.NewFromFloat32(0)},
					{"Pomodori grigliati", Contorno, false, decimal.NewFromFloat32(0)},
					{"Insalata di zucca gialla con pomodori e olive", Vegetariano, false, decimal.NewFromFloat32(0)},
					{"Fantasia di verdure al vapore", Vegetariano, false, decimal.NewFromFloat32(0)},
					{"Mozzarelle", Vegetariano, false, decimal.NewFromFloat32(0)},
					{"Macedonia di frutta fresca", Frutta, false, decimal.NewFromFloat32(0)},
					{"Macedonia di frutta fresca piccola", Frutta, false, decimal.NewFromFloat32(0)},
					{"Frutta a tocchi", Frutta, false, decimal.NewFromFloat32(0)},
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
					{"Penne all'amatriciana", Primo, false, decimal.NewFromFloat32(0)},
					{"Sedani salsiccia e olive", Primo, false, decimal.NewFromFloat32(0)},
					{"Paccheri zucchine e speck", Primo, false, decimal.NewFromFloat32(0)},
					{"Farro alla sorrentina (freddo)", Primo, false, decimal.NewFromFloat32(0)},
					{"Spaghetti allo scoglio", Primo, false, decimal.NewFromFloat32(0)},
					{"Pasta olio", Primo, false, decimal.NewFromFloat32(0)},
					{"Pasta al ragù", Primo, false, decimal.NewFromFloat32(0)},
					{"Pasta al pomodoro", Primo, false, decimal.NewFromFloat32(0)},
					{"Riso olio", Primo, false, decimal.NewFromFloat32(0)},
					{"Spiedini di carne", Secondo, false, decimal.NewFromFloat32(0)},
					{"Roastbeef", Secondo, false, decimal.NewFromFloat32(0)},
					{"Pollo ripieno", Secondo, false, decimal.NewFromFloat32(0)},
					{"Tagliata di tonno", Secondo, false, decimal.NewFromFloat32(0)},
					{"Salmone al vapore", Secondo, false, decimal.NewFromFloat32(0)},
					{"Tonno sott'olio", Secondo, false, decimal.NewFromFloat32(0)},
					{"Bresaola", Secondo, false, decimal.NewFromFloat32(0)},
					{"Prociutto crudo", Secondo, false, decimal.NewFromFloat32(0)},
					{"Peperoni alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Melanzane alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Belga alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Finocchi alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Radicchio alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Broccoli al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Cavolfiore al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Carote al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Fagiolini al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Pomodori", Contorno, false, decimal.NewFromFloat32(0)},
					{"Insalata", Contorno, false, decimal.NewFromFloat32(0)},
					{"Patate arrosto", Contorno, false, decimal.NewFromFloat32(0)},
					{"Piselli", Contorno, false, decimal.NewFromFloat32(0)},
					{"Spinaci saltati", Contorno, false, decimal.NewFromFloat32(0)},
					{"Taccole al pomodoro", Contorno, false, decimal.NewFromFloat32(0)},
					{"Primosale con insalata mista", Vegetariano, false, decimal.NewFromFloat32(0)},
					{"Dadolata di verdure al forno", Vegetariano, false, decimal.NewFromFloat32(0)},
					{"Mozzarelle", Vegetariano, false, decimal.NewFromFloat32(0)},
					{"Macedonia di frutta fresca", Frutta, false, decimal.NewFromFloat32(0)},
					{"Macedonia di frutta fresca piccola", Frutta, false, decimal.NewFromFloat32(0)},
					{"Frutta a tocchi", Frutta, false, decimal.NewFromFloat32(0)},
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
		{
			"Test new format",
			args{filepath.Join("test-fixtures", "testmenuv3.xlsx")},
			2019,
			&Menu{

				[]MenuRow{
					{"Fusilli con ricotta rucola e pinoli (freddo) + macedonia", Primo, true, decimal.NewFromFloat32(8.9)},
					{"Couscous con tonno pomodori e olive(freddo)", Primo, false, decimal.NewFromFloat32(7)},
					{"Fusilli con ricotta rucola e pinoli (freddo)", Primo, false, decimal.NewFromFloat32(7)},
					{"Sedani all'amatriciana", Primo, false, decimal.NewFromFloat32(7)},
					{"Paella catalana", Primo, false, decimal.NewFromFloat32(10)},
					{"Paccheri alla Carloforte", Primo, false, decimal.NewFromFloat32(8.5)},
					{"Pasta olio", Primo, false, decimal.NewFromFloat32(5)},
					{"Pasta al pesto", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al ragù", Primo, false, decimal.NewFromFloat32(7)},
					{"Pasta al pomodoro", Primo, false, decimal.NewFromFloat32(6)},
					{"Riso olio", Primo, false, decimal.NewFromFloat32(5)},

					{"Roastbeef con contorno a piacere + macedonia", Secondo, true, decimal.NewFromFloat32(10.9)},
					{"Insalata con mozzarella, tonno, pomodori (o scegli tu fra: uovo sodo, mais, semi vari)", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Cosciotto di maiale del Mugello", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Roastbeef", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"Tasca di tacchinoalla ligure", Secondo, false, decimal.NewFromFloat32(9.5)},
					{"polpo con piselli e olive", Secondo, false, decimal.NewFromFloat32(12)},
					{"Baccalà alla livornese", Secondo, false, decimal.NewFromFloat32(12)},

					{"Peperoni alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Melanzane alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Belga alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Finocchi alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Radicchio alla griglia", Contorno, false, decimal.NewFromFloat32(0)},
					{"Broccoli al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Cavolfiore al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Carote al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Fagiolini al vapore", Contorno, false, decimal.NewFromFloat32(0)},
					{"Pomodori", Contorno, false, decimal.NewFromFloat32(0)},
					{"Insalata mista", Contorno, false, decimal.NewFromFloat32(0)},
					{"Taccole con pomodorini", Contorno, false, decimal.NewFromFloat32(0)},
					{"Dadolata di verdure al forno", Contorno, false, decimal.NewFromFloat32(0)},
					{"Patate arrosto", Contorno, false, decimal.NewFromFloat32(0)},
					{"Spinaci saltati", Contorno, false, decimal.NewFromFloat32(0)},
					{"Ceci", Contorno, false, decimal.NewFromFloat32(0)},
					{"Spinaci con patate", Contorno, false, decimal.NewFromFloat32(0)},

					{"Insalata greca", Vegetariano, false, decimal.NewFromFloat32(9.5)},
					{"Verdure al vapore", Vegetariano, false, decimal.NewFromFloat32(9.5)},

					{"Macedonia di frutta fresca", Frutta, false, decimal.NewFromFloat32(4)},
					{"Macedonia di frutta fresca piccola", Frutta, false, decimal.NewFromFloat32(2)},
					{"Frutta a tocchi", Frutta, false, decimal.NewFromFloat32(4)},

					{"Schiacciata con l'uva", Dolce, false, decimal.NewFromFloat32(2.5)},
					{"Shiacciata con i fichi", Dolce, false, decimal.NewFromFloat32(2.5)},

					{"Diametro 12 mortadella", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 12 crudo pecorino e rucola", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Diametro 8 bresaola rucola e brie", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Diametro 8 vegetariano", Panino, false, decimal.NewFromFloat32(3.5)},
					{"Tubo 15 tonno maionese e pomodoro", Panino, false, decimal.NewFromFloat32(3.8)},
					{"Tubo 15 praga radicchi e grana", Panino, false, decimal.NewFromFloat32(3.8)},
				},
				time.Date(2019, 9, 20, 0, 0, 0, 0, loc),
			},
			false,
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

			if tt.want != nil && tt.want.String() != got.String() {
				t.Errorf("ParseMenuFile(): %v, want: %v", got.String(), tt.want.String())
			}
		})
	}
}
