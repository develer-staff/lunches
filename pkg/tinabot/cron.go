package tinabot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nlopes/slack"

	"github.com/develersrl/lunches/pkg/slackbot"
	"github.com/robfig/cron"

	"github.com/go-redis/redis"
)

func (t *TinaBot) Cron(bot slackbot.BotInterface, msg *slackbot.BotMsg, user *slack.User, args ...string) {

	var crontab []string
	if args[1] != "" {
		crontab = strings.SplitN(strings.TrimSpace(args[1]), " ", 2)
	} else {
		crontab = nil
	}

	if crontab == nil {
		var sched []string
		err := t.brain.Get("cron", &sched)
		if err == redis.Nil || len(sched) == 0 {
			bot.Message(msg.Channel, "Non c'è nessun cron impostato")
		} else {
			reply := "```"
			for i, s := range sched {
				reply += fmt.Sprintf("%d - %s\n", i, s)
			}
			bot.Message(msg.Channel, reply+"```")
		}
	} else {
		if len(crontab) < 2 {
			bot.Message(msg.Channel, "Argomenti insufficienti!")
			return
		}

		switch strings.ToLower(crontab[0]) {
		case "add":
			cmd := strings.Split(crontab[1], ";")
			if len(cmd) < 2 {
				bot.Message(msg.Channel, "Argomenti insufficienti!")
				return
			}
			_, err := cron.ParseStandard(cmd[0])
			if err != nil {
				bot.Message(msg.Channel, "Errore di formato cron: "+err.Error())
				return
			}
			var sched []string
			t.brain.Get("cron", &sched)
			sched = append(sched, crontab[1])
			t.brain.Set("cron", &sched)
			bot.Message(msg.Channel, fmt.Sprintf("Ok, cron aggiunto:```%d - %s```", len(sched)-1, crontab[1]))
		case "rm":
			var sched []string
			err := t.brain.Get("cron", &sched)
			if err != nil || len(sched) == 0 {
				bot.Message(msg.Channel, "Nessun cron impostato!")
				return
			}
			n, err := strconv.Atoi(crontab[1])
			if err != nil {
				bot.Message(msg.Channel, "Errore di parsing indice: "+err.Error())
				return
			}
			if n > len(sched)-1 || n < 0 {
				bot.Message(msg.Channel, "Indice inesistente!")
				return
			}
			bot.Message(msg.Channel, fmt.Sprintf("Ok, rimosso cron ```%d - %s```", n, sched[n]))
			sched = append(sched[:n], sched[n+1:]...)
			t.brain.Set("cron", sched)
		}
	}
}
