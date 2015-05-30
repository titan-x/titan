// Package jsonrpc implements JSON-RPC 2.0 protocol for Neptulon framework.
package jsonrpc

import "encoding/json"

// JSON RPC 2.0 dialect, where version field is ommited for brevity.

// Request is a JSON RPC 2.0 request object.
type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Notification is a JSON RPC 2.0 notification object.
type Notification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Response is a JSON RPC 2.0 response object.
type Response struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result"`
}

// Error is a JSON RPC 2.0 error response object.
type Error struct {
	ID    string       `json:"id"`
	Error ErrorPayload `json:"error"`
}

// ErrorPayload is a JSON RPC 2.0 error response object's 'error' field.
type ErrorPayload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ------------------------- legacy ----------------------

// // ReqMsg is a JSON RPC 2.0 request/notification object. Version field is ommited for brevity.
// type ReqMsg struct {
// 	ID     string          `json:"id,omitempty"`
// 	Method string          `json:"method"`
// 	Params json.RawMessage `json:"params,omitempty"`
// }
//
// // ResMsg is a JSON RPC 2.0 response object. Version field is ommited for brevity.
// type ResMsg struct {
// 	ID     string      `json:"id"`
// 	Result interface{} `json:"result,omitempty"`
// 	Error  ResError    `json:"error,omitempty"`
// }
//
// // ResError is a JSON RPC 2.0 response error object.
// type ResError struct {
// 	Code    int         `json:"code"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// }
