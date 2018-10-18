package main

import (
	"fmt"
	"github.com/ikawaha/slackbot"
	"log"
	"os"
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
		msg, err := bot.GetMessage()
		if err != nil {
			log.Printf("receive error, %v", err)
		}
		if bot.ID == msg.MentionID() && msg.Type == "message" && msg.SubType == "" {
			go func(m slackbot.Message) {
				m.Text = m.TextBody()
				bot.PostMessage(m)
			}(msg)
		}
	}
}
