package slackbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

var (
	reMsg = regexp.MustCompile(`(?:<@(.+)>)?(?::?\s+)?(.*)`)
)

// Bot represents a slack bot.
type Bot struct {
	ID       string
	Name     string
	Users    map[string]string
	Channels map[string]string
	Ims      map[string]string
	socket   *websocket.Conn
	counter  uint64
}

type connectResponse struct {
	OK       bool                        `json:"ok"`
	Error    string                      `json:"error"`
	URL      string                      `json:"url"`
	Self     struct{ ID, Name string }   `json:"self"`
	Users    []struct{ ID, Name string } `json:"users"`
	Channels []struct {
		ID, Name string
		IsMember bool `json:"is_member"`
	} `json:"channels"`
	Ims []struct {
		ID     string
		UserID string `json:"user"`
	} `json:"ims"`
}

// New creates a slack bot from API token.
// https://[YOURTEAM].slack.com/services/new/bot
func New(token string) (*Bot, error) {
	bot := Bot{
		Users:    map[string]string{},
		Channels: map[string]string{},
		Ims:      map[string]string{},
	}

	// access slack api
	resp, err := bot.rtmStart(token)
	if err != nil {
		return nil, fmt.Errorf("api connection error, %v", err)
	}
	if !resp.OK {
		return nil, fmt.Errorf("connection error, %v", resp.Error)
	}

	// get realtime connection
	if e := bot.dial(resp.URL); e != nil {
		return nil, e
	}

	// save properties
	bot.ID = resp.Self.ID
	bot.Name = resp.Self.Name
	for _, u := range resp.Users {
		bot.Users[u.ID] = u.Name
	}
	for _, c := range resp.Channels {
		if c.IsMember {
			bot.Channels[c.ID] = c.Name
		}
	}
	for _, im := range resp.Ims {
		bot.Ims[im.ID] = im.UserID
	}
	return &bot, nil
}

func (b Bot) rtmStart(token string) (*connectResponse, error) {
	q := url.Values{}
	q.Set("token", token)
	u := &url.URL{
		Scheme:   "https",
		Host:     "slack.com",
		Path:     "/api/rtm.start",
		RawQuery: q.Encode(),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with code %d", resp.StatusCode)
	}
	var body connectResponse
	dec := json.NewDecoder(resp.Body)
	if e := dec.Decode(&body); e != nil {
		return nil, fmt.Errorf("response decode error, %v", err)
	}
	return &body, nil
}

func (b *Bot) dial(url string) error {
	ws, err := websocket.Dial(url, "", "https://api.slack.com/")
	if err != nil {
		return fmt.Errorf("dial error, %v", err)
	}
	b.socket = ws
	return nil
}

// UserName returns a slack username from the user id.
func (b Bot) UserName(uid string) string {
	name, _ := b.Users[uid]
	return name
}

// GetMessage receives a message from the slack channel.
func (b Bot) GetMessage() (Message, error) {
	var msg Message
	err := websocket.JSON.Receive(b.socket, &msg)
	return msg, err
}

// PostMessage sends a message to the slack channel.
func (b *Bot) PostMessage(m Message) error {
	m.ID = atomic.AddUint64(&b.counter, 1)
	return websocket.JSON.Send(b.socket, m)
}

// Close implements the io.Closer interface.
func (b *Bot) Close() error {
	return b.socket.Close()
}

// Message represents a message.
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Channel string `json:"channel"`
	UserID  string `json:"user"`
	Text    string `json:"text"`
}

// TextBody returns the body of the message.
func (m Message) TextBody() string {
	matches := reMsg.FindStringSubmatch(m.Text)
	if len(matches) == 3 {
		return matches[2]
	}
	return ""
}

// MentionID returns the mention id of this message.
func (m Message) MentionID() string {
	matches := reMsg.FindStringSubmatch(m.Text)
	if len(matches) == 3 {
		return matches[1]
	}
	return ""
}
