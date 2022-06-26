package slackbot

import (
	"github.com/ikawaha/slackbot/socketmode"
	"github.com/ikawaha/slackbot/webapi"
)

type config struct {
	webapiClientOptions     []webapi.Option
	socketmodeClientOptions []socketmode.Option
}

func (c *config) AddWebapiOption(o webapi.Option) {
	c.webapiClientOptions = append(c.webapiClientOptions, o)
}

func (c *config) AddSocketmodeOption(o socketmode.Option) {
	c.socketmodeClientOptions = append(c.socketmodeClientOptions, o)
}

// Option represents the client's option.
type Option func(*config) error

// CacheUsers lists all users in a Slack team and caches it.
func CacheUsers() Option {
	return func(c *config) error {
		c.AddWebapiOption(webapi.CacheUsers())
		return nil
	}
}

func Debug() Option {
	return func(c *config) error {
		c.AddWebapiOption(webapi.Debug())
		c.AddSocketmodeOption(socketmode.Debug())
		return nil
	}
}
