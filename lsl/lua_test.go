package lsl

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLuaLog(t *testing.T) {
	buff := bytes.NewBufferString("")
	llog := &LuaLog{Output: buff}
	llog.Error("Test")
	if buff.String() != "ERROR: Test\n" {
		t.Fatalf("Expected \n'ERROR: Test\n'\n however got \n'%s'", buff.String())
	}
	buff.Reset()
	llog.Info("Test")
	if buff.String() != "INFO: Test\n" {
		t.Fatalf("Expected 'INFO: Test' however got %s", buff.String())
	}

	buff.Reset()
	llog.Debug("Test")
	if buff.String() != "DEBUG: Test\n" {
		t.Fatalf("Expected 'DEBUG: Test\n however got '%s'", buff.String())
	}
	buff.Reset()
	llog.Log("Test")
	if buff.String() != "Test\n" {
		t.Fatalf("Expected 'Test\n' however got '%s'", buff.String())
	}
}

func TestLuaLoader(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewLuaLoader(readBuff, writeBuff, "./lua/env.lua")
	loader.Log.Log("Test")
	if loader.envBuilt {
		t.Fatal("envBuilt should not have already happen")
	}
	loader.BuildEnv()
	if !loader.envBuilt {
		t.Fatal("Environment should be built")
	}

	if writeBuff.String() != "Test\n" {
		t.Fatalf("Failed to do any logging expected 'Test' and got %s", writeBuff.String())
	}
	writeBuff.Reset()
	// Simple log test
	err := loader.Code("log('code test')")
	if err != nil {
		t.Fatalf("Failed to run code %s", err)
	}
	if writeBuff.String() != "code test\n" {
		t.Fatalf("Failed to do any logging expected 'code test' and got %s", writeBuff.String())
	}
}

func TestLuaLoaderTest(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewLuaLoader(readBuff, writeBuff, "./lua/env.lua")

	cleanupFiles := make([]string, 0)
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
	defer testFile.Close()

	cleanupFiles = append(cleanupFiles, testFilePath)
	testFile.WriteString(`local ext = require "examplet"
test.example_pass = function() 
	local v = ext.add(1, 2)
	test.equal(v, 3)
end
`)
	testFile.Sync()
	err = loader.Test(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	output := writeBuff.String()
	if !strings.Contains(output, "PASSED") {
		t.Errorf("Expected successful test but output is \n%s", output)
	}
	// Start writing at head of file
	testFile.Truncate(0)
	writeBuff.Reset()
	// Write a failing test
	testFile.WriteAt([]byte(`local ext = require "examplet"
test.example_fail = function() 
	local v = ext.add(1, 1)
	test.equal(v, 3)
end
`), 0)
	testFile.Sync()
	err = loader.Test(testFilePath)
	if err != nil {
		t.Fatalf("Failed test did with error %s", err)
	}

	output = writeBuff.String()
	if !strings.Contains(output, "FAILED") {
		t.Errorf("Expected Fail test but output is \n%s", output)
	}
}
