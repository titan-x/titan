// Package ccs provides GCM Cloud Connection Server (XMPP) client implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"errors"
	"fmt"
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
	xmppConn                   *xmpp.Client //xmpp.Connection
	isConnected                bool
}

// New function connects to GCM CCS server and returns the connection object and an error (only if there were any).
// Accepts a CCS endpoint URI (production or staging) along with relevant credentials.
// Optionally debug mode can be enabled to dump all CSS communications to stdout.
func New(endpoint, senderID, apiKey string, debug bool) (*Conn, error) {
	if !strings.Contains(senderID, gcmDomain) {
		senderID += "@" + gcmDomain
	}

	c := &Conn{
		Endpoint: endpoint,
		SenderID: senderID,
		APIKey:   apiKey,
		Debug:    debug,
	}

	err := c.open()
	if debug {
		if err == nil {
			fmt.Printf("New CCS connection established with parameters: %+v\n", c)
		} else {
			if err == nil {
				fmt.Printf("New CCS connection failed established with parameters: %+v\n", c)
			}
		}
	}

	return c, err
}

func (c *Conn) open() error {
	xc, err := xmpp.NewClient(c.Endpoint, c.SenderID, c.APIKey, c.Debug)
	if err != nil {
		return err
	}

	c.xmppConn = xc
	c.isConnected = true
	return nil
}

// Listen starts listening for incoming messages from the CCS connection.
func (c *Conn) Listen(msgChan chan map[string]interface{}, errChan chan error) error {
	if !c.isConnected {
		return errors.New("no ccs connection")
	}

	for {
		event, err := c.xmppConn.Recv() //xmppConn.Listen
		if err != nil {
			c.xmppConn.Close()
			c.isConnected = false
			return err
		}

		go func(event interface{}) {
			switch v := event.(type) {
			case xmpp.Chat: //xmpp.Message
				isGcmMessage, message, err := c.handleMessage(v.Other[0])
				if err != nil {
					errChan <- err
					return
				}
				if isGcmMessage {
					return
				}
				msgChan <- message
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
	fmt.Printf("Incoming message: %v", msg)
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
			errFormat := "From: %v, MessageID: %v, Error: %v"
			result := fmt.Sprintf(errFormat, from, messageID, err)
			return true, nil, errors.New(result)
		}
	} else {
		ack := fmt.Sprintf(gcmACK, from, messageID)
		c.xmppConn.SendOrg(fmt.Sprintf(gcmXML, ack)) //xmppConn.SendRaw -or- just .Send (and .SendMsg for the other)
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
	c.xmppConn.SendOrg(res) //xmppConn.SendRaw

	return nil
}
