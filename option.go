package slackbot

import (
	"github.com/ikawaha/slackbot/socketmode"
	"github.com/ikawaha/slackbot/webapi"
)

type config struct {
	webAPIClientOptions     []webapi.Option
	socketModeClientOptions []socketmode.Option
}

func (c *config) AddWebAPIOption(o webapi.Option) {
	c.webAPIClientOptions = append(c.webAPIClientOptions, o)
}

func (c *config) AddSocketModeOption(o socketmode.Option) {
	c.socketModeClientOptions = append(c.socketModeClientOptions, o)
}

// Option represents the client's option.
type Option func(*config) error

// CacheUsers lists all users in a Slack team and caches it.
func CacheUsers() Option {
	return func(c *config) error {
		c.AddWebAPIOption(webapi.CacheUsers())
		return nil
	}
}

// Debug is the debug option.
func Debug() Option {
	return func(c *config) error {
		c.AddWebAPIOption(webapi.Debug())
		c.AddSocketModeOption(socketmode.Debug())
		return nil
	}
}
