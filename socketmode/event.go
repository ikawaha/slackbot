package socketmode

import (
	"encoding/json"
)

// Envelope represents the response of the Slack event API.
type Envelope struct {
	Type                   string          `json:"type"`
	EnvelopeID             string          `json:"envelope_id"`
	Payload                json.RawMessage `json:"payload"`
	AcceptsResponsePayload bool            `json:"accepts_response_payload"`
	RetryAttempt           int             `json:"retry_attempt"`
	RetryReason            string          `json:"retry_reason"`
	Reason                 string          `json:"reason"`     // disconnect type
	DebugInfo              json.RawMessage `json:"debug_info"` // disconnect type
}

// EventPayload is a part of the Envelope.
type EventPayload struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	Event        Event  `json:"event"`
	EventTime    int    `json:"event_time"`
	EventContext string `json:"event_context"`
	//	APIAppID           string `json:"api_app_id"`
	//	Authorizations     []any  `json:"authorizations"`
	//	IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
	//	TeamID             string `json:"team_id"`
	//	Token              string `json:"token"`
}

// Event represents the Slack event.
// see. https://api.slack.com/apis/connections/events-api#event_type_structure
type Event struct {
	Type        string `json:"type"`
	Channel     string `json:"channel"`
	UserID      string `json:"user"`
	ClientMsgID string `json:"client_msg_id"`
	BotID       string `json:"bot_id"`
	TeamID      string `json:"team"`
	Text        string `json:"text"`
	TS          string `json:"ts"`
	EventTS     string `json:"event_ts"`
	ChannelType string `json:"channel_type"`
}

// Acknowledge represents the payload type of the response back to Slack acknowledging.
// see. https://api.slack.com/apis/connections/socket-implement#acknowledge
type Acknowledge struct {
	EnvelopeID string `json:"envelope_id"`
}
