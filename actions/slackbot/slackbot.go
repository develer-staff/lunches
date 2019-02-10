package slackbot

import (
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack/slackevents"

	"github.com/nlopes/slack"
)

type SimpleAction func(*Bot, *slackevents.MessageEvent, *slack.User)
type Action func(*Bot, *slackevents.MessageEvent, *slack.User, ...string)

type Bot struct {
	UserID string

	client *slack.Client

	actions map[*regexp.Regexp]Action
	defact  SimpleAction
}

func New(botID string, api *slack.Client) *Bot {

	bot := &Bot{
		UserID:  botID,
		client:  api,
		actions: make(map[*regexp.Regexp]Action),
	}

	return bot
}

func (bot *Bot) RespondTo(match string, action Action) {
	bot.actions[regexp.MustCompile(match)] = action
}

func (bot *Bot) DefaultResponse(action SimpleAction) {
	bot.defact = action
}

func (bot *Bot) Message(channel string, msg string) {
	bot.client.PostMessage(channel, slack.MsgOptionText(msg, false))
}

func (bot *Bot) validMessage(msg *slackevents.MessageEvent) bool {
	return msg.Type == "message" &&
		(msg.SubType != "message_deleted" && msg.SubType != "bot_message") &&
		msg.User != bot.UserID &&
		(strings.HasPrefix(msg.Text, "<@"+bot.UserID+">") || strings.HasPrefix(msg.Channel, "D"))
}

func (bot *Bot) cleanupMsg(msg string) string {
	return strings.TrimLeft(strings.TrimSpace(msg), "<@"+bot.UserID+"> ")
}

func (bot *Bot) HandleMsg(msg *slackevents.MessageEvent) {
	if !bot.validMessage(msg) {
		return
	}

	txt := bot.cleanupMsg(msg.Text)

	user, err := bot.client.GetUserInfo(msg.User)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for match, action := range bot.actions {
		if matches := match.FindAllStringSubmatch(txt, -1); matches != nil {
			action(bot, msg, user, matches[0]...)
			return
		}
	}

	if bot.defact != nil {
		bot.defact(bot, msg, user)
	}
}
