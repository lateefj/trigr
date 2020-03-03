package lsl

import (
	"bytes"
	"os"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaLog(t *testing.T) {
	buff := bytes.NewBufferString("")
	llog := &LuaLog{Output: buff}
	llog.Log("Test")
	if buff.String() != "Test\n" {
		t.Fatalf("Expected \n'Test\n'\n however got \n'%s'", buff.String())
	}
	buff.Reset()
}

func TestLuaLoader(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewLuaLoader(readBuff, writeBuff)
	loader.Log.Log("Test")
	if loader.envBuilt {
		t.Fatal("envBuilt should not have already happen")
	}
	loader.BuildEnv()
	if !loader.envBuilt {
		t.Fatal("Environment should be built")
	}
	loader.LoadAllStdLibs()

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

func TestLuaFunction(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewLuaLoader(readBuff, writeBuff)
	err := loader.LoadStdLibs([]string{"string"})
	if err != nil {
		t.Fatalf("Failed to load StdLib 'io' %s", err)
	}
	loader.BuildEnv()
	code := `
function first_char(st)
	log(string.sub(st, 0, 1))
end
	`
	fileBuf := bytes.NewBufferString(code)
	err = loader.Function(fileBuf, "first_char", lua.LString("test"))
	if err != nil {
		t.Fatalf("Could not call function 'first_char' with error: %s", err)
	}
	if strings.TrimSpace(writeBuff.String()) != "t" {
		t.Fatalf("Expected 't' but got '%s'", strings.TrimSpace(writeBuff.String()))
	}

}

func TestLuaLoaderTest(t *testing.T) {
	readBuff := bytes.NewBufferString("")
	writeBuff := bytes.NewBufferString("")
	loader := NewLuaLoader(readBuff, writeBuff)

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
	err = loader.TestFile(testFilePath)
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

	err = loader.TestFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed test did with error %s", err)
	}

	output = writeBuff.String()
	if !strings.Contains(output, "FAILED") {
		t.Errorf("Expected Fail test but output is \n%s", output)
	}
}
