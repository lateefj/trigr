package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
)

var (
	// Track total number of messages
	messageCount int64 = 0
	confFile           = flag.String("conf", "~/.trigr/config.json", "Path to configuration file")
)

func handleTrigger(path string, t *trigr.Trigger) {
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	// TODO: Should make this configurable
	luaPath := fmt.Sprintf("%s/.trigr/%s.lua", path, t.Type)
	log.Printf("Lua loading file %s\n", luaPath)
	if _, err := os.Stat(luaPath); err == nil {
		// TODO: Lua dependent files should embedded into the binary
		l := ext.NewTrigSL(in, out, "./lsl/lua")
		// Add the trig event to the context
		l.SetGlobalVar("trig", t)
		//t.Logs <- trigr.NewLog(fmt.Sprintf("Running lua file %s", luaPath))
		err = l.RunFile(luaPath, t, make(chan *trigr.Trigger))
		if err != nil {
			log.Printf("Failed to run dsl %s\n", err)
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
	/*go func() {
		for t := range TriggerChannel {

			messageCount = messageCount + 1
			// First send the trigger out clients
			b, err := json.Marshal(t)
			if err != nil {
				log.Printf("Failed to marshal trigger %s\n", err)
				return
			}
			ClientsConnected.Send(b)
			// Next send logs to clients
			go func() {
				for l := range t.Logs {
					b, err := json.Marshal(l)
					if err != nil {
						log.Printf("Failed to trigger log %v error:  %s\n", l, err)
						continue
					}
					ClientsConnected.Send(b)
				}
			}()
			handleTrigger(t)
		}
	}()
	*/
	// Default watch current directory
	//dw := NewDirectoryWatcher("./", TriggerChannel, true)
	//go dw.Watch()
	setupHandlers()
	log.Printf("Handled %d messages\n", messageCount)
}
