package slackbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

const (
	//DefaultTimeout represents the time to wait for a response from slack.
	DefaultTimeout = time.Minute

	eventTypePing = "ping"
)

// Client represents a slack client.
type Client struct {
	ID       string
	Name     string
	Users    map[string]string
	Channels map[string]string
	Ims      map[string]string
	socket   *websocket.Conn
	counter  uint64
	token    string
	timeout time.Duration
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
func New(token string) (*Client, error) {
	bot := Client{
		Users:    map[string]string{},
		Channels: map[string]string{},
		Ims:      map[string]string{},
		token:    token,
		timeout:  DefaultTimeout,
	}

	// access slack api
	resp, err := bot.rtmStart(token)
	if err != nil {
		return nil, fmt.Errorf("api connection error, %v", err)
	}
	if !resp.OK {
		return nil, fmt.Errorf("connection error, %v", resp.Error)
	}

	// get real time connection
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

// SetTimeout sets client timeout.
// If you try to set a timeout less than DefaultTimeout, DefaultTimeout is set.
func (c *Client)SetTimeout(timeout time.Duration) {
	if timeout < time.Minute {
		c.timeout = timeout
	}
}

// Timeout returns the client timeout setting.
func (c Client) Timeout() time.Duration {
	return c.timeout
}

func (c Client) rtmStart(token string) (*connectResponse, error) {
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

func (c *Client) dial(url string) error {
	ws, err := websocket.Dial(url, "", "https://api.slack.com/")
	if err != nil {
		return fmt.Errorf("dial error, %v", err)
	}
	c.socket = ws
	return nil
}

// UserName returns a slack username from the user id.
func (c Client) UserName(uid string) string {
	name, _ := c.Users[uid]
	return name
}

// GetMessage receives a message from the slack channel.
func (c Client) GetMessage(ctx context.Context) (Message, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	ch := make(chan error, 1)

	go func(ctx context.Context, waiting time.Duration) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		select {
		case <-ctx.Done():
		case <-time.After(waiting):
			if err := websocket.JSON.Send(c.socket, &Message{Type: eventTypePing, Time: time.Now().Unix()}); err != nil {
				log.Printf("ping error, %v", err)
			}
		}
	}(ctx, c.timeout-time.Second)

	var msg Message
	go func() {
		ch <- websocket.JSON.Receive(c.socket, &msg)
	}()

	select {
	case err := <-ch:
		return msg, err
	case <-ctx.Done():
		return msg, fmt.Errorf("connection lost timeout")
	}
	return msg, nil
}

var (
	metaTag     = regexp.MustCompile(`<.*?>`)
	parentheses = strings.NewReplacer("&lt;","<","&gt;",">")
)

// PlainMessageText resolves meta tags of the message text and return it.
func (c Client) PlainMessageText(msg string) string {
	txt := metaTag.ReplaceAllStringFunc(msg, func(s string) string {
		var id string
		for i :=0; i < len(s)-2; i++ {
			if s[i] == '@' {
				id = s[i+1:len(s)-1]
				break
			}
		}
		if v, ok := c.Users[id]; ok {
			return "@"+v
		}
		if id !=""{
			return "@"+id
		}
		return s
	})
	return parentheses.Replace(txt)
}

// PostMessage sends a message to the slack channel.
func (c Client) PostMessage(m Message) error {
	m.ID = atomic.AddUint64(&c.counter, 1)
	return websocket.JSON.Send(c.socket, m)
}

// UploadImage uploads a image by files.upload API.
func (c Client) UploadImage(channels []string, title, fileName, fileType, comment string, img io.Reader) error {
	if c.token == "" {
		return fmt.Errorf("slack token is empty")

	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, err := mw.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("multipart create from file error, %v, %v", title, err)
	}
	if _, err := io.Copy(part, img); err != nil {
		return fmt.Errorf("file copy error, %v, %v", title, err)
	}
	// for slack settings
	settings := map[string]string{
		"token":           c.token,
		"channels":        strings.Join(channels, ","),
		"filetype":        fileType,
		"title":           title,
		"initial_comment": comment,
	}
	for k, v := range settings {
		if err := mw.WriteField(k, v); err != nil {
			return fmt.Errorf("write field error, %v:%v, %v", k, v, err)
		}
	}
	if err := mw.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", &buf)
	if err != nil {
		return fmt.Errorf("slack files.uplad new request error, %v", err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	cl := &http.Client{Timeout: 10 * time.Second}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("slack files.upload error, %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response error, %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack files.upload, %v, %v", resp.Status, body)
	}
	return nil
}

// Close implements the io.Closer interface.
func (c *Client) Close() error {
	return c.socket.Close()
}
