package tinabot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/develersrl/lunches/pkg/slackbot"
	"github.com/nlopes/slack"
)

func Mark(user, food string) error {
	markURL := os.Getenv("MARK_URL")
	if markURL == "" {
		return errors.New("no mark URL found")
	}
	url := strings.Replace(markURL, "<USER>", user, -1)
	url = strings.Replace(url, "<FOOD>", food, -1)

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	var err error
	for i := 0; i < 3; i++ {
		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			break
		} else {
			log.Printf("ERROR marking user %s: %s, retrying\n", user, err.Error())
		}
	}
	return err
}

func MarkUser(user *slack.User, food string) error {
	mail := user.Profile.Email
	if strings.Contains(mail, "@develer.com") {
		nick := strings.TrimSuffix(mail, "@develer.com")
		return Mark(nick, food)
	}
	return errors.New("user does not have a Develer mail")
}

func (t *TinaBot) Mark(bot *slackbot.Bot, msg *slackbot.BotMsg, user *slack.User, args ...string) {
	food := strings.TrimSpace(args[1])

	validFood := []string{
		"P",
		"PS",
		"PD",
		"S",
		"SD",
		"D",
		"PSD",
		"Niente",
	}

	if food == "" {
		t.bot.Message(msg.Channel, "Cosa devo segnare? Consulta `aiuto` per sapere come fare")
		return
	}

	for _, f := range validFood {
		if strings.ToUpper(f) == strings.ToUpper(food) {
			// This can be slow so spawn a goroutine to give Slack a fast reply and avoid retrys
			go func() {
				err := MarkUser(user, f)
				if err != nil {
					t.bot.Message(msg.Channel, "errore: "+err.Error())
					return
				}
				t.bot.Message(msg.Channel, fmt.Sprintf("Ok, segnato '%s' per %s sul foglio dei pranzi", f, user.Name))
			}()
			return
		}
	}
	t.bot.Message(msg.Channel, fmt.Sprintf("Scusami, la stringa '%s' non è valida.\nStringhe valide sono: %s", food, strings.Join(validFood, ", ")))
}
