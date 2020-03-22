package slackbot

import (
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type BotMsg struct {
	Channel string
	User    string
	Text    string
}

type SimpleAction func(BotInterface, *BotMsg, *slack.User)
type Action func(BotInterface, *BotMsg, *slack.User, ...string)

type BotInterface interface {
	RespondTo(match string, action Action)
	DefaultResponse(action SimpleAction)
	Message(channel string, msg string)
	HandleMsg(channel, username, text string)
	FindUser(user string) *slack.User
	OpenDirectChannel(user string) (string, error)
}

type Bot struct {
	UserID string

	Client *slack.Client

	actions map[*regexp.Regexp]Action
	defact  SimpleAction
}

func New(botID string, api *slack.Client) *Bot {

	bot := &Bot{
		UserID:  botID,
		Client:  api,
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
	bot.Client.PostMessage(channel, slack.MsgOptionText(msg, false))
}

func (bot *Bot) validMessage(msg *BotMsg) bool {
	return msg.User != bot.UserID &&
		(strings.HasPrefix(msg.Text, "<@"+bot.UserID+">") || strings.HasPrefix(msg.Channel, "D"))
}

func (bot *Bot) cleanupMsg(msg string) string {
	return strings.TrimLeft(strings.TrimSpace(msg), "<@"+bot.UserID+"> ")
}

func (bot *Bot) HandleMsg(channel, username, text string) {
	msg := &BotMsg{channel, username, text}
	if !bot.validMessage(msg) {
		return
	}

	txt := bot.cleanupMsg(msg.Text)

	user, err := bot.Client.GetUserInfo(msg.User)
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

func (bot *Bot) FindUser(user string) *slack.User {
	if strings.HasPrefix(user, "<@") {
		user = strings.Trim(user, "<@>")
		u, err := bot.Client.GetUserInfo(user)
		if err != nil {
			log.Println(err)
			return nil
		}
		return u
	}

	users, err := bot.Client.GetUsers()
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, u := range users {
		if strings.ToLower(u.Name) == strings.ToLower(user) {
			return &u
		}
	}
	return nil
}

func (bot *Bot) OpenDirectChannel(user string) (string, error) {
	_, _, ch, err := bot.Client.OpenIMChannel(user)
	if err != nil {
		return "", err
	}
	return ch, nil
}
