package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ikawaha/slackbot"
)

// Bot your bot
type Bot struct {
	*slackbot.Client
}

// NewBot creates a Slack bot.
func NewBot(appToken, botToken, botName string) (*Bot, error) {
	c, err := slackbot.New(appToken, botToken, botName)
	if err != nil {
		return nil, err
	}
	return &Bot{Client: c}, err
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: bot app-level-token slack-bot-token bot-name\n")
		os.Exit(1)
	}
	// set your app-level-token, bot token and bot name!
	bot, err := NewBot(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Close()
	fmt.Println("^C exits")

	callPrefix := "<@" + bot.ID + ">"
	for {
		if err := bot.ReceiveMessage(context.TODO(), func(ctx context.Context, e *slackbot.Event) error {
			log.Printf("-->%#+v", e)
			if !e.IsMessage() {
				return nil
			}
			u, ok := bot.User(e.UserID)
			if !ok || u.IsBot {
				return nil
			}
			if !strings.HasPrefix(e.Text, callPrefix) {
				return nil
			}
			msg := "Hi, " + u.Name + ": " + strings.TrimPrefix(e.Text, callPrefix)
			if err := bot.PostMessage(ctx, e.Channel, msg); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Printf("%v", err)
		}
	}
}
