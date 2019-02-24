package grifts

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/develersrl/lunches/actions/brain"
	"github.com/go-redis/redis"
	. "github.com/markbates/grift/grift"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
)

var _ = Namespace("tinabot", func() {

	Desc("Cron", "Execute scheduled tasks")
	Add("cron", func(c *Context) error {
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			log.Fatalln("No redis URL found!")
		}
		token := os.Getenv("SLACK_BOT_TOKEN")
		if token == "" {
			log.Fatalln("No slackbot token found!")
		}

		timerInterval := 10 * time.Minute
		interval := os.Getenv("INTERVAL_MINUTES")
		if interval != "" {
			n, err := strconv.Atoi(interval)
			if err != nil {
				timerInterval = time.Duration(n) * time.Minute
			}
		}

		foodChannel := os.Getenv("FOOD_CHANNEL")
		if foodChannel == "" {
			log.Fatalln("No food channel set!")
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
			txt := r[1]
			loc, err := time.LoadLocation("CET")
			now := time.Now().In(loc)
			log.Println("Time now:", now)
			next := sch.Next(now)

			if now.Add(timerInterval).Sub(next) > 0 {
				log.Printf("Executing cron #%d - %s", i, s)
				api := slack.New(token)
				api.PostMessage(foodChannel, slack.MsgOptionText(txt, false))
			}

		}
		return nil
	})
})
