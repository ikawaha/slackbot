package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-haiku" // awesome!

	"github.com/ikawaha/slackbot"
)

const (
	botResponseSleepTime = 3 * time.Second
)

// haiku sample
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

	r575 := []int{5, 7, 5}
	for {
		msg, err := bot.GetMessage()
		if err != nil {
			log.Printf("receive error, %v", err)
		}
		log.Printf("bot_id: %v, msguser_id: %v, msg:%+v\n", bot.ID, msg.UserID, msg)
		if msg.Type != "message" {
			continue
		}
		go func(m slackbot.Message) {
			t := m.TextBody()
			hs := haiku.Find(t, r575)
			if len(hs) < 1 {
				return
			}
			var tmp []string
			for _, h := range hs {
				tmp = append(tmp, fmt.Sprintf("```%v```", h))
			}
			m.Text = strings.Join(tmp, "\n")
			m.Text += " 575だ"
			time.Sleep(botResponseSleepTime)
			bot.PostMessage(m)
		}(msg)
	}
}
