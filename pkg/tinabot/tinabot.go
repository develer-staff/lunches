package tinabot

import (
	"fmt"
	"log"
	"net/url"
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

func sanitize(s string) string {
	s = strings.Replace(s, "“", "\"", -1)
	s = strings.Replace(s, "”", "\"", -1)
	s = strings.Replace(s, "‘", "'", -1)
	s = strings.Replace(s, "’", "'", -1)
	return s
}

type TinaBot struct {
	bot   *slackbot.Bot
	brain *brain.Brain
}

func New(bot *slackbot.Bot, b *brain.Brain) *TinaBot {
	return &TinaBot{bot, b}
}

func (t *TinaBot) AddCommands() {

	t.bot.DefaultResponse(func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User) {
		t.bot.Message(msg.Channel, "Mi dispiace "+user.Name+", purtroppo non posso farlo.\nProva con `aiuto` per vedere l'elenco delle cose che posso fare.")
	})

	t.bot.RespondTo("^(?i)(help|aiut).*$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		t.bot.Message(msg.Channel, strings.Replace(HelpStr, "‘", "`", -1))
	})

	t.bot.RespondTo("^(?i)per (\\S+) (.*)$", t.For)

	t.bot.RespondTo("^(?i)ordine$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(t.brain)
		t.bot.Message(msg.Channel, "Ecco l'ordine:\n"+order.String())
	})

	t.bot.RespondTo("^(?i)conto$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(t.brain)
		t.bot.Message(msg.Channel, "Ecco il conto:\n"+order.Bill())
	})

	t.bot.RespondTo("^(?i)cancella ordine$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := NewOrder()
		order.Save(t.brain)
		t.bot.Message(msg.Channel, "Ordine cancellato")
	})

	t.bot.RespondTo("^(?i)email$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		order := getOrder(t.brain)
		subj := "Ordine Develer del giorno " + order.Timestamp.Format("02/01/2006")
		body := order.Format(false, false)

		out := subj + "\n" + body + "\n\n" +
			"<mailto:info@tuttobene-bar.it,sara@tuttobene-bar.it" +
			"?subject=" + url.PathEscape(subj) +
			"&body=" + url.PathEscape(body) +
			"|Link `mailto` clickabile>"
		t.bot.Message(msg.Channel, out)
	})

	t.bot.RespondTo("^(?i)menu([\\s\\S]*)?", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {

		showPrices := false

		if strings.TrimSpace(args[1]) == "price" {
			showPrices = true
		} else if args[1] != "" {
			t.bot.Message(msg.Channel, "Se stai cercando di impostare il menù, usa il comando `setmenu`\nPer vedere il menù corrente, usa il comando `menu` senza argomenti.")
			return
		}

		var m tuttobene.Menu
		err := t.brain.Get("menu", &m)
		if err == redis.Nil {
			t.bot.Message(msg.Channel, "Non c'è nessun menù impostato!")
		} else {
			t.bot.Message(msg.Channel, "Ecco il menù:\n"+m.Format(showPrices))
		}
	})

	t.bot.RespondTo("^(?i)setmenu([\\s\\S]*)?", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		if args[1] != "" {
			menu := strings.Split(strings.TrimSpace(sanitize(args[1])), "\n")
			m, err := tuttobene.ParseMenuCells(menu, []string{})
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

	t.bot.RespondTo("^(?i)remind(.*)$", t.Remind)

	t.bot.RespondTo("^(?i)segna(.*)$", t.Mark)

	t.bot.RespondTo("^(?i)rmorder (.*)$", func(b *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
		u := args[1]
		name := User{u, ""}
		finduser := t.bot.FindUser(u)
		if finduser != nil {
			name = User{finduser.Name, finduser.ID}
		}
		order := getOrder(t.brain)
		old := order.ClearUser(name)
		if old != "" {
			t.bot.Message(msg.Channel, fmt.Sprintf("Ok, cancello ordine di %s:\n%s", name.Name, old))
		} else {
			t.bot.Message(msg.Channel, fmt.Sprintf("%s non aveva ordinato nulla", name.Name))
		}

		order.Save(t.brain)
	})
}
