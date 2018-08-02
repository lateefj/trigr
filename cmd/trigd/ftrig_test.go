package main

import (
	"testing"
	"time"
)

func TestDuplicateLimiter(t *testing.T) {
	dl := newDuplicatLimiter()
	key := "foo"
	go func() {
		if !dl.add(key) {
			t.Fatal("Fist key add should succeed")
		}
	}()

	// Make sure the go routine starts
	time.Sleep(5 * time.Millisecond)
	if dl.add(key) {
		t.Fatal("Second key add should fail")
	}

	time.Sleep(2 * changeDelay * time.Millisecond)
	if !dl.add(key) {
		t.Fatal("After sleep should be able to add but failed to add")
	}
}
