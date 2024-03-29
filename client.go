package slackbot

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ikawaha/slackbot/socketmode"
	"github.com/ikawaha/slackbot/webapi"
)

// Client represents a slack client.
type Client struct {
	Name             string
	ID               string
	webAPIClient     *webapi.Client
	socketModeClient *socketmode.Client
}

type (
	// Event is an alias type of the socket mode event.
	Event = socketmode.Event

	// User is an alias type of the web api user.
	User = webapi.User
)

// New creates a slack bot from app-level token and API token.
func New(appLevelToken, apiToken string, opts ...Option) (*Client, error) {
	var c config
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	c.webAPIClientOptions = append(c.webAPIClientOptions, webapi.CacheUsers())
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
	if c.searchBotID {
		id := a.UserID(c.botName)
		if id == "" {
			return nil, fmt.Errorf("bot-name not found: %s", c.botName)
		}
		ret.Name = c.botName
		ret.ID = id
	}
	return &ret, nil
}

var (
	metaTag     = regexp.MustCompile(`<.*?>`)
	parentheses = strings.NewReplacer("&lt;", "<", "&gt;", ">")
)

// ReceiveMessage receives a message and passes it to a handler for processing.
func (c Client) ReceiveMessage(ctx context.Context, handler func(ctx context.Context, e *Event) error) error {
	return c.socketModeClient.ReceiveMessage(ctx, handler)
}

// PostMessage sends a message to the Slack channel.
func (c Client) PostMessage(ctx context.Context, channelID, msg string) error {
	_, err := c.webAPIClient.PostMessage(ctx, channelID, msg)
	return err
}

// RespondToCommand responds to the Slack command.
func (c Client) RespondToCommand(ctx context.Context, responseURL string, msg string, visible bool) error {
	return c.webAPIClient.RespondToCommand(ctx, responseURL, msg, visible)
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
		if v, ok := c.webAPIClient.User(id); ok {
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
func (c Client) UploadImage(ctx context.Context, channels []string, title, fileName, fileType, comment string, img io.Reader) error {
	return c.webAPIClient.UploadImage(ctx, channels, title, fileName, fileType, comment, img)
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	return c.socketModeClient.Close()
}

// UsersList lists all users in a Slack team.
// see. https://api.slack.com/methods/users.list
func (c Client) UsersList(ctx context.Context) ([]User, error) {
	return c.webAPIClient.UsersList(ctx)
}

// Users lists all users in a Slack team and returns it's userID map.
func (c Client) Users(ctx context.Context) (map[string]User, error) {
	return c.webAPIClient.Users(ctx)
}

// RefreshUsersCache updates the client's cached user map.
func (c *Client) RefreshUsersCache(ctx context.Context) error {
	return c.webAPIClient.RefreshUsersCache(ctx)
}

// User returns the user corresponding to user ID from the client's user cache.
func (c *Client) User(id string) (User, bool) {
	return c.webAPIClient.User(id)
}
