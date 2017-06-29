package ext

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/lateefj/trigr"
)

func TestRunCode(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewTrigSL(readBuff, writeBuff, "./lua/env.lua")
	loader.Log.Log("Test")
	if writeBuff.String() != "Test\n" {
		t.Fatalf("Failed to do any logging expected 'Test' and got %s", writeBuff.String())
	}
	writeBuff.Reset()
	data := make(map[string]interface{})
	tig := trigr.NewTrigger("file", data)
	tig.Logs = make(chan *trigr.Log, 1)
	tstream := make(chan *trigr.Trigger, 0)
	err := loader.RunCode("log('code test')", tig, tstream)
	if err != nil {
		t.Fatalf("Failed to run code %s", err)
	}
	if writeBuff.String() != "code test\n" {
		t.Fatal("Failed to write code test to out")
	}
	fmt.Printf("Now going to waith on logs .....\n")
	l, more := <-tig.Logs
	if !more {
		t.Fatal("Channel closed before first log message is on it....")
	}
	if l.Text != "code test" {
		t.Fatalf("Failed to get log expected \n'code test\n'\n and got \n'%s'", l.Text)
	}
	if writeBuff.String() != "code test\n" {
		t.Fatalf("Failed to do any logging expected \n'code test\n'\n and got \n'%s'", writeBuff.String())
	}
	fmt.Printf("Done with log test ....\n")
}

func TestRunTest(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewTrigSL(readBuff, writeBuff, "./lua/env.lua")

	cleanupFiles := make([]string, 0)
	data := make(map[string]interface{})
	tig := trigr.NewTrigger("file", data)
	tig.Logs = make(chan *trigr.Log, 100)

	tstream := make(chan *trigr.Trigger, 0)
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
	local v = ext.add(1, 2)
	test.equal(v, 3)
end
`)
	testFile.Sync()
	err = loader.RunTest(testFilePath, tig, tstream)
	if err != nil {
		t.Fatalf("Failed to run test %s", err)
	}
	for {
		l, more := <-tig.Logs
		if !more {
			t.Fatal("Failed to find successful log")
		}
		fmt.Printf("log is %s\n", l.Text)
		if strings.Contains(l.Text, "PASSED") {
			break
		}
	}

	output := writeBuff.String()
	if !strings.Contains(output, "PASSED") {
		t.Errorf("Expected successful test but output is \n%s", output)
	}
	writeBuff.Reset()
	fmt.Printf("done with first passing test \n")
	// Write a failing test
	testFile.WriteAt([]byte(`local ext = require "examplet"
test.example_fail = function() 
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
	if !strings.Contains(output, "FAILED") {
		t.Errorf("Expected failed test but output is \n%s", output)
	}
	for {
		l, more := <-tig.Logs
		if !more {
			t.Fatal("Failed to find failed log")
		}
		if strings.Contains(l.Text, "FAILED") {
			break
		}
	}
}
