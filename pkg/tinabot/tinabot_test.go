package tinabot

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/develersrl/lunches/pkg/slackbot"
	"github.com/nlopes/slack"
)

func TestSplitSep(t *testing.T) {

	tests := map[string][]string{
		"a&b":              {"a", "b"},
		"a&b&c&d":          {"a", "b", "c", "d"},
		"a\\&b&c\\&d":      {"a&b", "c&d"},
		"a\\&b\\&c\\&d":    {"a&b&c&d"},
		"abcd":             {"abcd"},
		"&ab&cd":           {"", "ab", "cd"},
		"&ab&cd&":          {"", "ab", "cd", ""},
		"&ab&&&cd&":        {"", "ab", "", "", "cd", ""},
		"&ab&\\&&cd&":      {"", "ab", "&", "cd", ""},
		"&ab&\\&\\\\q&cd&": {"", "ab", "&\\q", "cd", ""},
		"a\\\\&b":          {"a\\&b"},
	}

	for i := range tests {
		out := splitEsc(i, "&")
		for j := range out {
			if out[j] != tests[i][j] {
				t.Fatalf("Error, wanted %v, got %v", tests[i], out)
			}
		}

	}
}

func getTestUserProfile() slack.UserProfile {
	return slack.UserProfile{
		StatusText:            "testStatus",
		StatusEmoji:           ":construction:",
		RealName:              "Test Real Name",
		RealNameNormalized:    "Test Real Name Normalized",
		DisplayName:           "Test Display Name",
		DisplayNameNormalized: "Test Display Name Normalized",
		Email:                 "test@develer.com",
		Image24:               "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_24.jpg",
		Image32:               "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_32.jpg",
		Image48:               "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_48.jpg",
		Image72:               "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_72.jpg",
		Image192:              "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-10-18/92962080834_ef14c1469fc0741caea1_192.jpg",
		Fields:                slack.UserProfileCustomFields{},
	}
}

func NewMockSlackUser(username string) *slack.User {
	return &slack.User{
		ID:                "UXXXXXXXX",
		Name:              username,
		Deleted:           false,
		Color:             "9f69e7",
		RealName:          "testuser",
		TZ:                "America/Los_Angeles",
		TZLabel:           "Pacific Daylight Time",
		TZOffset:          -25200,
		Profile:           getTestUserProfile(),
		IsBot:             false,
		IsAdmin:           false,
		IsOwner:           false,
		IsPrimaryOwner:    false,
		IsRestricted:      false,
		IsUltraRestricted: false,
		Has2FA:            false,
	}
}

type MockBrain struct{}

func (b MockBrain) Set(key string, val interface{}) error {
	return nil
}
func (b MockBrain) Read(key string) (string, error) {
	return "", nil
}

func (b MockBrain) Get(key string, q interface{}) error {
	return nil
}

func (b MockBrain) Close() error {
	return nil
}

type SentMessage struct {
	channel string
	message string
}

type MockBot struct {
	actions  map[*regexp.Regexp]slackbot.Action
	defact   slackbot.SimpleAction
	messages []SentMessage
}

func NewMockBot() *MockBot {
	return &MockBot{
		actions:  make(map[*regexp.Regexp]slackbot.Action),
		messages: make([]SentMessage, 0, 1),
	}
}

func (b *MockBot) RespondTo(match string, action slackbot.Action) {
	b.actions[regexp.MustCompile(match)] = action
}
func (b *MockBot) DefaultResponse(action slackbot.SimpleAction) {
	b.defact = action
}
func (b *MockBot) Message(channel string, msg string) {
	b.messages = append(b.messages, SentMessage{channel, msg})
}
func (b *MockBot) HandleMsg(channel, username, text string) {
	msg := &slackbot.BotMsg{channel, username, text}
	txt := strings.TrimLeft(strings.TrimSpace(text), "<@"+username+"> ")
	user := NewMockSlackUser(username)
	for match, action := range b.actions {
		if matches := match.FindAllStringSubmatch(txt, -1); matches != nil {
			action(b, msg, user, matches[0]...)
			return
		}
	}

	if b.defact != nil {
		b.defact(b, msg, user)
	}
}
func (b *MockBot) FindUser(user string) *slack.User {
	return nil
}
func (b *MockBot) OpenDirectChannel(user string) (string, error) {
	return "", nil
}

func TestCommands(t *testing.T) {
	// setup
	bot := NewMockBot()
	tina := New(bot, &MockBrain{})
	tina.AddCommands()
	url := "http://localhost:9998"
	os.Setenv("MARK_URL", url)
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "OK")
	})
	go http.ListenAndServe(url, nil)

	// action
	tina.bot.HandleMsg("chat_cibo", "pippo", "segna PD")
	// we have to wait because the Mark action happens in another goroutine
	// is there a better way to handle this case?
	time.Sleep(300 * time.Millisecond)

	// assert
	if len(bot.messages) != 1 {
		t.Error("Received no messages")
	}
	if strings.Index(bot.messages[0].message, "Ok") == -1 {
		t.Errorf("Received error: %v", bot.messages[0].message)
	}
}
