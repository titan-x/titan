package neptulon_test

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/neptulon/neptulon"
)

const debug = false

// Example demonstrating the Neptulon server.
func Example() {
	type SampleMsg struct {
		Message string `json:"message"`
	}

	// start the server and echo incoming messages back to the sender
	s := neptulon.NewServer("127.0.0.1:3000")
	s.MiddlewareFunc(func(ctx *neptulon.ReqCtx) error {
		var msg SampleMsg
		if err := ctx.Params(&msg); err != nil {
			return err
		}
		ctx.Res = msg
		return ctx.Next()
	})
	go s.ListenAndServe()
	time.Sleep(time.Millisecond * 50) // let server goroutine to warm up
	defer s.Close()

	// connect to the server and send a message
	c, err := neptulon.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Connect("ws://127.0.0.1:3000"); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	_, err = c.SendRequest("echo", SampleMsg{Message: "Hello!"}, func(ctx *neptulon.ResCtx) error {
		wg.Done()
		var msg SampleMsg
		if err := ctx.Result(&msg); err != nil {
			return err
		}
		fmt.Println("Server says:", msg.Message)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	// Output: Server says: Hello!
}
