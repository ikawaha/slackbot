package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	id := "<@" + bot.ID + ">"
	for {
		m, err := bot.GetMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v\n", m)
		if m.Type == "message" && strings.HasPrefix(m.Text, id) {
			msg := strings.TrimPrefix(m.Text, id)
			if strings.HasPrefix(msg, ": ") {
				msg = strings.TrimPrefix(msg, ": ")
			}
			m.Text = msg
			go bot.PostMessage(m)
		}
	}
}
