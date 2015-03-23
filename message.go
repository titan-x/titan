package main

import "encoding/json"

// ReqMsg is a JSON RPC 2.0 request/notification object. Version field is ommited for brevity.
type ReqMsg struct {
	ID     string          `json:"id,omitempty"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// AuthGoogReqParams is the Google+ OAuth token wrapper.
type AuthGoogReqParams struct {
	token string
}

// ResMsg is a JSON RPC 2.0 response object. Version field is ommited for brevity.
type ResMsg struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  ResError    `json:"error,omitempty"`
}

// ResError is a JSON RPC 2.0 response error object.
type ResError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
