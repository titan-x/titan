package devastator

// Queue is a message queue for connected devices.
// Messages are mapped from user ID to []interface{} array which may contain request, response, or notification messages.
type Queue map[uint32][]interface{}

// Send sends a message immediately to .
func (q *Queue) Send() {}

// Queue .
func (q *Queue) Queue() {}

// LocalQueue is an in-memory queue.
type LocalQueue struct{}
