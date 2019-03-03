package tinabot

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type DataStore interface {
	Set(string, interface{}) error
	Get(string, interface{}) error
}

// Order is a structure holding Tinabot orders
type Order struct {
	Timestamp time.Time
	Dishes    map[string][]string     //map dishes with users
	Users     map[string][]UserChoice //map each user to his/her dishes
}

// NewOrder returns a new empty order
func NewOrder() *Order {
	return &Order{
		Timestamp: time.Now(),
		Dishes:    make(map[string][]string),
		Users:     make(map[string][]UserChoice),
	}
}

// ClearUser clear the user order, returns the cleared dishes, if any
func (order *Order) ClearUser(user string) string {
	delete(order.Users, user)
	var deleted []string

	for d, users := range order.Dishes {
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
	return brain.Set("order", *order)
}

// Set set the current order for user to her choice, returns a string array of what she ordered
func (order *Order) Set(user string, choice []UserChoice) []string {
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
	return order.Format(true)
}

// Format convert the order to a string, with or without the user names
func (order *Order) Format(withUserNames bool) string {
	var r []string
	for _, d := range order.sorted() {
		l := fmt.Sprintf("%d %s", len(order.Dishes[d]), d)
		if withUserNames {
			l += " [" + strings.Join(order.Dishes[d], ", ") + "]"
		}
		r = append(r, l)
	}

	return strings.Join(r, "\n")
}
