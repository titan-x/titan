// GCM Cloud Connection Server (XMPP) client implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"github.com/soygul/nbusy-server/xmpp"
	"errors"
	"encoding/json"
	"strings"
	"fmt"
)

const (
	gcmXML    = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmACK    = `{"to": "%v", "message_id": "%v", "message_type": "ack"}`
	gcmDomain = "gcm.googleapis.com"
)

type connection struct {
	Endpoint, SenderID, APIKey string
	Debug                      bool
	xmppClient       *xmpp.Client
	isConnected                bool
}

// Connects to GCM CCS server and returns the connection object and an error, if there were any.
// Accepts a CCS endpoint URI (production or staging) along with relevant credentials.
// Optionally debug mode can be enabled to dump all CSS communications to stdout.
func New(endpoint, senderID, apiKey string, debug bool) (*connection, error) {
	if (!strings.Contains(senderID, gcmDomain)) {
		senderID += "@" + gcmDomain
	}

	conn := &connection{
		Endpoint: endpoint,
		SenderID: senderID,
		APIKey:   apiKey,
		Debug: debug,
	}

	if (debug) {
		fmt.Printf("Connection: %+v\n", conn)
	}

	err := conn.connect()
	return conn, err
}

func (conn *connection) connect() error {
	xmppClient, err := xmpp.NewClient(conn.Endpoint, conn.SenderID, conn.APIKey, conn.Debug)
	if err != nil {
		return err
	}

	conn.xmppClient = xmppClient
	conn.isConnected = true
	return nil
}

func (c *Client) Receive(msgCh chan map[string]interface{}, errCh chan error) error {
	if !c.isConnected {
		return errors.New("XMPP connection was closed. Cannot receive further from this channel.")
	}

	for {
		chat, err := c.xmppClient.Recv()
		if err != nil {
			c.xmppClient.Close()
			c.isConnected = false
			return err
		}

		go func(chat interface{}) {
			switch v := chat.(type) {
			case xmpp.Chat:
				isGcmMessage, message, err := c.handleRecvMessage(v.Other[0])
				if err != nil {
					errCh <- err
					return
				}
				if isGcmMessage {
					return
				}
			msgCh <- message
			}
		}(chat)
	}
}


func (c *Client) handleRecvMessage(msg string) (isGcmMessage bool, message map[string]interface{}, err error) {
	jsonData, err := json.NewJson([]byte(msg))
	if err != nil {
		return false, nil, errors.New("unknow message")
	}

	from, _ := jsonData.Get("from").String()
	messageId, _ := jsonData.Get("message_id").String()
	if _, ok := jsonData.CheckGet("message_type"); ok {
		err, _ := jsonData.Get("error").String()

		switch v, _ := jsonData.Get("message_type").String(); v {
		case "ack":
			return true, nil, nil
		case "nack":
			errFormat := "From: %v, MessageID: %v, Error: %v"
			result := fmt.Sprintf(errFormat, from, messageId, err)
			return true, nil, errors.New(result)
		}
		// ccsRateLimit = atomic.AddInt32(&ccsRateLimit, 1)
	} else {
		ack := fmt.Sprintf(GCM_ACK, from, messageId)
		c.xmppClient.SendOrg(fmt.Sprintf(GCM_XML, ack))
	}

	if _, ok := jsonData.CheckGet("from"); ok {
		data, _ := jsonData.Get("data").Map()
		return false, data, nil
	}

	return false, nil, errors.New("unknow message")

}

func (c *Client) Send(message *Message) error {
	if !c.isConnected {
		return errors.New("no connection")
	}
	// if ccsRateLimit <= 0 {
	// 	return errors.New("ccs rate limit")
	// }

	result := fmt.Sprintf(GCM_XML, message)
	c.xmppClient.SendOrg(result)

	// ccsRateLimit = atomic.AddInt32(&ccsRateLimit, -1)

	return nil
}
