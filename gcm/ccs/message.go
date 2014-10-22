package ccs

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

// Message is an XMPP <message> stanzas used in sending messages to the GCM CCS server.
// https://developer.android.com/google/gcm/ccs.html#format
type Message struct {
	To             string            `json:"to"`
	MessageID      string            `json:"message_id"`
	MessageType    string            `json:"message_type,omitempty"`
	Data           map[string]string `json:"data,omitempty"`
	CollapseKey    string            `json:"collapse_key,omitempty"`
	TimeToLive     int               `json:"time_to_live,omitempty"`
	DelayWhileIdle bool              `json:"delay_while_idle,omitempty"`
	ReturnReceipt  bool              `json:"delivery_receipt_requested,omitempty"`
}

// IncomingMessage is an XMPP <message> stanzas coming from the CCS server.
type IncomingMessage struct {
	From        string            `json:"from"`
	MessageID   string            `json:"message_id"`
	MessageType string            `json:"message_type"`
	Data        map[string]string `json:"data"`
	Error       string            `json:"error"`
}

// NewMessage creates a CCS message.
func NewMessage(id string) *Message {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Message{
		To:        id,
		MessageID: "m-" + strconv.Itoa(r.Intn(100000)),
		Data:      make(map[string]string),
	}
}

// SetData is a legacy function not to be used..
func (m *Message) SetData(key string, value string) {
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	m.Data[key] = value
}

func (m *Message) String() string {
	result, _ := json.Marshal(m)
	return string(result)
}
