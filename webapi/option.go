package webapi

import (
	"context"
)

// Option represents the client's option.
type Option func(*Client) error

// CacheUsers lists all users in a Slack team and caches it.
// required scopes: `users:read`
func CacheUsers() Option {
	return func(c *Client) error {
		return c.RefreshUsersCache(context.TODO())
	}
}

// Debug is the debug option.
func Debug() Option {
	return func(c *Client) error {
		c.debug = true
		return nil
	}
}
