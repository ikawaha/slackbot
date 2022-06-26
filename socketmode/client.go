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
func (c Client) Close() error {
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

func (c *Client) reconnect() error {
	_ = c.Close()
	wss, err := connectionOpen(context.TODO(), c.token)
	if err != nil {
		return err
	}
	return c.dial(wss)
}

// ReceiveMessage receives a message and passes it to a handler for processing.
func (c Client) ReceiveMessage(ctx context.Context, handler func(context.Context, *Event) error) error {
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
		event, err := c.openEnvelope(msg)
		if err != nil {
			log.Println(err, "reconnect...")
			if err := c.reconnect(); err != nil {
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

func (c Client) openEnvelope(msg interface{}) (*Event, error) {
	switch t := msg.(type) {
	default:
		return nil, fmt.Errorf("unknown message type: %T, %+v", msg, msg)
	case error:
		return nil, t
	case *Envelope:
		return c.processEnvelope(t)
	}
	return nil, fmt.Errorf("unknown message type: %T, %+v", msg, msg)
}

func (c Client) processEnvelope(ev *Envelope) (*Event, error) {
	if c.debug {
		dump, _ := json.MarshalIndent(ev, "", "  ")
		log.Printf("envelope:%s", dump)
	}
	// ack
	if ev.EnvelopeID != "" {
		if err := websocket.JSON.Send(c.socket, Acknowledge{EnvelopeID: ev.EnvelopeID}); err != nil {
			return nil, fmt.Errorf("acknowledge error: %w", err)
		}
	}
	switch ev.Type {
	case "events_api":
		ret, err := extractEvent(ev)
		if err != nil {
			return nil, fmt.Errorf("dicpatch error: %w", err)
		}
		return ret, nil
	case "disconnect":
		log.Printf("refresh: event_type: %s, %q", ev.Type, ev.Payload)
		return nil, c.reconnect()
	default:
		log.Printf("skip: event_type: %s, payload: %q", ev.Type, ev.Payload)
		return nil, nil
	}
}

func extractEvent(e *Envelope) (*Event, error) {
	var p EventPayload
	if err := json.Unmarshal(e.Payload, &p); err != nil {
		return nil, err
	}
	return &p.Event, nil
}
