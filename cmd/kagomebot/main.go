package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ikawaha/kagome/tokenizer"

	"github.com/ikawaha/slackbot"
)

const (
	botResponseSleepTime = 3 * time.Second
)

func init() {
	_ = tokenizer.SysDic()
}

func tokenize(sen string) string {
	t := tokenizer.New()
	tokens := t.Tokenize(sen)
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "```")
	for i := 1; i < len(tokens); i++ {
		if tokens[i].Class == tokenizer.DUMMY {
			fmt.Fprintf(&buf, "%s\n", tokens[i].Surface)
			continue
		}
		features := strings.Join(tokens[i].Features(), ",")
		fmt.Fprintf(&buf, "%s\t%v\n", tokens[i].Surface, features)
	}
	fmt.Fprintln(&buf, "```")
	return buf.String()
}

// ja-morph sample
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
				m.Text = tokenize(m.TextBody())
				if e := bot.PostMessage(m); e != nil {
					log.Printf("post error, %v", e)
				}
			}(msg)
		}
	}
}
