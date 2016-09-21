package inmem

import "github.com/titan-x/titan/data"

type middlewareChan struct {
	userID, connID string
}

type addReqChan struct {
	userID    string
	queuedReq queuedReq
}

func (q *Queue) worker() {
	for {
		select {
		case mid := <-q.middlewareChan:
			// start queue gorutine only once per connection
			if _, ok := q.conns[mid.userID]; !ok {
				q.conns[mid.userID] = mid.connID
				data.UserCount.Add(1)
				go q.processQueue(q.getQueueChan(mid.userID), mid.userID, mid.connID)
			}

		case userID := <-q.remUserChan:
			if _, ok := q.conns[userID]; ok {
				q.getQueueChan(userID).quit <- true
				delete(q.conns, userID)
				data.UserCount.Add(-1)
			}

		case req := <-q.addReqChan:
			data.QueueLength.Add(1)
			q.getQueueChan(req.userID).req <- req.queuedReq

		case userID := <-q.delQueueChan:
			delete(q.reqChans, userID)
		}
	}
}
