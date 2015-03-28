package main

// Token is an encrypted identifier for connecting devices.
type Token struct {
	ID uint32
	IV []byte
}
