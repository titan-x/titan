// GCM Cloud Connection Server (XMPP) client implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"github.com/soygul/nbusy-server/xmpp"
	json "github.com/bitly/go-simplejson"
	"errors"
	"strings"
	"fmt"
)

const (
	gcmXML    = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmACK    = `{"to": "%v", "message_id": "%v", "message_type": "ack"}`
	gcmDomain = "gcm.googleapis.com"
)

// Connects to GCM CCS server and returns the connection object and an error (only if there were any).
// Accepts a CCS endpoint URI (production or staging) along with relevant credentials.
// Optionally debug mode can be enabled to dump all CSS communications to stdout.
func New(endpoint, senderID, apiKey string, debug bool) (*Connection, error) {
	if (!strings.Contains(senderID, gcmDomain)) {
		senderID += "@"+gcmDomain
	}

	conn := &Connection{
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

//
type Connection struct {
	Endpoint, SenderID, APIKey string
	Debug                      bool
	xmppConn       *xmpp.Client //xmpp.Connection
	isConnected                bool
}

func (conn *Connection) connect() error {
	xmppConn, err := xmpp.NewClient(conn.Endpoint, conn.SenderID, conn.APIKey, conn.Debug)
	if err != nil {
		return err
	}

	conn.xmppConn = xmppConn
	conn.isConnected = true
	return nil
}

// Start listening for incoming messages from the CCS connection.
func (c *Connection) Listen(msgChan chan map[string]interface{}, errChan chan error) error {
	if !c.isConnected {
		return errors.New("no ccs connection")
	}

	for {
		event, err := c.xmppConn.Recv()//xmppConn.Listen
		if err != nil {
			c.xmppConn.Close()
			c.isConnected = false
			return err
		}

		go func(event interface{}) {
			switch v := event.(type) {
			case xmpp.Chat://xmpp.Message
				isGcmMessage, message, err := c.handleRecvMessage(v.Other[0])
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

func (c *Connection) handleRecvMessage(msg string) (isGcmMessage bool, message map[string]interface{}, err error) {
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
	} else {
		ack := fmt.Sprintf(gcmACK, from, messageId)
		c.xmppConn.SendOrg(fmt.Sprintf(gcmXML, ack))//xmppConn.SendRaw -or- just .Send (and .SendMsg for the other)
	}

	if _, ok := jsonData.CheckGet("from"); ok {
		data, _ := jsonData.Get("data").Map()
		return false, data, nil
	}

	return false, nil, errors.New("unknow message")
}

//func (c *connection) Send(message *Message) error {
//	if !c.isConnected {
//		return errors.New("no connection")
//	}
//
//	result := fmt.Sprintf(GCM_XML, message)
//	c.xmppConn.SendOrg(result) //xmppConn.SendRaw
//
//	return nil
//}
