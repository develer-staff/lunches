package actions

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/develersrl/lunches/actions/brain"
	"github.com/develersrl/lunches/actions/slackbot"
	"github.com/develersrl/lunches/actions/tinabot"
	"github.com/gobuffalo/buffalo"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// SlackHandler default implementation.
func SlackHandler(c buffalo.Context) error {
	//return c.Render(200, r.HTML("slack/handler.html"))
	slackToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackToken == "" {
		log.Fatalln("No SLACK_BOT_TOKEN found!")
	}
	accessToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
	if accessToken == "" {
		log.Fatalln("No SLACK_VERIFICATION_TOKEN found!")
	}
	botID := os.Getenv("BOT_ID")
	if accessToken == "" {
		log.Fatalln("No BOT_ID found!")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatalln("No redis URL found!")
	}

	api := slack.New(slackToken)

	brain := brain.New(redisURL)
	defer brain.Close()

	bot := slackbot.New(botID, api)
	tinabot.Tinabot(bot, brain)

	w := c.Response()
	r := c.Request()
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: accessToken}))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent

		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			bot.HandleMsg(ev.Channel, ev.User, ev.Text)
		case *slackevents.MessageEvent:
			bot.HandleMsg(ev.Channel, ev.User, ev.Text)
		}

	}

	return nil
}
