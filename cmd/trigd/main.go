package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
)

type TriggerExec struct {
	trigr.Trigger
	Input  io.ReadCloser
	Output io.WriteCloser
}

var (
	// Track total number of messages
	messageCount   int64 = 0
	TriggerChannel chan *trigr.Trigger
)

func init() {
	TriggerChannel = make(chan *trigr.Trigger, 0)
}

func handleTrigger(t *trigr.Trigger) {
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	// TODO: Should make this configurable
	luaPath := fmt.Sprintf("./.trigr/%s.lua", t.Type)
	//log.Printf("Lua loading file %s\n", luaPath)
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
	log.Printf("Starting trigd\n")
	go func() {
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

	// Default watch current directory
	dw := NewDirectoryWatcher("./", TriggerChannel, true)
	go dw.Watch()
	setupHandlers()
	log.Printf("Handled %d messages\n", messageCount)
}
