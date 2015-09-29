package main

import (
	"fmt"
	"log"
	"os"

	slackbot "github.com/ikawaha/slackbot"
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
		panic(err)
	}
	fmt.Println("^C exits")

	for {
		m, err := bot.GetMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("msg:%+v\n", m)
		if bot.ID == m.UserID() && m.Type == "message" {
			m.Text = m.TextBody()
			go func(msg slackbot.Message) {
				bot.PostMessage(msg)
			}(m)
		}
	}
}
