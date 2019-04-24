package tinabot

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/nlopes/slack"

	"github.com/develersrl/lunches/pkg/brain"
	"github.com/develersrl/lunches/pkg/slackbot"
	"github.com/develersrl/lunches/pkg/tuttobene"
	"github.com/go-redis/redis"
)

func getOrder(brain *brain.Brain) *Order {
	var order Order

	if order.Load(brain) != nil {
		return NewOrder()
	}

	if !order.IsUpdated() {
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
	for _, m := range menu.Rows {
		if strings.EqualFold(m.Content, dish) {
			return []tuttobene.MenuRow{m}
		}

		if fuzzyMatch(dish, m.Content) {
			matches = append(matches, m)
		}
	}
	return matches
}

func Sanitize(s string) string {
	s = strings.Replace(s, "“", "\"", -1)
	s = strings.Replace(s, "”", "\"", -1)
	s = strings.Replace(s, "‘", "'", -1)
	s = strings.Replace(s, "’", "'", -1)
	return s
}

func Unescape(s, sep string) string {

	s = strings.Replace(s, "\\"+sep, sep, -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}

func SplitEsc(s, sep string) []string {
	escC := byte('\\')

	n := strings.Count(s, sep)
	var a []string
	i := 0
	start := 0
	startcp := 0

	for i < n {
		m := strings.Index(s[start:], sep)
		if m < 0 {
			break
		}
		m += start
		if m == 0 || (m > 0 && s[m-1] != escC) {
			a = append(a, Unescape(s[startcp:m], sep))
			startcp = m + len(sep)
		}
		start = m + len(sep)
		i++
	}

	a = append(a, Unescape(s[startcp:], sep))
	return a
}

type TinaBot struct {
	bot   *slackbot.Bot
	brain *brain.Brain
}

func New(bot *slackbot.Bot, b *brain.Brain) *TinaBot {
	return &TinaBot{bot, b}
}

func getUserInfo(api *slack.Client, user string) *slack.User {
	if strings.HasPrefix(user, "<@") {
		user = strings.Trim(user, "<@>")
		u, err := api.GetUserInfo(user)
		if err != nil {
			log.Println(err)
			return nil
		}
		return u
	}

	users, err := api.GetUsers()
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, u := range users {
		if strings.ToLower(u.Name) == strings.ToLower(user) {
			return &u
		}
	}
	return nil
}

func (t *TinaBot) AddCommands() {

	t.bot.DefaultResponse(func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User) {
		t.bot.Message(msg.Channel, "Mi dispiace "+user.Name+", purtroppo non posso farlo.\nProva con `aiuto` per vedere l'elenco delle cose che posso fare.")
	})

	t.bot.RespondTo("^(?i)(help|aiut).*$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		t.bot.Message(msg.Channel, strings.Replace(HelpStr, "‘", "`", -1))
	})

	t.bot.RespondTo("^(?i)per (\\S+) (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		dest := args[1]
		dish := Sanitize(args[2])

		destUser := user
		destName := user.Name
		destCh := ""

		if strings.ToLower(dest) != "me" {
			destUser = getUserInfo(t.bot.Client, dest)
			if destUser != nil {
				destName = destUser.Name
				_, _, ch, err := b.Client.OpenIMChannel(destUser.ID)
				if err != nil {
					log.Println(err)
				} else {
					destCh = ch
				}
			} else {
				if !strings.HasPrefix(dest, "guest_") {
					t.bot.Message(msg.Channel, fmt.Sprintf("Utente '%s' non trovato. Se vuoi ordinare per conto di un ospite usa il prefisso *guest_* nel nome", dest))
					return
				}
				destName = dest
			}
		}

		if strings.ToLower(dish) == "niente" {
			order := getOrder(t.brain)
			old := order.ClearUser(destName)
			order.Save(t.brain)

			t.bot.Message(msg.Channel, fmt.Sprintf("Ok, cancello ordine per %s:\n%s", destName, old))
			if destCh != "" {
				t.bot.Message(destCh, fmt.Sprintf("Mi spiace disturbarti, volevo informarti che <@%s> ha appena cancellato il tuo ordine:\n%s", user.ID, old))
			}
			return
		}

		var menu tuttobene.Menu
		err := t.brain.Get("menu", &menu)
		if err != nil {
			t.bot.Message(msg.Channel, "Nessun menù impostato!")
			return
		}

		if !menu.IsUpdated() {
			t.bot.Message(msg.Channel, "Non puoi ordinare, il menù non è quello di oggi, riporta la data del "+menu.Date.Format("02/01/2006"))
			return
		}

		var choice []UserChoice
		reqs := SplitEsc(dish, "+")

		reply := ""
		for _, req := range reqs {
			dishes := SplitEsc(req, "&amp;")
			var currChoice UserChoice
			for _, dish := range dishes {
				dish = strings.TrimSpace(dish)

				quoted := (dish[0] == '"' && dish[len(dish)-1] == '"')
				dish = strings.Trim(dish, "\"")

				found := findDishes(menu, dish)
				nDish := len(found)

				if quoted && nDish != 1 {
					p := tuttobene.MenuRow{
						Content:         dish,
						Type:            tuttobene.Empty,
						IsDailyProposal: false,
					}
					reply = reply + fmt.Sprintf("Aggiungo testualmente: '%s'\n", dish)
					currChoice.Add(p)
				} else if nDish == 0 {
					t.bot.Message(msg.Channel, reply+"Non ho trovato nulla nel menù che corrisponda a '"+dish+"'\nOrdine non aggiunto!")
					return
				} else if nDish > 1 {
					var matches []string
					for _, d := range found {
						matches = append(matches, d.Content)
					}

					t.bot.Message(msg.Channel, reply+"Cercando per '"+dish+"' ho trovato i seguenti piatti:\n"+strings.Join(matches, "\n")+"\n----\nOrdine non aggiunto, prova ad essere più preciso!")
					return
				} else { // nDish == 1
					d := found[0]
					reply = reply + "Trovato: " + d.Content + fmt.Sprintf(" (%s)\n", tuttobene.Titles[d.Type])

					err := currChoice.Add(d)
					if err != nil {
						t.bot.Message(msg.Channel, reply+"Errore nella personalizzazione: "+err.Error()+"\nOrdine non aggiunto!")
						return
					}
				}
			}
			if currChoice.Customized() {
				reply = reply + "Piatto personalizzato: " + currChoice.String() + "\n"
			}
			choice = append(choice, currChoice)
		}

		order := getOrder(t.brain)
		list := order.Set(destName, choice)
		order.Save(t.brain)

		l := len(choice)
		c := "o"
		if l > 1 {
			c = "i"
		}
		t.bot.Message(msg.Channel, reply+fmt.Sprintf("Ok, aggiunt%s %d piatt%s per %s", c, l, c, destName))
		if destCh != "" {
			t.bot.Message(destCh, fmt.Sprintf("Ti volevo informare che <@%s> ha ordinato i seguenti piatti per conto tuo:\n%s", user.ID, strings.Join(list, "\n")))
		}
	})

	t.bot.RespondTo("^(?i)ordine$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(t.brain)
		t.bot.Message(msg.Channel, "Ecco l'ordine:\n"+order.String())
	})

	t.bot.RespondTo("^(?i)cancella ordine$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := NewOrder()
		order.Save(t.brain)
		t.bot.Message(msg.Channel, "Ordine cancellato")
	})

	t.bot.RespondTo("^(?i)email$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(t.brain)
		subj := "Ordine Develer del giorno " + order.Timestamp.Format("02/01/2006")
		body := order.Format(false)

		out := subj + "\n" + body + "\n\n" +
			"<mailto:info@tuttobene-bar.it,sara@tuttobene-bar.it" +
			"?subject=" + url.PathEscape(subj) +
			"&body=" + url.PathEscape(body) +
			"|Link `mailto` clickabile>"
		t.bot.Message(msg.Channel, out)
	})

	t.bot.RespondTo("^(?i)menu([\\s\\S]*)?", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {

		if args[1] != "" {
			t.bot.Message(msg.Channel, "Se stai cercando di impostare il menù, usa il comando `setmenu`\nPer vedere il menù corrente, usa il comando `menu` senza argomenti.")
			return
		}

		var m tuttobene.Menu
		err := t.brain.Get("menu", &m)
		if err == redis.Nil {
			t.bot.Message(msg.Channel, "Non c'è nessun menù impostato!")
		} else {
			t.bot.Message(msg.Channel, "Ecco il menù:\n"+m.String())
		}
	})

	t.bot.RespondTo("^(?i)setmenu([\\s\\S]*)?", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		if args[1] != "" {
			menu := strings.Split(strings.TrimSpace(Sanitize(args[1])), "\n")
			m, err := tuttobene.ParseMenuRows(menu)
			if err != nil {
				t.bot.Message(msg.Channel, "Menu parse error: "+err.Error())
				return
			}
			t.brain.Set("menu", *m)
			t.bot.Message(msg.Channel, "Ok, menù impostato:\n"+m.String())
		} else {
			t.bot.Message(msg.Channel, "Non hai indicato nessun nuovo menù!")
		}
	})

	t.bot.RespondTo("^set (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		ar := strings.Split(args[1], " ")
		key := ar[0]
		val := ar[1]
		err := t.brain.Set(key, val)
		if err != nil {
			t.bot.Message(msg.Channel, "Error: "+err.Error())
		} else {
			t.bot.Message(msg.Channel, "Ok")
		}
	})

	t.bot.RespondTo("^get (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		key := args[1]
		var val string
		err := t.brain.Get(key, &val)
		if err != nil {
			t.bot.Message(msg.Channel, "Error: "+err.Error())
		} else {
			t.bot.Message(msg.Channel, key+": "+val)
		}
	})

	t.bot.RespondTo("^read (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		key := args[1]

		val, err := t.brain.Read(key)
		if err != nil {
			t.bot.Message(msg.Channel, "Error: "+err.Error())
		} else {
			t.bot.Message(msg.Channel, key+": "+val)
		}
	})

	t.bot.RespondTo("^(?i)cron(.*)$", t.Cron)

	t.bot.RespondTo("^(?i)rmorder (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		u := args[1]

		order := getOrder(t.brain)
		old := order.ClearUser(u)
		t.bot.Message(msg.Channel, fmt.Sprintf("Ok, cancello ordine di %s:\n%s", u, old))
		order.Save(t.brain)
	})
}
