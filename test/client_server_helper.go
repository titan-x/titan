package test

import "testing"

// ClientServerHelper wraps both client and server helper object for convenience.
type ClientServerHelper struct {
	Client *ClientHelper
	Server *ServerHelper
}

// NewClientServerHelper creates a new client-server wrapper.
func NewClientServerHelper(t *testing.T, useClientCert bool) *ClientServerHelper {
	s := NewServerHelper(t)
	c := NewClientHelper(t, useClientCert)

	return &ClientServerHelper{Client: c, Server: s}
}

// func (cs *ClientServerHelper) Close() {
//
// }

// Close closes the client connection and then stops the server instace.
func (cs *ClientServerHelper) Close() {
	cs.Client.Close()
	cs.Server.Stop()
}
