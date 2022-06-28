package slackbot

import (
	"github.com/ikawaha/slackbot/socketmode"
	"github.com/ikawaha/slackbot/webapi"
)

type config struct {
	webAPIClientOptions     []webapi.Option
	socketModeClientOptions []socketmode.Option
	searchBotID             bool
	botName                 string
}

// AddWebAPIOption adds an option to the Web API client.
func (c *config) AddWebAPIOption(o webapi.Option) {
	c.webAPIClientOptions = append(c.webAPIClientOptions, o)
}

// AddSocketModeOption adds an option to the Socket Mode client.
func (c *config) AddSocketModeOption(o socketmode.Option) {
	c.socketModeClientOptions = append(c.socketModeClientOptions, o)
}

// Option represents the client's option.
type Option func(*config) error

// CacheUsers lists all users in a Slack team and caches it.
// required scopes: `users:read`
func CacheUsers() Option {
	return func(c *config) error {
		c.AddWebAPIOption(webapi.CacheUsers())
		return nil
	}
}

// SetBotID sets bot ID and name to a client. When this option is specified,
// all user data in a Slack team will be cached even if the CacheUsers option is not set.
// required scopes: `users:read`
func SetBotID(name string) Option {
	return func(c *config) error {
		c.searchBotID = true
		c.botName = name
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
