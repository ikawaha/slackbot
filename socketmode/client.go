package socketmode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	// DefaultTimeout represents the time to wait for a response from slack.
	DefaultTimeout = 10 * time.Second
)

const (
	appsConnectionsOpenEndpoint = `https://slack.com/api/apps.connections.open`
)

// Client represents a Slack client.
type Client struct {
	mux     sync.Mutex
	socket  *websocket.Conn
	token   string
	timeout time.Duration
	debug   bool
}

// New creates a slack bot with an app-level token.
func New(token string, opts ...Option) (*Client, error) {
	ret := Client{
		token:   token,
		timeout: DefaultTimeout,
	}
	wss, err := connectionOpen(context.TODO(), token)
	if err != nil {
		return nil, fmt.Errorf("api connection error, %w", err)
	}
	if err := ret.dial(wss); err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(&ret); err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// Close closes the client.
func (c *Client) Close() error {
	return c.socket.Close()
}

type socketOpenResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	URL   string `json:"url"`
}

func connectionOpen(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, appsConnectionsOpenEndpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("websocket requsest access failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("websocket request failed: %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("websocket response error: %w", err)
	}
	var r socketOpenResponse
	if err := json.Unmarshal(b, &r); err != nil {
		return "", fmt.Errorf("websocket response decode error: %w", err)
	}
	if !r.OK {
		return "", fmt.Errorf("websocket open error: %s", r.Error)
	}
	return r.URL, nil
}

func (c *Client) dial(url string) error {
	ws, err := websocket.Dial(url, "", "https://api.slack.com/")
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	defer c.mux.Unlock()
	c.mux.Lock()
	c.socket = ws
	return nil
}

func (c *Client) reconnect(ctx context.Context) error {
	_ = c.Close()
	wss, err := connectionOpen(ctx, c.token)
	if err != nil {
		return err
	}
	return c.dial(wss)
}

// ReceiveMessage receives a message and passes it to a handler for processing.
func (c *Client) ReceiveMessage(ctx context.Context, handler func(context.Context, *Event) error) error {
	ch := make(chan interface{}, 1)
	go func() {
		var e Envelope
		if err := websocket.JSON.Receive(c.socket, &e); err != nil {
			ch <- fmt.Errorf("receive error: %w", err)
		}
		ch <- &e
	}()
	select {
	case msg := <-ch:
		event, err := c.openEnvelope(ctx, msg)
		if err != nil {
			log.Println(err, ", reconnect...")
			if err := c.reconnect(ctx); err != nil {
				return err
			}
		}
		if event != nil {
			if err := handler(ctx, event); err != nil {
				return err
			}
		}
	case <-ctx.Done():
		return fmt.Errorf("context done")
	}
	return nil
}

func (c *Client) openEnvelope(ctx context.Context, msg interface{}) (*Event, error) {
	switch t := msg.(type) {
	case error:
		return nil, t
	case *Envelope:
		return c.processEnvelope(ctx, t)
	default:
		return nil, fmt.Errorf("unknown message type: %T, %+v", msg, msg)
	}
}

func (c *Client) processEnvelope(ctx context.Context, el *Envelope) (*Event, error) {
	if c.debug {
		dump, err := json.MarshalIndent(el, "", "  ")
		if err != nil {
			log.Printf("envelope marshal error: %v", err)
		}
		log.Printf("envelope:%s", dump)
	}
	// ack
	if el.EnvelopeID != "" {
		if err := websocket.JSON.Send(c.socket, Acknowledge{EnvelopeID: el.EnvelopeID}); err != nil {
			return nil, fmt.Errorf("acknowledge error: %w", err)
		}
	}
	switch EnvelopeType(el.Type) {
	case EventsAPI:
		return &el.Payload.Event, nil
	case SlashCommands:
		return newSlashCommandEvent(&el.Payload), nil
	case Disconnect:
		log.Printf("refresh: event_type: %s, %#+v", el.Type, el.Payload)
		return nil, c.reconnect(ctx)
	case Hello:
		log.Println("event_type: hello, client has successfully connected to the server")
	default:
		log.Printf("skip: event_type: %s, payload: %#+v", el.Type, el.Payload)
	}
	return nil, nil
}

func newSlashCommandEvent(p *EventPayload) *Event {
	return &Event{
		Type:        SlashCommand,
		UserID:      p.UserID,
		AppID:       p.APIAppID,
		TeamID:      p.TeamID,
		Text:        p.Text,
		Command:     p.Command,
		UserName:    p.UserName,
		ResponseURL: p.ResponseURL,
		TriggerID:   p.TriggerID,
	}
}
