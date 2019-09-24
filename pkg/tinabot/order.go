package tinabot

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type DataStore interface {
	Set(string, interface{}) error
	Get(string, interface{}) error
}

// User data
type User struct {
	Name string
	ID   string
}

func (u User) MarshalText() ([]byte, error) {
	return []byte(u.Name + "&&&&" + u.ID), nil
}

func (u *User) UnmarshalText(text []byte) error {
	js := string(text)

	if js == "null" {
		return nil
	}

	f := strings.Split(js, "&&&&")
	if len(f) != 2 {
		return errors.New("invalid User field")
	}

	*u = User{f[0], f[1]}
	return nil
}

// Order is a structure holding Tinabot orders
type Order struct {
	Timestamp time.Time
	Dishes    map[string][]User        //map dishes with users
	Users     map[User]UserChoiceArray //map each user to his/her dishes
}

// NewOrder returns a new empty order
func NewOrder() *Order {
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Println("LoadLocation error: ", err)
		return nil
	}

	return &Order{
		Timestamp: time.Now().In(loc),
		Dishes:    make(map[string][]User),
		Users:     make(map[User]UserChoiceArray),
	}
}

// ClearUser clear the user order, returns the cleared dishes, if any
func (order *Order) ClearUser(user User) string {
	var deleted []string

	for _, d := range order.sorted() {
		users := order.Dishes[d]
		for i, u := range users {
			if u == user {
				deleted = append(deleted, d)
				order.Dishes[d] = append(order.Dishes[d][:i], order.Dishes[d][i+1:]...)
				break
			}
		}
		if len(order.Dishes[d]) == 0 {
			delete(order.Dishes, d)
		}
	}

	delete(order.Users, user)
	return strings.Join(deleted, "\n")
}

// sorted return an array of ordered dished sorted by dish type, dishname
func (order *Order) sorted() []string {
	// Create a map of ordered string -> rendered string
	dishmap := make(map[string]string)
	for _, choices := range order.Users {
		for _, c := range choices {
			dishmap[c.OrdString()] = c.String()
		}
	}

	// extract from the map all the ordered strings
	var ordstring []string
	for k := range dishmap {
		ordstring = append(ordstring, k)
	}

	// sort them
	sort.Slice(ordstring, func(i, j int) bool {
		return strings.Compare(ordstring[i], ordstring[j]) < 0
	})

	// return the ordered rendered strings
	var out []string
	for _, d := range ordstring {
		out = append(out, dishmap[d])
	}
	return out
}

// Load loads order from redis brain
func (order *Order) Load(brain DataStore) error {
	err := brain.Get("order", order)
	if err != nil {
		return err
	}

	return nil
}

// Save saves order to redis brain
func (order *Order) Save(brain DataStore) error {
	fmt.Println("save")
	return brain.Set("order", *order)
}

// Set set the current order for user to her choice, returns a string array of what she ordered
func (order *Order) Set(user User, choice []UserChoice) []string {
	order.ClearUser(user)
	var list []string
	for _, c := range choice {
		order.Dishes[c.String()] = append(order.Dishes[c.String()], user)
		order.Users[user] = append(order.Users[user], c)
		list = append(list, c.String())
	}

	return list
}

func (order *Order) String() string {
	return order.Format(true, false)
}

func (order *Order) Bill() string {
	return order.Format(true, true)
}

// Format convert the order to a string, with or without the user names
func (order *Order) Format(withUserNames, withPrices bool) string {
	var r []string
	var noPrice []string
	total := decimal.Zero

	for _, d := range order.sorted() {
		l := fmt.Sprintf("%d %s", len(order.Dishes[d]), d)
		if withUserNames {
			//gather names
			var names []string
			for _, u := range order.Dishes[d] {
				names = append(names, u.Name)
			}
			l += " [" + strings.Join(names, ", ") + "]"
		}

		if withPrices {
			cnt := len(order.Dishes[d])
			mul := decimal.New(int64(cnt), 0)
			priceFound := false

			u := order.Dishes[d][0]
			for _, dish := range order.Users[u] {
				if dish.String() == d {
					row := dish.Price().Mul(mul)
					total = total.Add(row)
					if !row.IsZero() {
						l += " -> €" + row.String()
						priceFound = true
						break
					}
				}
			}

			if !priceFound {
				l += " -> *prezzo non disponibile!*"
				noPrice = append(noPrice, d)
			}
		}
		r = append(r, l)
	}

	if withPrices {
		r = append(r, fmt.Sprintf("*Prezzo TOTALE: €%s*", total.String()))
		if len(noPrice) > 0 {
			r = append(r, "I seguenti piatti non hanno un prezzo indicato:")
			r = append(r, noPrice...)
		}
	}

	return strings.Join(r, "\n")
}

// IsUpdated returns true if it's today's order, false otherwise
func (order *Order) IsUpdated() bool {
	loc, err := time.LoadLocation("Europe/Rome")

	if err != nil {
		log.Println("LoadLocation error: ", err)
		return false
	}

	y, m, d := time.Now().In(loc).Date()
	ts := order.Timestamp
	return (y == ts.Year() && m == ts.Month() && d == ts.Day())
}
