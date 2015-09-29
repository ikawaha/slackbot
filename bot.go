package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// Bot represents slack bot.
type Bot struct {
	ID      string
	Users   map[string]string
	socket  *websocket.Conn
	counter uint64
}

type connectResponse struct {
	OK    bool                        `json:"ok"`
	Error string                      `json:"error"`
	Url   string                      `json:"url"`
	Self  struct{ ID string }         `json:"self"`
	Users []struct{ ID, Name string } `json:"users"`
}

// New creates a slack bot from API token.
// https://[YOURTEAM].slack.com/services/new/bot
func New(token string) (*Bot, error) {
	bot := Bot{Users: map[string]string{}}

	// access slack api
	resp, err := bot.connect(token)
	if err != nil {
		return nil, fmt.Errorf("api connection error, %v", err)
	}
	if !resp.OK {
		return nil, fmt.Errorf("connection error, %v", resp.Error)
	}

	// get realtime connection
	if e := bot.dial(resp.Url); e != nil {
		return nil, e
	}

	// save properties
	bot.ID = resp.Self.ID
	for _, u := range resp.Users {
		bot.Users[u.ID] = u.Name
	}

	return &bot, nil
}

func (b Bot) connect(token string) (*connectResponse, error) {
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

func (b *Bot) dial(url_ string) error {
	ws, err := websocket.Dial(url_, "", "https://api.slack.com/")
	if err != nil {
		return fmt.Errorf("dial error, %v", err)
	}
	b.socket = ws
	return nil
}

// Message represents a message.
type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// GetMessage receives a message from the slack channel.
func (b Bot) GetMessage() (Message, error) {
	var msg Message
	err := websocket.JSON.Receive(b.socket, &msg)
	return msg, err
}

// PostMessage sends a message to the slack channel.
func (b *Bot) PostMessage(m Message) error {
	m.Id = atomic.AddUint64(&b.counter, 1)
	return websocket.JSON.Send(b.socket, m)
}

// Close implements the io.Closer interface.
func (b *Bot) Close() error {
	return b.socket.Close()
}
