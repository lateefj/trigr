package lsl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

var (
	// Embed Lua source code
	luaBox packr.Box
	// Embed Lua test framework
	uTestBox   packr.Box
	lslFrame   string
	uTestFrame string

	// List of standard library functions
	stdLibs = map[string]lua.LGFunction{
		lua.LoadLibName:      lua.OpenPackage,
		lua.TabLibName:       lua.OpenTable,
		lua.IoLibName:        lua.OpenIo,
		lua.OsLibName:        lua.OpenOs,
		lua.StringLibName:    lua.OpenString,
		lua.MathLibName:      lua.OpenMath,
		lua.DebugLibName:     lua.OpenDebug,
		lua.ChannelLibName:   lua.OpenChannel,
		lua.CoroutineLibName: lua.OpenCoroutine,
	}
)

func init() {
	var err error
	luaBox = packr.NewBox("./lua")
	lslFrame, err = luaBox.MustString("env.lua")
	if err != nil {
		log.Fatalf("Failed to load lua env %s", err)
	}
	uTestBox = packr.NewBox("./lua/u-test")
	uTestFrame, err = uTestBox.MustString("u-test.lua")
	if err != nil {
		log.Fatalf("Failed to test env %s", err)
	}
}

// OutputWriter ... function for writing
type OutputWriter func(txt string)

// InputReader ... function for reading input
type InputReader func() string

// LuaLog ...  Log capture system
type LuaLog struct {
	Output io.Writer
}

// Log ... Just output some text
func (ll *LuaLog) Log(v interface{}) {
	msg := fmt.Sprintf("%v\n", v)
	log.Print(msg)
	ll.Output.Write([]byte(msg))
}

// Print ... Print out without log
func (ll *LuaLog) Print(v interface{}) {
	ll.Output.Write([]byte(fmt.Sprintf("%v\n", v)))
}

// LuaLoader ... Simple loader system
type LuaLoader struct {
	State    *lua.LState
	Input    io.Reader
	Output   io.Writer
	Log      *LuaLog
	EnvMap   map[string]lua.LValue
	env      *lua.LTable
	envBuilt bool
}

// NewLuaLoader ... creates a new LuaLoader
func NewLuaLoader(in io.Reader, out io.Writer) *LuaLoader {

	state := lua.NewState()
	ll := &LuaLoader{State: state, Input: in, Output: out, Log: &LuaLog{out}, EnvMap: make(map[string]lua.LValue), env: state.NewTable(), envBuilt: false}

	return ll
}

// SetGlobalVar ... Push into environment
func (ll *LuaLoader) SetGlobalVar(n string, v lua.LValue) {
	// If already exists then remove it first
	if _, exists := ll.EnvMap[n]; exists {
		delete(ll.EnvMap, n)
	}
	ll.env.RawSet(lua.LString(n), v)
	ll.State.SetFEnv(v, ll.env)
}

// BuildEnv ... Initializes the environment needed
func (ll *LuaLoader) BuildEnv() {
	// Don't build the more than once
	if ll.envBuilt {
		return
	}

	// Log wrappers
	ll.State.SetGlobal("log", luar.New(ll.State, ll.Log.Log))
	ll.State.SetGlobal("print", luar.New(ll.State, ll.Log.Print))

	// Setup the DSL
	for k, v := range ll.EnvMap {
		ll.env.RawSetString(k, v)
		ll.State.SetFEnv(v, ll.env)
	}

	ll.envBuilt = true
}

// LoadStdLibs ... Provide a list of libraries that are part of the standard library
func (ll *LuaLoader) LoadStdLibs(exp []string) error {
	for _, m := range exp {
		fnct, exists := stdLibs[m]
		if !exists {
			return errors.New(fmt.Sprintf("There is not standard library %s", m))
		}
		ll.State.Push(ll.State.NewFunction(fnct))
		ll.State.Push(lua.LString(m))
		ll.State.Call(1, 0)
	}
	return nil
}

// LoadAllStdLibs ... Wrapper around exposing all the OpenLibs()
func (ll *LuaLoader) LoadAllStdLibs() {
	ll.State.OpenLibs()
}

// DoFile ... Setups up DSL file
func (ll *LuaLoader) DoFile(path string) error {
	ll.BuildEnv()
	return ll.State.DoFile(path)
}

// Function ... Execute code path
func (ll *LuaLoader) Function(file io.Reader, name string, params ...lua.LValue) error {
	code, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	ff, err := ll.State.LoadString(string(code))
	if err != nil {
		return err
	}
	ll.BuildEnv()
	ll.State.Push(ff)
	ll.State.Push(lua.LString("llFunction"))
	ll.State.Call(1, 0)

	return ll.State.CallByParam(lua.P{
		Fn:      ll.State.GetGlobal(name),
		NRet:    1,
		Protect: true,
	}, params...)
}

// Code ... Execute code path
func (ll *LuaLoader) Code(code string) error {
	ff, err := ll.State.LoadString(string(code))
	if err != nil {
		return err
	}
	ll.BuildEnv()
	ll.State.Push(ff)
	ll.State.Push(lua.LString("llCode"))
	ll.State.Call(1, 0)
	return nil

}

// File ... Runs a DSL file
func (ll *LuaLoader) File(path string) error {
	ff, err := ll.State.LoadFile(path)
	if err != nil {
		return err
	}
	ll.BuildEnv()
	ll.State.Push(ff)
	ll.State.Push(lua.LString("llFile"))
	ll.State.Call(1, 0)
	return nil
}

// TestFile ... Just pass the path to the file to run the tests on
func (ll *LuaLoader) TestFile(path string) error {
	// Make sure the test code has everything available
	// require is a key function
	lua.OpenBase(ll.State)
	// Load all the standard library
	ll.LoadAllStdLibs()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	testPath := filepath.Dir(path)
	// To be able to import functions from the file it is testing need to modify the package path
	pc := fmt.Sprintf("package.path = package.path .. \";\" .. \"%s\" .. \"/?.lua\"\n", testPath)

	pcf, err := ll.State.LoadString(pc)
	if err != nil {
		return err
	}
	ll.State.Push(pcf)
	ll.State.Push(lua.LString("packageStuff"))
	ll.State.Call(1, 0)
	return ll.Test(f)
	return nil
}

// Test ... Execute test file in dsl mode
func (ll *LuaLoader) Test(file io.Reader) error {

	// Add u-test lua testing framework
	ut, err := ll.State.LoadString(uTestFrame)
	if err != nil {
		log.Fatalf("Failed to load u-test framework %s", err)
	}
	ll.State.SetGlobal("utest", ut)

	code, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	ff, err := ll.State.LoadString(string(code))
	if err != nil {
		return err
	}
	ll.BuildEnv()
	preCode := `
print(package.path)
test = utest()
	`
	utf, err := ll.State.LoadString(preCode)
	if err != nil {
		return err
	}
	ll.State.Push(utf)
	ll.State.Push(lua.LString("utestSetup"))
	ll.State.Call(1, 0)

	ll.State.Push(ff)
	ll.State.Push(lua.LString("llFile"))
	ll.State.Call(1, 0)
	postCode := `
  -- Call for the results of the test
  local tests, failed =  test.result()
  -- Print out the summary of results
  test.summary()
	`
	ll.State.LoadString(postCode)
	return nil
}

// Close ... End execution and exit
func (ll *LuaLoader) Close() {
	ll.State.Close()
}
