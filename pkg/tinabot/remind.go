package tinabot

import (
	"strings"

	"github.com/go-redis/redis"
	"github.com/nlopes/slack"

	"github.com/develersrl/lunches/pkg/slackbot"
)

func formatReminder(mask int) string {
	weekNames := []string{
		"domenica",
		"lunedì",
		"martedì",
		"mercoledì",
		"giovedì",
		"venerdì",
		"sabato",
	}

	if mask == 0 {
		return "Reminder disattivato"
	}

	reply := "Reminder attivo "
	if mask == 0xff {
		reply += "tutti i giorni"
	} else {
		var days []string
		for i := uint(0); i < 7; i++ {
			if ((1 << i) & mask) != 0 {
				days = append(days, weekNames[i])
			}
		}
		reply += strings.Join(days, ", ")
	}
	return reply
}

func (t *TinaBot) Remind(bot *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
	weekMask := map[string]int{
		"off": 0,
		"dis": 0,
		"fal": 0,
		"0":   0,
		"on":  0xff,
		"ena": 0xff,
		"tru": 0xff,
		"1":   0xff,
		"sem": 0xff,
		"tut": 0xff,
		"all": 0xff,

		"dom": 1 << 0,
		"lun": 1 << 1,
		"mar": 1 << 2,
		"mer": 1 << 3,
		"gio": 1 << 4,
		"ven": 1 << 5,
		"sab": 1 << 6,
	}

	if args[1] == "" {
		var remind map[string]int
		err := t.brain.Get("remind", &remind)
		if err == redis.Nil || len(remind) == 0 {
			bot.Message(msg.Channel, "Non c'è nessun reminder impostato")
		} else {
			if val, ok := remind[user.Name]; ok {
				bot.Message(msg.Channel, formatReminder(val))
			} else {
				bot.Message(msg.Channel, "Non c'è nessun reminder impostato")
			}
		}

	} else {
		arg := strings.ToLower(args[1])
		mask := 0

		cmdFound := false
		for _, d := range strings.Split(arg, ",") {
			d = strings.TrimSpace(d)
			if len(d) > 3 {
				d = d[:3]
			}

			if m, ok := weekMask[d]; ok {
				mask |= m
				cmdFound = true
			}
		}

		if !cmdFound {
			bot.Message(msg.Channel, "Mi spiace, ma non ho capito cosa mi stai chiedendo di ricordare")
			return
		}

		remind := make(map[string]int)
		t.brain.Get("remind", &remind)
		if mask == 0 {
			delete(remind, user.ID)
		} else {
			remind[user.ID] = mask
		}
		t.brain.Set("remind", remind)

		bot.Message(msg.Channel, formatReminder(mask))
	}
}
