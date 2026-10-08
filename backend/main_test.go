package main

import (
	"context"
	"sync"
	"testing"
)

// 1000 clicks at the same moment must give back 1..1000, each exactly once.
func TestIncrementConcurrent(t *testing.T) {
	const clicks = 1000
	m := &mutationResolver{&Resolver{}}

	results := make(chan int64, clicks)
	var wg sync.WaitGroup
	for range clicks {
		wg.Go(func() {
			v, _ := m.Increment(context.Background())
			results <- v
		})
	}
	wg.Wait()
	close(results)

	seen := make(map[int64]bool)
	for v := range results {
		if v < 1 || v > clicks || seen[v] {
			t.Fatalf("got %d: out of range or returned twice", v)
		}
		seen[v] = true
	}
	if len(seen) != clicks {
		t.Fatalf("got %d distinct values, want %d", len(seen), clicks)
	}
}
