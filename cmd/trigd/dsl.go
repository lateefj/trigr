package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// Configure the dsl
func setupDsl(env []string, p *Project, t *trigr.Trigger, input *bytes.Buffer, output *bytes.Buffer, l *ext.TrigSL) {
	l.SetGlobalVar("exec", func(cmd, directory string) string {
		split := strings.Split(cmd, " ")
		c := split[0]
		args := make([]string, 0)
		if len(split) > 1 {
			args = split[1:]
		}
		p := exec.Command(c, args...)
		p.Dir = directory
		p.Env = env
		t.Logs <- trigr.NewLog(fmt.Sprintf("running: %s ", cmd))
		output, err := p.CombinedOutput()
		if err != nil {
			t.Logs <- trigr.NewLog(err.Error())
		}
		return string(output)
	})

	l.SetGlobalVar("publish_trigr", func(t *trigr.Trigger) {
		p.Triggers <- t
	})
}

// handleTrigger mainly for local project configuration
func handleTrigger(env []string, p *Project, t *trigr.Trigger) {
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	if p.LocalSource != nil {
		path := p.LocalSource.Path
		// Handle streaming all trigr events to a single file
		streamPath := fmt.Sprintf("%s/.trigr/stream.lua", path)
		if _, err := os.Stat(streamPath); err == nil {
			println("Stream file exists sending all trigr event to it")
			l := ext.NewTrigSL(in, out, "./lsl/lua")
			setupDsl(env, p, t, in, out, l)
			/*err := l.DoFile(streamPath)
			if err != nil {
				msgErr := fmt.Sprintf("Failed to run dsl file %s with error %s\n", streamPath, err)
				t.Logs <- trigr.NewLog(msgErr)
				log.Print(msgErr)
			}*/
			if err := l.State.CallByParam(lua.P{
				Fn: l.State.GetGlobal("run_func_with_env"),
				//Fn:      l.State.GetGlobal("handle_trigr"),
				NRet:    1,
				Protect: true,
				//}, luar.New(l.State, t)); err != nil {
			}, lua.LString(streamPath), lua.LString("handle_trigr"), luar.New(l.State, t)); err != nil {
				msgErr := fmt.Sprintf("Failed to run dsl file %s with error %s\n", streamPath, err)
				t.Logs <- trigr.NewLog(msgErr)
				log.Print(msgErr)
			}
		}

		// Handle specific types of trigr events to a file
		luaPath := fmt.Sprintf("%s/.trigr/%s.lua", path, t.Type)
		if _, err := os.Stat(luaPath); err == nil {
			l := ext.NewTrigSL(in, out, "./lsl/lua")
			setupDsl(env, p, t, in, out, l)
			// Add the trig event to the context
			l.SetGlobalVar("trig", t)
			err = l.RunFile(luaPath, t, make(chan *trigr.Trigger))
			if err != nil {
				msgErr := fmt.Sprintf("Failed to run dsl file %s error %s\n", luaPath, err)
				t.Logs <- trigr.NewLog(msgErr)
				log.Print(msgErr)
			}
		}
	}
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func setAppend(items []string, e string) []string {
	if !contains(items, e) {
		items = append(items, e)
	}
	return items
}
