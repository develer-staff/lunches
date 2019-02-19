package tinabot

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/nlopes/slack"

	"github.com/develersrl/lunches/actions/brain"
	"github.com/develersrl/lunches/actions/slackbot"
	"github.com/develersrl/lunches/pkg/tuttobene"
	"github.com/go-redis/redis"
)

type UserChoice struct {
	DishMask uint
	Dishes   []tuttobene.MenuRow
}

func (u *UserChoice) Add(dish tuttobene.MenuRow) error {
	if (dish.Type == tuttobene.Primo && u.DishMask != 0) ||
		(dish.Type == tuttobene.Secondo && (u.DishMask&^(1<<uint(tuttobene.Contorno))) != 0) ||
		(dish.Type == tuttobene.Contorno && (u.DishMask&^(1<<uint(tuttobene.Contorno)|1<<uint(tuttobene.Secondo))) != 0) {
		return errors.New("è possibile solo comporre piatti formati da un secondo e contorni")
	}

	u.DishMask |= (1 << uint(dish.Type))
	u.Dishes = append(u.Dishes, dish)
	return nil
}

func (u *UserChoice) Sort() {
	sort.Slice(u.Dishes, func(i, j int) bool {
		si := fmt.Sprintf("%d%s", u.Dishes[i].Type, u.Dishes[i].Content)
		sj := fmt.Sprintf("%d%s", u.Dishes[j].Type, u.Dishes[j].Content)
		return strings.Compare(si, sj) < 0
	})
}

func (u *UserChoice) String() string {
	u.Sort()
	var main []string
	var side []string
	for _, d := range u.Dishes {
		if d.Type == tuttobene.Contorno {
			side = append(side, d.Content)
		} else {
			main = append(main, d.Content)
		}
	}
	out := strings.Join(main, ",")
	if len(side) > 0 {
		if len(main) > 0 {
			out += " con contorno di: "
		}
		out += strings.Join(side, ",")
	}
	return out
}

type Order struct {
	Timestamp time.Time
	Dishes    map[string][]string     //map dishes with users
	Users     map[string][]UserChoice //map each user to his/her dishes
}

func NewOrder() *Order {
	return &Order{
		Timestamp: time.Now(),
		Dishes:    make(map[string][]string),
		Users:     make(map[string][]UserChoice),
	}
}

func getOrder(brain *brain.Brain) *Order {
	var order Order
	err := brain.Get("order", &order)
	if err != nil {
		return NewOrder()
	}

	if time.Since(order.Timestamp).Hours() > 13 {
		log.Println("Deleting old order")
		return NewOrder()
	}
	return &order
}

func fuzzyMatch(dish, menuline string) bool {
	dish = strings.ToLower(dish)

	key := regexp.MustCompile(strings.Replace(regexp.QuoteMeta(dish), " ", ".*", -1))

	return key.MatchString(strings.ToLower(menuline))
}

func findDishes(menu tuttobene.Menu, dish string) []tuttobene.MenuRow {
	dish = strings.TrimSpace(strings.ToLower(dish))

	var matches []tuttobene.MenuRow
	for _, m := range menu {
		if strings.EqualFold(m.Content, dish) {
			return []tuttobene.MenuRow{m}
		}

		if fuzzyMatch(dish, m.Content) {
			matches = append(matches, m)
		}
	}
	return matches
}

