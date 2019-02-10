package actions

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/develersrl/lunches/actions/slackbot"
	"github.com/gobuffalo/buffalo"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// SlackHandler default implementation.
func SlackHandler(c buffalo.Context) error {
	//return c.Render(200, r.HTML("slack/handler.html"))
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	accessToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
	botID := os.Getenv("BOT_ID")

	bot := slackbot.New(botID, api)
	bot.DefaultResponse(func(b *slackbot.Bot, msg *slackevents.MessageEvent, user *slack.User) {
		bot.Message(msg.Channel, "Mi dispiace "+user.Name+", purtroppo non posso farlo.")
	})

	bot.RespondTo("^(?i)menu([\\s\\S]*)?", func(b *slackbot.Bot, msg *slackevents.MessageEvent, user *slack.User, args ...string) {
		bot.Message(msg.Channel, "Non c'è nessun menu impostato!")
	})

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
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		case *slackevents.MessageEvent:
			log.Println(ev)
			bot.HandleMsg(ev)
		}

	}

	return nil
}
