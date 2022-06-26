package slackbot

import (
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/ikawaha/slackbot/socketmode"
	"github.com/ikawaha/slackbot/webapi"
)

// Client represents a slack client.
type Client struct {
	webAPIClient     *webapi.Client
	socketModeClient *socketmode.Client
	members          map[string]webapi.User // cache
}

// New creates a slack bot from app-level token and API token.
func New(appLevelToken, apiToken string, opts ...Option) (*Client, error) {
	var c config
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	a, err := webapi.New(apiToken, c.webAPIClientOptions...)
	if err != nil {
		return nil, err
	}
	s, err := socketmode.New(appLevelToken, c.socketModeClientOptions...)
	if err != nil {
		return nil, err
	}
	ret := Client{
		webAPIClient:     a,
		socketModeClient: s,
	}
	return &ret, nil
}

var (
	metaTag     = regexp.MustCompile(`<.*?>`)
	parentheses = strings.NewReplacer("&lt;", "<", "&gt;", ">")
)

// ReceiveMessage receives a message and passes it to a handler for processing.
func (c Client) ReceiveMessage(ctx context.Context, handler func(ctx context.Context, event *socketmode.Event) error) error {
	return c.socketModeClient.ReceiveMessage(ctx, handler)
}

// PostMessage sends a message to the Slack channel.
func (c Client) PostMessage(ctx context.Context, channelID, msg string) error {
	_, err := c.webAPIClient.PostMessage(ctx, channelID, msg)
	return err
}

// PlainMessageText resolves meta tags of the message text and return it.
func (c Client) PlainMessageText(msg string) string {
	txt := metaTag.ReplaceAllStringFunc(msg, func(s string) string {
		var id string
		for i := 0; i < len(s)-2; i++ {
			if s[i] == '@' {
				id = s[i+1 : len(s)-1]
				break
			}
		}
		if v, ok := c.members[id]; ok {
			return "@" + v.Name
		}
		if id != "" {
			return "@" + id
		}
		return s
	})
	return parentheses.Replace(txt)
}

// UploadImage uploads an image by files.upload API.
// see. https://api.slack.com/methods/files.upload
func (c Client) UploadImage(channels []string, title, fileName, fileType, comment string, img io.Reader) error {
	return c.webAPIClient.UploadImage(channels, title, fileName, fileType, comment, img)
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	return c.socketModeClient.Close()
}

// UsersList lists all users in a Slack team.
// see. https://api.slack.com/methods/users.list
func (c Client) UsersList() ([]webapi.User, error) {
	return c.webAPIClient.UsersList()
}

// Users lists all users in a Slack team and returns it's userID map.
func (c Client) Users() (map[string]webapi.User, error) {
	return c.webAPIClient.Users()
}

// RefreshUsersCache updates the client's cached user map.
func (c *Client) RefreshUsersCache() error {
	return c.webAPIClient.RefreshUsersCache()
}

// User returns the user corresponding to user ID from the client's user cache.
func (c *Client) User(id string) (webapi.User, bool) {
	return c.webAPIClient.User(id)
}