func clearUserOrder(order *Order, user string) string {
	dishes := order.Users[user]
	delete(order.Users, user)
	var deleted []string
	for _, uc := range dishes {

		d := uc.String()
		deleted = append(deleted, d)
		// Find and remove the user
		for i, v := range order.Dishes[d] {
			if v == user {
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

func renderMenu(menu tuttobene.Menu) string {
	menutype := tuttobene.Unknonwn

	out := ""
	for _, r := range menu {
		if r.Type != menutype {
			out = out + "\n*" + strings.ToUpper(tuttobene.Titles[r.Type]) + "*\n"
			menutype = r.Type
		}
		out = out + r.Content + "\n"
	}
	return out
}

func Tinabot(bot *slackbot.Bot) {

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatalln("No redis URL found!")
	}
	log.Printf("Redis URL: %s\n", redisURL)

	brain := brain.New(redisURL)

	bot.DefaultResponse(func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User) {
		bot.Message(msg.Channel, "Mi dispiace "+user.Name+", purtroppo non posso farlo.")
	})

	bot.RespondTo("^(?i)per me (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		dish := args[1]
		order := getOrder(brain)

		if strings.ToLower(dish) == "niente" {
			old := clearUserOrder(order, user.Name)
			bot.Message(msg.Channel, "Ok, cancello ordine:\n"+old)
			brain.Set("order", order)
			return
		}
		var menu tuttobene.Menu
		err := brain.Get("menu", &menu)
		if err != nil {
			bot.Message(msg.Channel, "Nessun menu impostato!")
			return
		}

		var choice []UserChoice
		dishes := strings.Split(dish, "&amp;&amp;")

		reply := ""
		for _, dish := range dishes {
			dishes := findDishes(menu, dish)

			if len(dishes) == 0 {
				bot.Message(msg.Channel, reply+"Non ho trovato nulla nel menu che corrisponda a '"+dish+"'\nOrdine non aggiunto!")
				return
			} else if len(dishes) > 1 {
				var matches []string
				for _, d := range dishes {
					matches = append(matches, d.Content)
				}

				bot.Message(msg.Channel, reply+"Cercando per '"+dish+"' ho trovato i seguenti piatti:\n"+strings.Join(matches, "\n")+"\n----\nOrdine non aggiunto, prova ad essere più preciso!")
				return
			} else {
				d := dishes[0]
				reply = reply + "Trovato: " + d.Content + fmt.Sprintf(" (%s)\n", tuttobene.Titles[d.Type])
				var c UserChoice
				err := c.Add(d)
				if err != nil {
					bot.Message(msg.Channel, reply+"errore nell'aggiunta "+err.Error())
					return
				}
				choice = append(choice, c)
			}
		}
		clearUserOrder(order, user.Name)
		u := user.Name
		for _, c := range choice {
			order.Dishes[c.String()] = append(order.Dishes[c.String()], u)
			order.Users[u] = append(order.Users[u], c)
		}
		brain.Set("order", order)
		l := len(choice)
		c := "o"
		if l > 1 {
			c = "i"
		}
		bot.Message(msg.Channel, reply+fmt.Sprintf("Ok, aggiunt%s %d piatt%s per %s", c, l, c, u))
	})

	bot.RespondTo("^(?i)ordine$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(brain)

		r := ""
		for d := range order.Dishes {
			l := fmt.Sprintf("%d %s ", len(order.Dishes[d]), d)
			l += "[ " + strings.Join(order.Dishes[d], ",") + " ]\n"
			r = r + l
		}

		bot.Message(msg.Channel, "Ecco l'ordine:\n"+r)
	})

	bot.RespondTo("^(?i)email$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(brain)
		subj := "Ordine Develer del giorno " + order.Timestamp.Format("02/01/2006")
		body := ""
		for d := range order.Dishes {
			body += fmt.Sprintf("%d %s\n", len(order.Dishes[d]), d)
		}
		out := subj + "\n" + body + "\n\n" +
			"<mailto:info@tuttobene-bar.it,sara@tuttobene-bar.it" +
			"?subject=" + url.PathEscape(subj) +
			"&body=" + url.PathEscape(body) +
			"|Link `mailto` clickabile>"
		bot.Message(msg.Channel, out)
	})

	bot.RespondTo("^(?i)menu([\\s\\S]*)?", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		var menu []string
		if args[1] != "" {
			menu = strings.Split(strings.TrimSpace(args[1]), "\n")
		} else {
			menu = nil
		}

		if menu == nil {
			var m tuttobene.Menu
			err := brain.Get("menu", &m)
			if err == redis.Nil {
				bot.Message(msg.Channel, "Non c'è nessun menu impostato!")
			} else {
				bot.Message(msg.Channel, "Il menu è:\n"+renderMenu(m))
			}
		} else {
			m, err := tuttobene.ParseMenuRows(menu)
			if err != nil {
				bot.Message(msg.Channel, "Menu parse error: "+err.Error())
				return
			}
			brain.Set("menu", *m)
			bot.Message(msg.Channel, "Ok, il menu è:\n"+renderMenu(*m))
		}
	})

	bot.RespondTo("^set (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		ar := strings.Split(args[1], " ")
		key := ar[0]
		val := ar[1]
		err := brain.Set(key, val)
		if err != nil {
			bot.Message(msg.Channel, "Error: "+err.Error())
		} else {
			bot.Message(msg.Channel, "Ok")
		}
	})

	bot.RespondTo("^get (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		key := args[1]
		var val string
		err := brain.Get(key, &val)
		if err != nil {
			bot.Message(msg.Channel, "Error: "+err.Error())
		} else {
			bot.Message(msg.Channel, key+": "+val)
		}
	})

	bot.RespondTo("^read (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		key := args[1]

		val, err := brain.Read(key)
		if err != nil {
			bot.Message(msg.Channel, "Error: "+err.Error())
		} else {
			bot.Message(msg.Channel, key+": "+val)
		}
	})
}
