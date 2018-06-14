package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	log "github.com/Sirupsen/logrus"
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
	TriggerChannel = make(chan *trigr.Trigger)
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
		log.Errorf("Error %s receiving reply", err)
		return
	}
	log.Debug("Received back from client trigr: " + reply)
	var t *trigr.Trigger
	err = json.Unmarshal([]byte(reply), t)
	if err != nil {
		log.Errorf("Failed to unmarshal %s", err)
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
			log.Errorf("Error %s receiving reply", err)
			break
		}
		log.Debug("Received back from client: " + reply)
		atomic.AddInt64(&messageCount, 1)
		err = json.Unmarshal([]byte(reply), &tl)
		if err != nil {
			log.Errorf("Failed to unmarshal log %s", err)
			continue
		}
	}
}

func handleTrigger(t *trigr.Trigger) {
	in, out, err := os.Pipe()
	if err != nil {
		log.Errorf("Failed to get a pipe %s", err)
		return
	}
	luaPath := fmt.Sprintf("%s.lua", t.Type)
	log.Debugf("Lua loading file %s", luaPath)
	if _, err := os.Stat(luaPath); err == nil {

		l := ext.NewTrigSL(in, out, "./lua")
		err = l.RunFile(luaPath, t, make(chan *trigr.Trigger))
		if err != nil {
			log.Errorf("Failed ot run dsl %s", err)
		}
	}
}

func main() {
	messageCount = 0
	log.SetLevel(log.DebugLevel)
	log.Debug("Hmmmmm starting trigrd")
	go func() {
		for t := range TriggerChannel {

			handleTrigger(t)
			b, err := json.Marshal(t)
			if err != nil {
				log.Errorf("Failed to marshal trigger %s", err)
			}
			ClientsConnected.Send(string(b))
		}
	}()

	dw := DirectoryWatcher{"./", TriggerChannel}
	go dw.Watch()
	setupHandlers()
}
