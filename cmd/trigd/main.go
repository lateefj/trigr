package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
)

var (
	// Track total number of messages
	messageCount int64 = 0
	confFile           = flag.String("conf", "~/.trigr/config.json", "Path to configuration file")
)

func setOutput() string {

	cmd := exec.Command("bash", "-c", "set")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("FAILED WITH ERROR: %s\n", err)
	}
	return string(output)
}

func handleTrigger(path string, t *trigr.Trigger) {
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	// TODO: Should make this configurable
	luaPath := fmt.Sprintf("%s/.trigr/%s.lua", path, t.Type)
	//log.Printf("Lua loading file %s\n", luaPath)
	if _, err := os.Stat(luaPath); err == nil {
		// TODO: Lua dependent files should embedded into the binary
		l := ext.NewTrigSL(in, out, "./lsl/lua")
		l.SetGlobalVar("exec", func(cmd string) string {
			t.Logs <- trigr.NewLog(setOutput())
			p := exec.Command(cmd)
			p.Dir = path
			p.Env = os.Environ()
			t.Logs <- trigr.NewLog(fmt.Sprintf("Env is \n%s\n", p.Env))
			t.Logs <- trigr.NewLog(fmt.Sprintf("running: %s ", cmd))
			output, err := p.CombinedOutput()
			if err != nil {
				t.Logs <- trigr.NewLog(err.Error())
			}
			return string(output)
		})
		// Add the trig event to the context
		l.SetGlobalVar("trig", t)
		//t.Logs <- trigr.NewLog(fmt.Sprintf("Running lua file %s", luaPath))
		err = l.RunFile(luaPath, t, make(chan *trigr.Trigger))
		if err != nil {
			msgErr := fmt.Sprintf("Failed to run dsl %s\n", err)
			t.Logs <- trigr.NewLog(msgErr)
			log.Print(msgErr)
		}
	}
}

func main() {
	flag.Parse()
	log.Printf("Starting trigd\n")
	if _, err := os.Stat(*confFile); err == nil {
		log.Printf("Config file path %s\n", *confFile)
		bits, err := ioutil.ReadFile(*confFile)
		if err != nil {
			log.Printf("ERROR: Reading file %s error: %s\n", *confFile, err)
		}
		pm, err := LoadProjectManager(bits)
		if err != nil {
			log.Printf("ERROR: Loading file %s error: %s\n", *confFile, err)
		}
		Manager = pm
	} else {
		log.Printf("ERROR: Configuration file does not exist: %s\n", *confFile)
	}
	// Default watch current directory
	//dw := NewDirectoryWatcher("./", TriggerChannel, true)
	//go dw.Watch()
	setupHandlers()
	log.Printf("Handled %d messages\n", messageCount)
}
