// Package ccs provides GCM Cloud Connection Server (XMPP) client implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

/*
Connection*, Channel

Connect*, Dial, Open

Accept, Listen*

Receive, Read*

Close*
*/

import (
	"errors"
	"fmt"
	"log"
	"strings"

	json "github.com/bitly/go-simplejson"
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
	MessageChan                chan map[string]interface{}
	ErrorChan                  chan error
	xmppConn                   *xmpp.Client
	isConnected                bool
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
		c.isConnected = true
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

// Listen starts listening for incoming messages from the CCS connection.
func (c *Conn) Listen() error {
	for {
		event, err := c.xmppConn.Recv()
		if err != nil {
			c.xmppConn.Close()
			c.isConnected = false
			return err
		}

		go func(event interface{}) {
			switch v := event.(type) {
			case xmpp.Chat:
				isGcmMessage, message, err := c.handleMessage(v.Other[0])
				if err != nil {
					c.ErrorChan <- err
					return
				}
				if isGcmMessage {
					return
				}
				c.MessageChan <- message
			}
		}(event)
	}
}

// Close a CSS connection.
func (c *Conn) Close() error {
	c.isConnected = false
	return c.xmppConn.Close()
}

func (c *Conn) handleMessage(msg string) (isGcmMessage bool, message map[string]interface{}, err error) {
	log.Printf("Incoming CCS message: %v\n", msg)
	jsonData, err := json.NewJson([]byte(msg))
	if err != nil {
		return false, nil, errors.New("unknow message")
	}

	from, _ := jsonData.Get("from").String()
	messageID, _ := jsonData.Get("message_id").String()
	if _, ok := jsonData.CheckGet("message_type"); ok {
		err, _ := jsonData.Get("error").String()

		switch v, _ := jsonData.Get("message_type").String(); v {
		case "ack":
			return true, nil, nil
		case "nack":
			errFormat := "From: %v, Message ID: %v, Error: %v"
			result := fmt.Sprintf(errFormat, from, messageID, err)
			return true, nil, errors.New(result)
		}
	} else {
		ack := fmt.Sprintf(gcmACK, from, messageID)
		c.xmppConn.SendOrg(fmt.Sprintf(gcmXML, ack))
	}

	if _, ok := jsonData.CheckGet("from"); ok {
		data, _ := jsonData.Get("data").Map()
		return false, data, nil
	}

	return false, nil, errors.New("unknow message")
}

// Send a GCM CCS message.
func (c *Conn) Send(message *Message) error {
	if !c.isConnected {
		return errors.New("no connection")
	}

	res := fmt.Sprintf(gcmXML, message)
	c.xmppConn.SendOrg(res)

	return nil
}
