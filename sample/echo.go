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
