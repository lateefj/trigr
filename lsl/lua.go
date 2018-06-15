package lsl

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

var (
	luaBox     packr.Box
	uTestBox   packr.Box
	lslFrame   string
	uTestFrame string
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

// Error ... Log error
func (ll *LuaLog) Error(txt string) {
	msg := fmt.Sprintf("ERROR: %s", txt)
	ll.Log(msg)
}

// Info ... Info log wrapper
func (ll *LuaLog) Info(txt string) {
	msg := fmt.Sprintf("INFO: %s", txt)
	ll.Log(msg)
}

// Debug ... Very low level effort
func (ll *LuaLog) Debug(txt string) {
	msg := fmt.Sprintf("DEBUG: %s", txt)
	ll.Log(msg)
}

// Log ... Just output some text
func (ll *LuaLog) Log(v interface{}) {
	msg := fmt.Sprintf("%v\n", v)
	log.Print(msg)
	ll.Output.Write([]byte(msg))
}

// LuaLoader ... Simple loader system
type LuaLoader struct {
	State    *lua.LState
	Input    io.Reader
	Output   io.Writer
	Log      *LuaLog
	EnvMap   map[string]interface{}
	DslPath  string
	envBuilt bool
}

// NewLuaLoader ... creates a new LuaLoader
func NewLuaLoader(in io.Reader, out io.Writer, dslPath string) *LuaLoader {
	ll := &LuaLoader{State: lua.NewState(), Input: in, Output: out, Log: &LuaLog{out}, EnvMap: make(map[string]interface{}), DslPath: dslPath, envBuilt: false}
	return ll
}

// SetGlobalVar ... Push into enviroment
func (ll *LuaLoader) SetGlobalVar(n string, v interface{}) {
	// If already exists then remove it first
	if _, exists := ll.EnvMap[n]; exists {
		delete(ll.EnvMap, n)
	}
	ll.EnvMap[n] = v
}

func (ll *LuaLoader) BuildEnv() {
	// Don't build the more than once
	if ll.envBuilt {
		return
	}
	err := ll.State.DoString(lslFrame)
	if err != nil {
		log.Fatalf("Failed to load lsl framework %s", err)
	}
	// Log wrappers
	ll.SetGlobalVar("log_debug", luar.New(ll.State, ll.Log.Debug))
	ll.SetGlobalVar("log_info", luar.New(ll.State, ll.Log.Info))
	ll.SetGlobalVar("log_error", luar.New(ll.State, ll.Log.Error))
	ll.SetGlobalVar("log_output", luar.New(ll.State, ll.Log.Log))
	// EnvMap wrapper
	ll.State.SetGlobal("env_map", luar.New(ll.State, ll.EnvMap))
	ll.envBuilt = true
}

// Code ... Execute code path
func (ll *LuaLoader) Code(code string) error {
	ll.BuildEnv()
	return ll.State.CallByParam(lua.P{
		Fn:      ll.State.GetGlobal("run_code_with_env"),
		NRet:    1,
		Protect: true,
	}, lua.LString(code))
}

// File ... Runs a DSL file
func (ll *LuaLoader) File(path string) error {
	ll.BuildEnv()
	return ll.State.CallByParam(lua.P{
		Fn:      ll.State.GetGlobal("run_file_with_env"),
		NRet:    1,
		Protect: true,
	}, lua.LString(path))
}

// Test ... Execute test file in dsl mode
func (ll *LuaLoader) Test(path string) error {
	ll.State.SetGlobal("print", luar.New(ll.State, func(msg string) {
		ll.Output.Write([]byte(msg))
	}))

	ut, err := ll.State.LoadString(uTestFrame)
	if err != nil {
		log.Fatalf("Failed to load u-test framework %s", err)
	}
	ll.State.SetGlobal("utest", ut)
	ll.BuildEnv()
	testPath := filepath.Dir(path)
	return ll.State.CallByParam(lua.P{
		Fn:      ll.State.GetGlobal("run_test_with_env"),
		NRet:    1,
		Protect: true,
	}, lua.LString(path), lua.LString(testPath))
}

// Close ... End execution and exit
func (ll *LuaLoader) Close() {
	ll.State.Close()
}
