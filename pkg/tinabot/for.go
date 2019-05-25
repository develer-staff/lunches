package tinabot

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"

	"github.com/develersrl/lunches/pkg/slackbot"
	"github.com/develersrl/lunches/pkg/tuttobene"
)

func unescape(s, sep string) string {

	s = strings.Replace(s, "\\"+sep, sep, -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}

func splitEsc(s, sep string) []string {
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
			a = append(a, unescape(s[startcp:m], sep))
			startcp = m + len(sep)
		}
		start = m + len(sep)
		i++
	}

	a = append(a, unescape(s[startcp:], sep))
	return a
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

func (t *TinaBot) For(bot *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
	dest := args[1]
	dish := sanitize(args[2])

	destUser := user
	destName := user.Name
	destCh := ""

	if strings.ToLower(dest) != "me" {
		destUser = getUserInfo(t.bot.Client, dest)
		if destUser != nil {
			destName = destUser.Name
			_, _, ch, err := bot.Client.OpenIMChannel(destUser.ID)
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
	reply := ""

	// handle the "copy" order command
	if strings.HasPrefix(strings.ToLower(dish), "come") {
		l := strings.Split(dish, " ")
		if len(l) < 2 {
			t.bot.Message(msg.Channel, fmt.Sprintf("E' necessario specificare da chi vuoi copiare l'ordine"))
			return
		}
		name := l[1]
		if strings.HasPrefix(name, "<@") {
			user := getUserInfo(bot.Client, name)
			if user == nil {
				t.bot.Message(msg.Channel, fmt.Sprintf("Mi spiace, ma non riesco a trovare l'utente %s", name))
				return
			}
			name = user.Name
		}
		order := getOrder(t.brain)
		if newchoice, ok := order.Users[name]; ok {
			reply = reply + fmt.Sprintf("Ok, copio l'ordine di %s:\n", name)
			for _, c := range newchoice {
				reply = reply + c.String() + "\n"
			}
			choice = newchoice
		} else {
			t.bot.Message(msg.Channel, fmt.Sprintf("Mi spiace, ma non trovo l'utente '%s' nell'ordine", name))
			return
		}
	} else {
		reqs := splitEsc(dish, "+")

		for _, req := range reqs {
			dishes := splitEsc(req, "&amp;")
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
}
