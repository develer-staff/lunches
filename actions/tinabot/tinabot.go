package tinabot

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"

	"github.com/develersrl/lunches/actions/brain"
	"github.com/develersrl/lunches/actions/slackbot"
	"github.com/go-redis/redis"
)

type Order struct {
	Timestamp time.Time
	Dishes    map[string][]string //map dishes with users
	Users     map[string][]string //map each user to his/her dishes
}

func NewOrder() *Order {
	return &Order{
		Timestamp: time.Now(),
		Dishes:    make(map[string][]string),
		Users:     make(map[string][]string),
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

func findDishes(menu, dish string) []string {
	dish = strings.TrimSpace(strings.ToLower(dish))
	menus := strings.Split(strings.TrimSpace(menu), "\n")
	var matches []string
	for _, m := range menus {
		if strings.ToLower(m) == dish {
			return []string{m}
		}

		if fuzzyMatch(dish, m) {
			matches = append(matches, m)
		}
	}
	return matches
}

func clearUserOrder(order *Order, user string) string {
	dishes := order.Users[user]
	delete(order.Users, user)
	for _, d := range dishes {

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
	return strings.Join(dishes, "\n")
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
		var menu string
		err := brain.Get("menu", &menu)
		if err != nil {
			bot.Message(msg.Channel, "Nessun menu impostato!")
			return
		}

		var choice []string
		dishes := strings.Split(dish, "&amp;&amp;")

		reply := ""
		for _, dish := range dishes {
			dishes := findDishes(menu, dish)

			if len(dishes) == 0 {
				bot.Message(msg.Channel, reply+"Non ho trovato nulla nel menu che corrisponda a '"+dish+"'\nOrdine non aggiunto!")
				return
			} else if len(dishes) > 1 {
				matches := strings.Join(dishes, "\n")
				bot.Message(msg.Channel, reply+"Cercando per '"+dish+"' ho trovato i seguenti piatti:\n"+matches+"\n----\nOrdine non aggiunto, prova ad essere più preciso!")
				return
			} else {
				d := dishes[0]
				reply = reply + "Trovato: " + d + "\n"
				choice = append(choice, d)
			}
		}
		clearUserOrder(order, user.Name)
		u := user.Name
		for _, c := range choice {
			order.Dishes[c] = append(order.Dishes[c], u)
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
		var menu string
		if len(args) > 1 {
			menu = strings.TrimSpace(args[1])
		} else {
			menu = ""
		}

		if menu == "" {
			err := brain.Get("menu", &menu)
			if err == redis.Nil {
				bot.Message(msg.Channel, "Non c'è nessun menu impostato!")
			} else {
				bot.Message(msg.Channel, "Il menu è:\n"+menu)
			}
		} else {
			brain.Set("menu", menu)
			bot.Message(msg.Channel, "Ok, il menu è:\n"+menu)
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
