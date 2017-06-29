package ext

import (
	"fmt"
	"io"
	"time"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/lsl"
	luar "layeh.com/gopher-luar"
)

type TrigSL struct {
	lsl.LuaLoader
}

// NewTrigSL ... creates a new TrigSL
func NewTrigSL(in io.Reader, out io.Writer, dslPath string) *TrigSL {
	ll := lsl.NewLuaLoader(in, out, dslPath)
	return &TrigSL{*ll}
}

func (ll *TrigSL) buildContext(trig *trigr.Trigger) {
	ll.BuildEnv()
	ll.SetGlobalVar("log", func(msg string) {
		fmt.Printf("Log function being called with msg %s\n", msg)
		ll.Log.Log(msg)
		trig.Logs <- &trigr.Log{Timestamp: time.Now().Unix(), Text: msg}
	})
	ll.SetGlobalVar("trig", luar.New(ll.State, trig))
}

// RunCode ... Execute code path
func (ll TrigSL) RunCode(code string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	defer close(trig.Logs)
	ll.buildContext(trig)
	return ll.Code(code)
}

// RunFile ... Runs a DSL file
func (ll *TrigSL) RunFile(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	defer close(trig.Logs)
	ll.buildContext(trig)
	return ll.File(path)
}

// RunTest ... Execute test file in dsl mode
func (ll *TrigSL) RunTest(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	defer close(trig.Logs)
	ll.buildContext(trig)
	return ll.Test(path)
}
