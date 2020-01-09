package actions

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/develersrl/lunches/pkg/brain"
	"github.com/develersrl/lunches/pkg/tuttobene"
	"github.com/gobuffalo/buffalo"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/nlopes/slack"
)

// EmailHandler default implementation.
func EmailHandler(c buffalo.Context) error {
	log.Println("Email received!")
	domain := os.Getenv("MAILGUN_DOMAIN")
	if domain == "" {
		log.Println("MAILGUN_DOMAIN not set")
		return nil
	}

	apiKey := os.Getenv("MAILGUN_API_KEY")
	if apiKey == "" {
		log.Println("MAILGUN_API_KEY not set")
		return nil
	}

	mg := mailgun.NewMailgun(domain, apiKey)

	verified, err := mg.VerifyWebhookSignature(mailgun.Signature{
		TimeStamp: c.Param("timestamp"),
		Token:     c.Param("token"),
		Signature: c.Param("signature"),
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	if !verified {
		log.Printf("Mailgun signature verification error")
		return nil
	}

	if !strings.HasPrefix(c.Param("Content-Type"), "multipart/mixed") {
		log.Printf("Wrong POST Content-Type: '%s'", c.Param("Content-Type"))
		return nil
	}
	if c.Param("attachment-count") == "" {
		log.Println("No attachment found")
		return nil
	}

	n, err := strconv.Atoi(c.Param("attachment-count"))
	if err != nil {
		log.Println(err)
		return nil
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		log.Fatalln("No slackbot token found!")
	}

	channel := os.Getenv("FOOD_CHANNEL")
	if channel == "" {
		log.Println("No channel found!")
		return nil
	}
	api := slack.New(token)

	for i := 0; i < n; i++ {
		f, h, err := c.Request().FormFile(fmt.Sprintf("attachment-%d", i+1))
		if err != nil {
			log.Println(err)
			return nil
		}
		name := strings.ToLower(h.Filename)
		if strings.Contains(name, ".xlsx") {
			if h.Size > 500000 {
				log.Println("Attachemnt too large!")
				api.PostMessage(channel, slack.MsgOptionText("Menu ricevuto, file in attachment di dimensioni eccessive!", false))
				return nil
			}
			buf := make([]byte, h.Size)

			_, err := f.Read(buf)
			if err != nil {
				log.Println(err)
				return nil
			}

			m, err := tuttobene.ParseMenuBytes(buf)

			if err != nil {
				log.Println("Menu parse error: ", err)
				api.PostMessage(channel, slack.MsgOptionText("Menu ricevuto, errore durante l'analisi: "+err.Error(), false))
				return nil
			}
			redisURL := os.Getenv("REDIS_URL")
			if redisURL == "" {
				log.Println("No redis URL found!")
				return nil
			}

			b := brain.New(redisURL)
			defer b.Close()

			b.Set("menu", *m)

			log.Println("Tuttobene menu parsed correctly")

			date := m.Date.Format("02/01/2006")
			api.PostMessage(channel, slack.MsgOptionText("Ho appena ricevuto e impostato correttamente il menu per il giorno "+date, false))
			return nil
		}

		log.Println("Unrecognized attachment:", h.Filename)
	}

	log.Println("No menu parsed from email")
	return nil
}
