package actions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// SlackHandler default implementation.
func SlackHandler(c buffalo.Context) error {
	//return c.Render(200, r.HTML("slack/handler.html"))
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	accessToken := os.Getenv("SLACK_OAUTH_TOKEN")

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
		}
	}

	return nil
}
