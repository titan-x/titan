package neptulon

import "github.com/nbusy/devastator/neptulon"

// Logger provides low level request logging, performance metrics, and other metrics data.
type Logger struct{}

func perfLoggerMiddleware(conn *neptulon.Conn, session *neptulon.Session, msg []byte) {
}
