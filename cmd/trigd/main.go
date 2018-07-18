package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
	"golang.org/x/net/websocket"
)

type TriggerExec struct {
	trigr.Trigger
	Input  io.ReadCloser
	Output io.WriteCloser
}

var (
	messageCount   int64
	TriggerChannel chan *trigr.Trigger
)

func init() {
	TriggerChannel = make(chan *trigr.Trigger, 0)
}

func ReadMessages(ws *websocket.Conn) {
	clientId, send := ClientsConnected.New()
	defer ClientsConnected.Remove(clientId)
	defer ws.Close()
	for m := range send {
		println("Sending message from socket", m)
		websocket.Message.Send(ws, m)
	}
}
func PublishTrigger(ws *websocket.Conn) {
	clientId, send := ClientsConnected.New()
	defer ClientsConnected.Remove(clientId)
	defer ws.Close()
	go func() {
		for m := range send {
			websocket.Message.Send(ws, m)
		}
	}()
	var reply string
	err := websocket.Message.Receive(ws, &reply)
	if err != nil {
		log.Printf("Error: %s receiving reply", err)
		return
	}
	log.Println("Received back from client trigr: " + reply)
	var t *trigr.Trigger
	err = json.Unmarshal([]byte(reply), t)
	if err != nil {
		log.Printf("Failed: to unmarshal %s", err)
		return
	}
	handleTrigger(t)
	var tl trigr.Log
	for {
		err = websocket.Message.Receive(ws, &reply)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error %s receiving reply\n", err)
			break
		}
		log.Println("Received back from client: " + reply)
		atomic.AddInt64(&messageCount, 1)
		err = json.Unmarshal([]byte(reply), &tl)
		if err != nil {
			log.Printf("Failed: to unmarshal log %s", err)
			continue
		}
	}
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
		err = l.RunFile(luaPath, t, make(chan *trigr.Trigger))
		if err != nil {
			log.Printf("Failed to run dsl %s\n", err)
		}
		//fmt.Printf(out.String())
	}
}

func main() {
	messageCount = 0
	log.Printf("Starting trigd\n")
	go func() {
		for t := range TriggerChannel {

			go func() {
				handleTrigger(t)
				b, err := json.Marshal(t)
				if err != nil {
					log.Printf("Failed to marshal trigger %s\n", err)
				}
				ClientsConnected.Send(string(b))
			}()
		}
	}()

	// Default watch current directory
	dw := NewDirectoryWatcher("./", TriggerChannel, true)
	go dw.Watch()
	setupHandlers()
}
