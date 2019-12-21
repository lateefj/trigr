package main

import (
	"testing"
	"time"
)

func TestDuplicateLimiter(t *testing.T) {
	dl := duplicateLimiter{changeDelay: 10 * time.Millisecond}
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

	time.Sleep(7 * time.Millisecond)
	if !dl.add(key) {
		t.Fatal("After sleep should be able to add but failed to add")
	}
}

/*
func TestDirectoryWatcher(t *testing.T) {
	testingPath := "/tmp/trigr_directory_watcher"
	os.MkdirAll(testingPath, 755)
	defer func() {
		os.RemoveAll(testingPath)
	}()

	triggers := make(chan *trigr.Trigger, 0)
	dw := NewDirectoryWatcher(testingPath, triggers, false)
	defer dw.Stop()
	go func() {
		err := dw.Watch()
		if err != nil {
			log.Printf("Failed to watch directory %s with error: %s", testingPath, err)
		}
	}()
}*/
