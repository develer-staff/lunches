package tinabot

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/develersrl/lunches/pkg/tuttobene"
)

// UserChoice is what the user choose to customize her dish from the menu (second with side dishes)
type UserChoice struct {
	DishMask uint
	Dishes   []tuttobene.MenuRow
}

// Clear clears the current user choice
func (u *UserChoice) Clear() {
	u.DishMask = 0
	u.Dishes = nil
}

// Customized returns true is the user choosed to customize her dish adding one or more side dishes
func (u *UserChoice) Customized() bool {
	return len(u.Dishes) > 1
}

// Add adds a dish the the choice.
func (u *UserChoice) Add(dish tuttobene.MenuRow) error {
	allowedMask := map[tuttobene.MenuRowType]uint{
		tuttobene.Empty:       0,
		tuttobene.Primo:       0,
		tuttobene.Secondo:     1<<uint(tuttobene.Contorno) | 1<<uint(tuttobene.Vegetariano),
		tuttobene.Contorno:    1<<uint(tuttobene.Secondo) | 1<<uint(tuttobene.Contorno) | 1<<uint(tuttobene.Vegetariano),
		tuttobene.Vegetariano: 1<<uint(tuttobene.Secondo) | 1<<uint(tuttobene.Contorno) | 1<<uint(tuttobene.Vegetariano),
		tuttobene.Frutta:      0,
		tuttobene.Panino:      0,
	}

	if u.DishMask&^allowedMask[dish.Type] != 0 {
		return errors.New("è possibile solo comporre piatti formati da un secondo e contorno/i")
	}

	u.DishMask |= (1 << uint(dish.Type))
	u.Dishes = append(u.Dishes, dish)
	return nil
}

func (u *UserChoice) sort() {
	sort.Slice(u.Dishes, func(i, j int) bool {
		si := fmt.Sprintf("%d%s", u.Dishes[i].Type, u.Dishes[i].Content)
		sj := fmt.Sprintf("%d%s", u.Dishes[j].Type, u.Dishes[j].Content)
		return strings.Compare(si, sj) < 0
	})
}

func (u *UserChoice) String() string {
	u.sort()
	var main []string
	var side []string
	for _, d := range u.Dishes {
		if d.Type == tuttobene.Secondo {
			main = append(main, d.Content)
		} else {
			side = append(side, d.Content)
		}
	}
	out := strings.Join(main, ", ")
	if len(side) > 0 {
		if len(main) > 0 {
			out += " con "
		}
		out += strings.Join(side, ", ")
	}
	return out
}

// OrdString return a string with a prefix that can be used to sort the dishes by category (first courses, second courses, fruit, etc... )
func (u *UserChoice) OrdString() string {
	return fmt.Sprintf("%04d-%s", u.DishMask, u.String())
}
