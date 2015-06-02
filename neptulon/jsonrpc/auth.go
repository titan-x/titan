package jsonrpc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nbusy/devastator/neptulon"
)

// todo: remove session.UserID and use session.data.UserID (but without locking as we know that UseID will be set only once)

func authMiddleware(conn *neptulon.Conn, msg *Message) {
	if conn.Session.Get("userid") != 0 {
		return
	}

	// client certificate authorization: certificate is verified by the TLS listener instance so we trust it
	peerCerts := conn.ConnectionState().PeerCertificates
	if len(peerCerts) > 0 {
		idstr := peerCerts[0].Subject.CommonName
		uid64, err := strconv.ParseUint(idstr, 10, 32)
		if err != nil {
			conn.Session.Set("error", fmt.Errorf("Cannot parse client message or method mismatched: %v", err))
			return
		}
		userID := uint32(uid64)
		log.Printf("Client connected with client certificate subject: %+v", peerCerts[0].Subject)
		conn.Session.Set("userid", userID)
	}
}
