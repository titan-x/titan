package titan

import "github.com/neptulon/jsonrpc"

// Client is a Titan API client.
type Client struct {
	client *jsonrpc.Client
}
