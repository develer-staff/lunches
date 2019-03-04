package grifts

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/develersrl/lunches/pkg/brain"
	"github.com/go-redis/redis"
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

		for i, s := range sched {
			r := strings.SplitN(s, ";", 2)
			if len(r) < 2 {
				log.Println("Malformed cron string: " + s)
				return nil
			}
			sch, err := cron.ParseStandard(r[0])
			if err != nil {
				log.Println(err)
				return nil
			}
			txt := strings.TrimSpace(r[1])
			loc, err := time.LoadLocation("Europe/Rome")
			if err != nil {
				log.Println("LoadLocation error: ", err)
				return nil
			}
			now := time.Now().In(loc)
			now = now.Add(-timerInterval / 2)
			next := sch.Next(now)

			if now.Add(timerInterval).Sub(next) > 0 {
				log.Printf("Executing cron #%d - %s", i, s)

				args := strings.Split(txt, " ")
				if len(args) < 1 {
					log.Println("No task specified!")
					return nil
				}
				task := "tinabot:" + args[0]
				ctx := NewContext(task)
				ctx.Args = args[1:]
				err := Run(task, ctx)
				if err != nil {
					log.Println(err)
				}
				return err

			}
		}
		return nil
	})

	Desc("post", "post on slack. Usage: post <channel> <message>")
	Add("post", func(c *Context) error {
		token := os.Getenv("SLACK_BOT_TOKEN")
		if token == "" {
			log.Fatalln("No slackbot token found!")
		}

		if len(c.Args) < 2 {
			log.Fatalln("Not enough arguments, usage: post <channel> <message>")
		}
		channel := c.Args[0]

		api := slack.New(token)
		api.PostMessage(channel, slack.MsgOptionText(strings.Join(c.Args[1:], " "), false))
		return nil
	})
})
