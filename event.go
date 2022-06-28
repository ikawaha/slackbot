package slackbot

import (
	"github.com/ikawaha/slackbot/socketmode"
)

type EventType = socketmode.EventType

const (
	// AppMention is a Slack event type.
	// Subscribe to only the message events that mention your app or bot.
	AppMention = socketmode.AppMention

	// Message is a Slack event type.
	// A message was sent to a channel.
	Message = socketmode.Message

	// SlashCommand is a slash command.
	SlashCommand = socketmode.SlashCommand
)
