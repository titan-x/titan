package devastator

type Conns struct {
	userToConn map[string]string // user ID -> conn ID
	connToUser map[string]string // conn ID -> user ID

}
