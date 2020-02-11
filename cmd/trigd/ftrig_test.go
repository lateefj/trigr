package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/lateefj/trigr"
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

// TestDirectoryWatcher ... make sure directory watching is working
func TestDirectoryWatcher(t *testing.T) {
	testingPath := "/tmp/trigr_directory_watcher"
	os.MkdirAll(testingPath, 755)
	defer func() {
		os.RemoveAll(testingPath)
	}()

	triggers := make(chan *trigr.Trigger, 0)
	dw := NewDirectoryWatcher(testingPath, triggers, false)

	go func() {
		err := dw.Watch()
		if err != nil {
			log.Printf("Failed to watch directory %s with error: %s", testingPath, err)
			t.Fatalf("Error calling Watch %s", err)
		}
	}()

	f, err := os.Create(fmt.Sprintf("%s/test.dat", testingPath))
	defer f.Close()
	if err != nil {
		t.Fatalf("Failed to create test file %s", err)
	}
	f.Write([]byte("Testing one two three"))
	stopped := make(chan bool, 0)
	go func() {
		dw.Stop()
		stopped <- true
	}()

	select {
	case <-stopped:
		// All good here
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Failed to stop the directory watcher")
	}

}
