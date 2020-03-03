package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
)

func TestDslPublishTrigr(t *testing.T) {
	proj := NewProject("test")
	// Need to create a buffered channel to not block
	proj.Triggers = make(chan *trigr.Trigger, 1)
	env := make([]string, 0)
	tr := trigr.NewTrigger("TEST", make(map[string]interface{}))
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	l := ext.NewTrigSL(in, out)
	setupDsl(env, proj, tr, in, out, l)
	code := `
local tt = new_trigr("LUA-TEST", {foo = "bar"})
publish_trigr(tt)
	`
	proj.Triggers = make(chan *trigr.Trigger, 1)
	err := l.RunCode(code, tr, proj.Triggers)
	if err != nil {
		t.Fatalf("Error running code %s\n error: %s", code, err)
	}
	var tt *trigr.Trigger
	var more bool
	select {
	case tt, more = <-proj.Triggers: // Assigned the variable
	case <-time.After(10 * time.Millisecond): // Timeout
		t.Fatal("Timed out waiting for trigger to publish")
	}
	if !more {
		t.Fatal("Channel of project triggers is closed before finished")
	}
	if tt == nil {
		t.Fatal("Trigr variable 'tt' test is null but should be assigned")
	}
	if tt.Type != "LUA-TEST" {
		t.Fatalf("Expected 'LUA-TEST' as trigr Type but it is %s", tt.Type)
	}
}

func TestCallStreamFunc(t *testing.T) {
	cleanupFiles := make([]string, 0)
	defer func() { // Remove any created files
		for _, name := range cleanupFiles {
			err := os.Remove(name)
			if err != nil {
				fmt.Printf("Failed to remove file with error: %s\n", err)
			}
		}
	}()
	codeFilePath := "/tmp/call_stream_func_test.lua"

	codeFile, err := os.Create(codeFilePath)
	if err != nil {
		t.Fatalf("Failed to create file %s", codeFilePath)
	}
	cleanupFiles = append(cleanupFiles, codeFilePath)
	_, err = codeFile.WriteString(`
function handle_trigr(trig)
	publish_trigr(trig)
end`)
	if err != nil {
		t.Fatalf("Failed to write test code with error %s", err)
	}
	err = codeFile.Sync()
	if err != nil {
		t.Fatalf("Failed to sync writes to file error: %s", err)
	}
	proj := NewProject("test")
	// Need to create a buffered channel to not block
	proj.Triggers = make(chan *trigr.Trigger, 1)
	env := make([]string, 0)
	tr := trigr.NewTrigger("TEST", make(map[string]interface{}))
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	l := ext.NewTrigSL(in, out)
	setupDsl(env, proj, tr, in, out, l)
	err = callPipeline(codeFilePath, l, proj, tr)
	if err != nil {
		t.Fatalf("Error calling stream func: %s", err)
	}
	var tt *trigr.Trigger
	select {
	case tt = <-proj.Triggers: // Assigned the variable
	case <-time.After(10 * time.Millisecond): // Timeout
		t.Fatal("Timed out waiting for trigger to publish")
	}
	if tt == nil {
		t.Fatal("Trigr variable 'tt' test is null but should be assigned")
	}
	if tt.Type != "TEST" {
		t.Fatalf("Expected 'TEST' as trigr Type but it is %s", tt.Type)
	}

}
