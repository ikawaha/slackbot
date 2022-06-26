package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ikawaha/slackbot"
	"github.com/ikawaha/slackbot/socketmode"
)

// Bot your bot
type Bot struct {
	*slackbot.Client
}

func NewBot(appToken, botToken string) (*Bot, error) {
	c, err := slackbot.New(appToken, botToken, slackbot.CacheUsers(), slackbot.Debug())
	if err != nil {
		return nil, err
	}
	return &Bot{Client: c}, err
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: bot app-level-token slack-bot-token\n")
		os.Exit(1)
	}

	bot, err := NewBot(os.Args[1], os.Args[2]) // set your app-level-token and bot token!
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()
	fmt.Println("^C exits")

	for {
		if err := bot.ReceiveMessage(context.TODO(), func(ctx context.Context, ev *socketmode.Event) error {
			u, ok := bot.User(ev.UserID)
			if !ok || u.IsBot {
				return nil
			}
			msg := "Hello, " + u.Name + ": " + ev.Text
			if err := bot.PostMessage(ctx, ev.Channel, msg); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Printf("%v", err)
		}
	}
}
