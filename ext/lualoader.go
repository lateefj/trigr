package ext

import (
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/lateefj/trigr"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

type LuaLoader struct {
	State  *lua.LState
	Input  io.Reader
	Output io.Writer
}

func NewLuaLoader(in io.Reader, out io.Writer) *LuaLoader {
	return &LuaLoader{State: lua.NewState(), Input: in, Output: out}
}

func (ll *LuaLoader) buildContext(trig *trigr.Trigger) {
	ll.State.SetGlobal("trig", luar.New(ll.State, trig))
	ll.State.SetGlobal("trig_out", luar.New(ll.State, ll.Output))
	ll.State.SetGlobal("trig_in", luar.New(ll.State, ll.Input))
}
func (ll *LuaLoader) Run(path string, trig *trigr.Trigger, out chan *trigr.Trigger) error {
	ll.buildContext(trig)
	if err := ll.State.DoFile(path); err != nil {
		return err
	}
	return nil
}

type LuaDslLoader struct {
	LuaLoader
	DslPath string
}

func NewLuaDslLoader(in io.Reader, out io.Writer, dslPath string) *LuaDslLoader {
	return &LuaDslLoader{LuaLoader: LuaLoader{State: lua.NewState(), Input: in, Output: out}, DslPath: dslPath}
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
