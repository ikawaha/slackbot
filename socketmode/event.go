package socketmode

import (
	"encoding/json"
)

// EnvelopeType is the Slack event envelope type.
type EnvelopeType string

const (
	// EventsAPI indicates that the data is events API payload.
	EventsAPI EnvelopeType = "events_api"

	// Disconnect is the message expecting disconnects to the WebSocket connection.
	Disconnect EnvelopeType = "disconnect"

	// Interactive indicates that data is sent from a modal.
	Interactive EnvelopeType = "interactive"

	// SlashCommands indicates that data is sent from a slash command.
	SlashCommands EnvelopeType = "slash_commands"

	// Hello indicates that the client has successfully connected to the server.
	Hello EnvelopeType = "hello"
)

// Envelope represents the response of the Slack event API.
type Envelope struct {
	Type                   string          `json:"type"`
	EnvelopeID             string          `json:"envelope_id"`
	Payload                EventPayload    `json:"payload"`
	AcceptsResponsePayload bool            `json:"accepts_response_payload"`
	RetryAttempt           int             `json:"retry_attempt"`
	RetryReason            string          `json:"retry_reason"`
	Reason                 string          `json:"reason"`     // disconnect type
	DebugInfo              json.RawMessage `json:"debug_info"` // disconnect type
}

// EventPayload is a part of the Envelope.
type EventPayload struct {
	Type               string        `json:"type"`
	EventID            string        `json:"event_id"`
	Event              Event         `json:"event"`
	EventTime          int           `json:"event_time"`
	EventContext       string        `json:"event_context"`
	APIAppID           string        `json:"api_app_id"`
	Authorizations     []interface{} `json:"authorizations"`
	IsExtSharedChannel bool          `json:"is_ext_shared_channel"`
	TeamID             string        `json:"team_id"`
	Token              string        `json:"token"`

	// for slash_commands
	Command     string `json:"command"`
	UserID      string `json:"user_id"`
	Text        string `json:"text"`
	UserName    string `json:"user_name"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}

// Event represents the Slack event.
// see. https://api.slack.com/apis/connections/events-api#event_type_structure
type Event struct {
	Type        string `json:"type"`
	Channel     string `json:"channel"`
	UserID      string `json:"user"`
	ClientMsgID string `json:"client_msg_id"`
	AppID       string `json:"app_id"`
	BotID       string `json:"bot_id"`
	TeamID      string `json:"team"`
	Text        string `json:"text"`
	TS          string `json:"ts"`
	EventTS     string `json:"event_ts"`
	ChannelType string `json:"channel_type"`

	// extended for slash_command
	Command     string `json:"command"`
	UserName    string `json:"user_name"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}

// Acknowledge represents the payload type of the response back to Slack acknowledging.
// see. https://api.slack.com/apis/connections/socket-implement#acknowledge
type Acknowledge struct {
	EnvelopeID string `json:"envelope_id"`
}

// EventType is the Slack event type.
type EventType string

const (
	// AppMention is a Slack event type.
	// Subscribe to only the message events that mention your app or bot.
	AppMention EventType = "app_mention"

	// Message is a Slack event type.
	// A message was sent to a channel.
	Message EventType = "message"

	// SlashCommand is a slash command.
	SlashCommand = "slash_command"
)

// Is returns true, if the event type equals tne given event type.
func (e Event) Is(t EventType) bool {
	return EventType(e.Type) == t
}

// IsAppMention returns true, if the event type is "app_mention".
func (e Event) IsAppMention() bool {
	return e.Is(AppMention)
}

// IsMessage returns true, if the event type is "message".
func (e Event) IsMessage() bool {
	return e.Is(Message)
}

// IsSlashCommand returns true, if the event type is "slash_command".
func (e Event) IsSlashCommand() bool {
	return e.Is(SlashCommand)
}
