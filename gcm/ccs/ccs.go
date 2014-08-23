// GCM Cloud Connection Server (XMPP) implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"errors"
	"fmt"
	json "github.com/bitly/go-simplejson"
	"time"
	"github.com/soygul/nbusy-server/xmpp"
)

const (
	GCM_HOST = "gcm.googleapis.com:5235"
	GCM_XML  = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`

	GCM_ACK = `{"to": "%v", "message_id": "%v", "message_type": "ack"}`

	RESET_CONN_ERR_NUM_TIME = 60
)

// var (
// 	ccsRateLimit int32 = 1000
// )

type Client struct {
	ID, Key     string
	Debug       bool
	xmppClient  *xmpp.Client
	isConnected bool
}

func New(id, key string, debug bool) (*Client, error) {
	id += "@gcm.googleapis.com"

	c := &Client{
		ID:    id,
		Key:   key,
		Debug: debug,
	}

	err := c.connect()
	return c, err

}

func (c *Client) connect() error {
	xmppClient, err := xmpp.NewClient(GCM_HOST, c.ID, c.Key, c.Debug)
	if err != nil {
		return err
	}

	c.xmppClient = xmppClient
	c.isConnected = true
	return nil

}

func (c *Client) recv(msgCh chan map[string]interface{}, errCh chan error) error {
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

func (c *Client) Recv(msgCh chan map[string]interface{}, errCh chan error) error {
	if !c.isConnected {
		return errors.New("no connection")
	}

	var errNum int

	go func(errNum int) {
		// reset errNum
		errNum = 0
		time.Sleep(RESET_CONN_ERR_NUM_TIME * time.Second)
	}(errNum)

Recv:
	err := c.recv(msgCh, errCh)
	for err != nil {
		if errNum > 3 {
			return err
		}

		// reconnect
		errNum++
		c.connect()
		goto Recv
	}
	panic("unreachable")
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
