package ext

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lateefj/trigr"
)

func TestRunCode(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewTrigSL(readBuff, writeBuff)
	loader.Log.Log("Test")
	if writeBuff.String() != "Test\n" {
		t.Fatalf("Failed to do any logging expected 'Test' and got %s", writeBuff.String())
	}
	writeBuff.Reset()
	data := make(map[string]interface{})
	tig := trigr.NewTrigger("file", data)
	tig.Logs = make(chan *trigr.Log, 1)
	defer close(tig.Logs)
	tstream := make(chan *trigr.Trigger, 0)
	err := loader.RunCode("log('code test')", tig, tstream)
	if err != nil {
		t.Fatalf("Failed to run code %s", err)
	}
	var l *trigr.Log
	var more bool
	select {
	case l, more = <-tig.Logs:
	case <-time.After(1 * time.Millisecond):
		t.Fatal("Timeout waiting to get logs")
	}

	if !more {
		t.Fatal("Channel closed before first log message is on it....")
	}
	if l.Text != "code test" {
		t.Fatalf("Failed to get log expected \n'code test\n'\n and got \n'%s'", l.Text)
	}
}

func TestNewTrigr(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewTrigSL(readBuff, writeBuff)
	data := make(map[string]interface{})
	tig := trigr.NewTrigger("file", data)
	tig.Logs = make(chan *trigr.Log, 1)
	defer close(tig.Logs)
	tstream := make(chan *trigr.Trigger, 1)
	err := loader.RunCode(`
print('hello...')
local nt = new_trigr("LUA-TEST", {foo = "bar"})
publish_trigr(nt)
`, tig, tstream)
	if err != nil {
		t.Fatalf("Failed to run code %s", err)
	}
	var tr *trigr.Trigger
	var more bool
	select {
	case tr, more = <-tstream:
	case <-time.After(2 * time.Millisecond):
		t.Fatal("Timeout waiting trigr")
	}

	if !more {
		t.Fatal("Channel closed trigr is on it....")
	}
	if tr.Type != "LUA-TEST" {
		t.Fatalf("Failed to get log expected 'LUA-TEST' and got '%s'", tr.Type)
	}
}

func TestRunTest(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewTrigSL(readBuff, writeBuff)

	cleanupFiles := make([]string, 0)
	data := make(map[string]interface{})
	tig := trigr.NewTrigger("file", data)
	tig.Logs = make(chan *trigr.Log, 100)

	tstream := make(chan *trigr.Trigger, 1)
	defer func() { // Remove any created files
		for _, name := range cleanupFiles {
			os.Remove(name)
		}
	}()
	codeFilePath := "/tmp/examplet.lua"

	codeFile, err := os.Create(codeFilePath)
	if err != nil {
		t.Fatalf("Failed to create file %s", codeFilePath)
	}
	cleanupFiles = append(cleanupFiles, codeFilePath)
	codeFile.WriteString(`module("examplet", package.seeall)
function add(a, b)
	return a + b
end`)
	codeFile.Sync()
	testFilePath := "/tmp/examplet_test.lua"
	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("Failed to create file %s", testFilePath)
	}
	cleanupFiles = append(cleanupFiles, testFilePath)
	testFile.WriteString(`local ext = require "examplet"
test.example_pass = function() 
	log("test")
	local v = ext.add(1, 2)
	test.equal(v, 3)
end
`)
	testFile.Sync()
	err = loader.RunTest(testFilePath, tig, tstream)
	if err != nil {
		t.Fatalf("Failed to run test %s", err)
	}
	output := writeBuff.String()
	if !strings.Contains(output, "PASSED") {
		t.Errorf("Expected successful test but output is \n%s", output)
	}
	for {

		var l *trigr.Log
		var more bool
		select {
		case l, more = <-tig.Logs:
		case <-time.After(1 * time.Millisecond):
			t.Fatal("Timeout waiting to get logs")
		}
		if !more {
			t.Fatal("Failed to find successful log")
		}
		if strings.Contains(l.Text, "test") {
			break
		}
	}

	writeBuff.Reset()
	// Write a failing test
	testFile.WriteAt([]byte(`local ext = require "examplet"
test.example_fail = function() 
	log("test failure")
	local v = ext.add(1, 1)
	test.equal(v, 3)
end
`), 0)
	testFile.Sync()
	tig.Logs = make(chan *trigr.Log, 100)
	err = loader.RunTest(testFilePath, tig, tstream)
	if err != nil {
		t.Fatalf("Failed test did with error %s", err)
	}
	output = writeBuff.String()
	if !strings.Contains(output, "FAIL") {
		t.Errorf("Expected failed test but output is \n%s", output)
	}
	for {
		l, more := <-tig.Logs
		if !more {
			t.Fatal("Failed to find failed log")
		}
		if strings.Contains(l.Text, "test failure") {
			break
		}
	}
}
