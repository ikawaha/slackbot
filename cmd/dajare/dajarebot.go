package main

import (
	"fmt"
	"log"
	"os"
	"time"

	//"github.com/kurehajime/dajarep" //awesome!
	"github.com/ikawaha/dajarep"
	"github.com/ikawaha/slackbot"
)

const (
	botResponseSleepTime = 3 * time.Second
	responseTemplate     = "ねぇねぇ，%v\nいまの ```%v``` ってダジャレ？ダジャレ？"
)

// dajare sample
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
		if msg.Type != "message" || msg.SubType != "" {
			continue
		}
		go func(m slackbot.Message) {
			t := m.TextBody()
			dug, daj := dajarep.Dajarep(t)
			log.Printf("msg: %v, dajare: %+v, debug: %+v\n", t, daj, dug)
			if len(daj) < 1 {
				return
			}
			m.Text = fmt.Sprintf(responseTemplate, bot.UserName(m.UserID), t)
			time.Sleep(botResponseSleepTime)
			bot.PostMessage(m)
		}(msg)
	}
}
