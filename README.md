[![Go Reference](https://pkg.go.dev/badge/github.com/ikawaha/slackbot.svg)](https://pkg.go.dev/github.com/ikawaha/slackbot)

# slackbot

Tiny slack bot client.

# Prepare

Please get the token for your slack bot.

`https://[YOURTEAM].slack.com/services/new/bot`

see. https://api.slack.com/

# Interface

Echo bot sample (see. `sample/echo.go`).
See also [kagome-bot](https://github.com/ikawaha/kagome-bot).
```Go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ikawaha/slackbot"
)

// your bot
type Bot struct {
	*slackbot.Client
}

func NewBot(token string) (*Bot, error) {
	c, err := slackbot.New(token)
	if err != nil {
		return nil, err
	}
	return &Bot{Client: c}, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: bot slack-bot-token\n")
		os.Exit(1)
	}

	bot, err := NewBot(os.Args[1]) // set your bot token!
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()
	fmt.Println("^C exits")

	for {
		msg, err := bot.ReceiveMessage(context.TODO())
		if err != nil {
			log.Printf("receive error, %v", err)
		}
		if strings.Contains(msg.Text, bot.ID) && msg.Type == "message" && msg.SubType == "" {
			go func(m slackbot.Message) {
				log.Print(m.Text)
				if err := bot.PostMessage(m); err != nil {
					log.Printf("post message failed: %v", err)
				}
			}(msg)
		}
	}
}


```

# Lisence

MIT
