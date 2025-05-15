package entity

import "time"

// Message represents a message in 
type Message struct {
	ID      string                `json:"ID"`
	From    Address        `json:"From"`
	To      []Address      `json:"To"`
	Content MessageContent `json:"Content"`
	Created time.Time             `json:"Created"`
	MIME    MIME           `json:"MIME"`
	Raw     Raw            `json:"Raw"`
}

// Address represents an email address
type Address struct {
	Relays  []string `json:"Relays"`
	Mailbox string   `json:"Mailbox"`
	Domain  string   `json:"Domain"`
	Params  string   `json:"Params"`
}

// MessageContent represents the content of a message
type MessageContent struct {
	Headers map[string][]string `json:"Headers"`
	Body    string              `json:"Body"`
	Size    int                 `json:"Size"`
}

// MIME represents MIME information
type MIME struct {
	Parts []any `json:"Parts"`
}

// Raw represents raw message data
type Raw struct {
	From string   `json:"From"`
	To   []string `json:"To"`
	Data string   `json:"Data"`
	Helo string   `json:"Helo"`
}

type ErrResp map[string]any