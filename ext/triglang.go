package ext

import (
	"io"
	"time"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/lsl"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// TrigSL ... Simple language wrapper around lsl
type TrigSL struct {
	lsl.LuaLoader
}

// NewTrigSL ... creates a new TrigSL
func NewTrigSL(in io.Reader, out io.Writer) *TrigSL {
	ll := lsl.NewLuaLoader(in, out)
	return &TrigSL{*ll}
}

// Context builder
func (ll *TrigSL) buildContext(trig *trigr.Trigger, out chan *trigr.Trigger) {
	ll.BuildEnv()
	ll.SetGlobalVar("log", luar.New(ll.State, func(msg string) {
		trig.Logs <- &trigr.Log{Timestamp: time.Now().Unix(), Text: msg}
	}))
	ll.SetGlobalVar("print", luar.New(ll.State, func(msg string) {
		trig.Logs <- &trigr.Log{Timestamp: time.Now().Unix(), Text: msg}
	}))
	ll.SetGlobalVar("new_trigr", luar.New(ll.State, func(tType string, data map[string]interface{}) *trigr.Trigger {
		return trigr.NewTrigger(tType, data)
	}))

	ll.SetGlobalVar("publish_trigr", luar.New(ll.State, func(t *trigr.Trigger) {
		out <- t
	}))
	ll.SetGlobalVar("trig", luar.New(ll.State, trig))

}

// RunFunction ... Execute code path
func (ll TrigSL) RunFunction(file io.Reader, name string, trig *trigr.Trigger, out chan *trigr.Trigger, params ...lua.LValue) error {
	ll.buildContext(trig, out)
	return ll.Function(file, name, params...)
}

// RunCode ... Execute code path
func (ll TrigSL) RunCode(code string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	ll.buildContext(trig, out)
	return ll.Code(code)
}

// RunFile ... Runs a DSL file
func (ll *TrigSL) RunFile(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	ll.buildContext(trig, out)
	return ll.File(path)
}

// RunTest ... Execute test file in dsl mode
func (ll *TrigSL) RunTest(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	ll.buildContext(trig, out)
	return ll.TestFile(path)
}
