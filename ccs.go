// GCM Cloud Connection Server (XMPP) implementation.
// https://developer.android.com/google/gcm/ccs.html
package main

import "github.com/soygul/nbusy-server/xmpp"

const (
	gcmXML  = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmACK = `{"to": "%v", "message_id": "%v", "message_type": "ack"}`
	resetConn = 60
)

type Client struct {
	ID, Key     string
	Debug       bool
	xmppClient  *xmpp.Client
	isConnected bool
}
