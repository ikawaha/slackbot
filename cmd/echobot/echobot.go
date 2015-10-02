package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ikawaha/slackbot"
)

// echo sample
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: bot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	bot, err := slackbot.New(os.Args[1])
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
		log.Printf("bot_id: %v, msguser_id: %v, msg:%+v\n", bot.ID, msg.UserID, msg)
		if bot.ID == msg.MentionID() && msg.Type == "message" && msg.SubType == "" {
			go func(m slackbot.Message) {
				m.Text = m.TextBody()
				bot.PostMessage(m)
			}(msg)
		}
	}
}
