// Package ccs provides GCM Cloud Connection Server (XMPP) client implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/soygul/nbusy-server/xmpp"
)

const (
	gcmXML    = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmACK    = `{"to": "%v", "message_id": "%v", "message_type": "ack"}`
	gcmDomain = "gcm.googleapis.com"
)

// Conn is a GCM CCS connection.
type Conn struct {
	Endpoint, SenderID, APIKey string
	Debug                      bool
	xmppConn                   *xmpp.Client
}

// Connect connects to GCM CCS server denoted by endpoint URI (production or staging) along with relevant credentials.
// Debug mode dumps all CSS communications to stdout.
func Connect(endpoint, senderID, apiKey string, debug bool) (Conn, error) {
	if !strings.Contains(senderID, gcmDomain) {
		senderID += "@" + gcmDomain
	}

	c := Conn{
		Endpoint: endpoint,
		SenderID: senderID,
		APIKey:   apiKey,
		Debug:    debug,
	}

	xc, err := xmpp.NewClient(c.Endpoint, c.SenderID, c.APIKey, c.Debug)
	if err == nil {
		c.xmppConn = xc
	}

	if debug {
		if err == nil {
			log.Printf("New CCS connection established with parameters: %+v\n", c)
		} else {
			log.Printf("New CCS connection failed to establish with parameters: %+v\n", c)
		}
	}

	return c, err
}

// Read reads the next incoming messages from the CCS connection.
func (c *Conn) Read() (map[string]interface{}, error) {
	event, err := c.xmppConn.Recv()
	if err != nil {
		c.Close()
		return nil, err
	}

	switch v := event.(type) {
	case xmpp.Chat:
		isGcmMessage, message, err := c.handleMessage(v.Other[0])
		if err != nil {
			return nil, err
		}
		if isGcmMessage {
			return nil, nil
		}
		return message, nil
	}

	return nil, nil
}

func (c *Conn) handleMessage(msg string) (isGcmMessage bool, message map[string]interface{}, err error) {
	log.Printf("Incoming raw CCS message: %v\n", msg)
	var m IncomingMessage
	err = json.Unmarshal([]byte(msg), &m)
	if err != nil {
		return false, nil, errors.New("unknow message")
	}

	if m.MessageType != "" {
		switch m.MessageType {
		case "ack":
			return true, nil, nil
		case "nack":
			errFormat := "From: %v, Message ID: %v, Error: %v, Error Description: %v"
			result := fmt.Sprintf(errFormat, m.From, m.ID, m.Err, m.ErrDesc)
			return true, nil, errors.New(result)
		}
	} else {
		ack := fmt.Sprintf(gcmACK, m.From, m.ID)
		c.xmppConn.SendOrg(fmt.Sprintf(gcmXML, ack))
	}

	if m.From != "" {
		return false, m.Data, nil
	}

	return false, nil, errors.New("unknow message")
}

// Send a GCM CCS message.
func (c *Conn) Send(message *Message) error {
	res := fmt.Sprintf(gcmXML, message)
	c.xmppConn.SendOrg(res)
	return nil
}

// Close a CSS connection.
func (c *Conn) Close() error {
	return c.xmppConn.Close()
}
