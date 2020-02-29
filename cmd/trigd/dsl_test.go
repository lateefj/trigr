package main

import (
	"bytes"
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
	l := ext.NewTrigSL(in, out, "../../lsl/lua")
	setupDsl(env, proj, tr, in, out, l)
	code := `
local test = new_trigr("LUA-TEST", {foo = "bar"})
publish_trigr(test)
	`
	tChan := make(chan *trigr.Trigger, 1)
	err := l.RunCode(code, tr, tChan)
	if err != nil {
		t.Fatalf("Error running code %s\n error: %s", code, err)
	}
	var tt *trigr.Trigger
	select {
	case tt = <-proj.Triggers: // Assigned the variable
	case <-time.After(100 * time.Millisecond): // Timeout
		t.Fatal("Timed out waiting for trigger to publish")
	}
	if tt == nil {
		t.Fatal("Trigr variable 'tt' test is null ")
	}
	if tt.Type != "LUA-TEST" {
		t.Fatalf("Expected 'LUA-TEST' as trigr Type but it is %s", tt.Type)
	}
}
