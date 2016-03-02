package ext

import (
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/lateefj/trigr"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

type LuaLog struct {
	Output io.Writer
}

func (ll *LuaLog) Error(L *lua.LState) int {
	log.Debugf("Error is being called ...")
	s := L.ToString(1)
	log.Error(s)
	ll.Output.Write([]byte(fmt.Sprintf("Error: %s", s)))
	return 1
}

func (ll *LuaLog) Log(L *lua.LState) int {
	log.Debugf("Log is being called ...")
	s := L.ToString(1)
	log.Info(s)
	ll.Output.Write([]byte(s))
	return 1
}

type LuaLoader struct {
	State  *lua.LState
	Input  io.Reader
	Output io.Writer
	Log    *LuaLog
}

func NewLuaLoader(in io.Reader, out io.Writer) *LuaLoader {
	return &LuaLoader{State: lua.NewState(), Input: in, Output: out, Log: &LuaLog{out}}
}
func Double(L *lua.LState) int {
	lv := L.ToInt(1)            /* get argument */
	L.Push(lua.LNumber(lv * 2)) /* push result */
	return 1                    /* number of results */
}
func (ll *LuaLoader) buildContext(trig *trigr.Trigger) {
	ll.State.SetGlobal("trig", luar.New(ll.State, trig))
	// Log wrappers
	ll.State.SetGlobal("trig_error", ll.State.NewFunction(ll.Log.Error))
	ll.State.SetGlobal("trig_log", ll.State.NewFunction(ll.Log.Log))
	ll.State.SetGlobal("double", ll.State.NewFunction(Double))
}
func (ll *LuaLoader) Run(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	ll.buildContext(trig)
	if err := ll.State.DoFile(path); err != nil {
		return err
	}
	return nil
}

type LuaDslLoader struct {
	*LuaLoader
	DslPath string
}

func NewLuaDslLoader(in io.Reader, out io.Writer, dslPath string) *LuaDslLoader {
	return &LuaDslLoader{LuaLoader: NewLuaLoader(in, out), DslPath: dslPath}
}

func (ldl *LuaDslLoader) RunDsl(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	ldl.buildContext(trig)

	ldl.State.SetGlobal("trig_dsl_path", luar.New(ldl.State, path))
	if err := ldl.State.DoFile(fmt.Sprintf("%s/dslify.lua", ldl.DslPath)); err != nil {
		log.Errorf("Failed to do file %s error: %s", path, err)
		return err
	}
	return nil
}
