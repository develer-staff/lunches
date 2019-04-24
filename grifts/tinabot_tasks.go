package grifts

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/develersrl/lunches/pkg/tuttobene"

	"github.com/develersrl/lunches/pkg/tinabot"

	"github.com/develersrl/lunches/pkg/brain"
	"github.com/go-redis/redis"
	"github.com/mailgun/mailgun-go/v3"
	. "github.com/markbates/grift/grift"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
)

var _ = Namespace("tinabot", func() {

	Desc("cron", "Execute scheduled tasks")
	Add("cron", func(c *Context) error {
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			log.Fatalln("No redis URL found!")
		}

		timerInterval := 10 * time.Minute
		interval := os.Getenv("INTERVAL_MINUTES")
		if interval != "" {
			n, err := strconv.Atoi(interval)
			if err != nil {
				timerInterval = time.Duration(n) * time.Minute
			}
		}

		brain := brain.New(redisURL)
		defer brain.Close()

		var sched []string
		err := brain.Get("cron", &sched)
		if err == redis.Nil || len(sched) == 0 {
			log.Println("No cron set")
			return nil
		}

		loc, err := time.LoadLocation("Europe/Rome")
		if err != nil {
			log.Println("LoadLocation error: ", err)
			return nil
		}

		for i, s := range sched {
			r := strings.SplitN(s, ";", 2)
			if len(r) < 2 {
				log.Println("Malformed cron string: " + s)
				continue
			}
			sch, err := cron.ParseStandard(r[0])
			if err != nil {
				log.Println(err)
				continue
			}
			txt := strings.TrimSpace(r[1])
			now := time.Now().In(loc)
			now = now.Add(-timerInterval / 2)
			next := sch.Next(now)

			if now.Add(timerInterval).Sub(next) > 0 {
				log.Printf("Executing cron #%d - %s", i, s)

				args := strings.Split(txt, " ")
				if len(args) < 1 {
					log.Println("No task specified!")
					continue
				}
				task := "tinabot:" + args[0]
				ctx := NewContext(task)
				ctx.Args = args[1:]
				err := Run(task, ctx)
				if err != nil {
					log.Println(err)
				}
			}
		}
		return nil
	})

	Desc("post", "post on slack. Usage: post <channel> [<options>] <message>")
	Add("post", func(c *Context) error {
		token := os.Getenv("SLACK_BOT_TOKEN")
		if token == "" {
			log.Fatalln("No slackbot token found!")
		}

		if len(c.Args) < 2 {
			log.Fatalln("Not enough arguments, usage: post <channel> [<options>] <message>")
		}
		channel := c.Args[0]
		onlyValidOrder := false
		onlyValidMenu := false
		startMsg := 1
		for i := 0; i < len(c.Args); i++ {
			opt := c.Args[i]
			if strings.HasPrefix(opt, "-") {
				startMsg = i + 1
				// Post only if there is a valid order
				if strings.Contains(opt, "o") {
					onlyValidOrder = true
				}

				// Post only if there is a valid menu
				if strings.Contains(opt, "m") {
					onlyValidMenu = true
				}
			}
		}

		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			log.Fatalln("No redis URL found!")
		}

		brain := brain.New(redisURL)
		defer brain.Close()

		var order tinabot.Order
		order.Load(brain)

		var menu tuttobene.Menu
		err := brain.Get("menu", &menu)
		if err == redis.Nil {
			log.Println("No menu found")
		}

		if onlyValidMenu && (err == redis.Nil || !menu.IsUpdated()) {
			return nil
		}

		if onlyValidOrder && !order.IsUpdated() {
			return nil
		}

		msg := strings.Join(c.Args[startMsg:], " ")
		msg = strings.Replace(msg, "$MENU", menu.String(), -1)
		msg = strings.Replace(msg, "$ORDER_NONAMES", order.Format(false), -1)
		msg = strings.Replace(msg, "$ORDER", order.Format(true), -1)
		msg = strings.Replace(msg, "\\n", "\n", -1)

		api := slack.New(token)
		api.PostMessage(channel, slack.MsgOptionText(msg, false))
		return nil
	})

	Desc("sendmail", "send the email of the lunch order to the given address(es)")
	Add("sendmail", func(c *Context) error {
		domain := os.Getenv("MAILGUN_DOMAIN")
		if domain == "" {
			log.Println("MAILGUN_DOMAIN not set")
			return nil
		}

		apiKey := os.Getenv("MAILGUN_API_KEY")
		if apiKey == "" {
			log.Println("MAILGUN_API_KEY not set")
			return nil
		}

		if len(c.Args) < 1 {
			log.Println("No recipients found!")
			return nil
		}

		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			log.Fatalln("No redis URL found!")
		}

		brain := brain.New(redisURL)
		defer brain.Close()

		var order tinabot.Order
		order.Load(brain)

		var menu tuttobene.Menu
		err := brain.Get("menu", &menu)
		if err == redis.Nil {
			log.Println("No menu found")
		}

		if !menu.IsUpdated() || !order.IsUpdated() {
			return nil
		}

		mg := mailgun.NewMailgun(domain, apiKey)
		var addresses []string
		for _, a := range c.Args {
			if strings.HasPrefix(a, "<mailto:") {
				a = strings.TrimPrefix(a, "<mailto:")
				a = strings.Split(a, "|")[0]
			}
			addresses = append(addresses, a)
		}
		to := strings.Join(addresses, ",")

		subj := "Ordine Develer del giorno " + order.Timestamp.Format("02/01/2006")
		from := "cibo@develer.com"
		body := order.Format(false)
		m := mg.NewMessage(from, subj, body, to)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		_, id, err := mg.Send(ctx, m)
		log.Println("Sendmail ID", id)
		return err
	})

	Desc("reminder", "send the users the reminder to order")
	Add("reminder", func(c *Context) error {
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			log.Fatalln("No redis URL found!")
		}

		brain := brain.New(redisURL)
		defer brain.Close()

		var remind map[string]tinabot.Remind
		err := brain.Get("remind", &remind)
		if err == redis.Nil || len(remind) == 0 {
			return nil
		}

		var order tinabot.Order
		order.Load(brain)

		var menu tuttobene.Menu
		err = brain.Get("menu", &menu)
		if err == redis.Nil {
			log.Println("No menu found")
			return nil
		}

		if !menu.IsUpdated() || !order.IsUpdated() {
			return nil
		}

		token := os.Getenv("SLACK_BOT_TOKEN")
		if token == "" {
			log.Fatalln("No slackbot token found!")
		}
		api := slack.New(token)

		loc, err := time.LoadLocation("Europe/Rome")
		if err != nil {
			log.Println("LoadLocation error: ", err)
			return nil
		}

		weekmask := 1 << uint(time.Now().In(loc).Weekday())

		fmtmsg := "Ciao %s, scusa il disturbo. Vedo che non hai ancora ordinato il pranzo e mi hai chiesto di ricordartelo. Ecco il menù di oggi:\n" + menu.String()
		for user, v := range remind {
			if v.Mask&weekmask != 0 {
				if _, ok := order.Users[user]; !ok {
					log.Printf("Sending reminder to %s\n", user)
					_, _, ch, err := api.OpenIMChannel(v.ID)
					if err != nil {
						log.Println(err)
						continue
					}

					txt := fmt.Sprintf(fmtmsg, user)
					api.PostMessage(ch, slack.MsgOptionText(txt, false))
				}
			}
		}

		return nil
	})
})
