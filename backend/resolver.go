package main

import "sync"

//go:generate go tool gqlgen generate

// Resolver holds the shared counter. One Resolver is created at startup and
// shared by every request, and requests run at the same time, so mu guards count.
type Resolver struct {
	mu    sync.Mutex
	count int64
}
