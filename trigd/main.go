package main

import (
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
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

var messageCount int64

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
	var t trigr.Trigger
	err = json.Unmarshal([]byte(reply), &t)
	if err != nil {
		log.Errorf("Failed to unmarshal %s", err)
		return
	}
	in, out, err := os.Pipe()
	if err != nil {
		log.Errorf("Failed to get a pipe %s", err)
		return
	}
	luaPath := fmt.Sprintf("%s.lua", t.Type)
	if _, err := os.Stat(luaPath); err == nil {

		l := ext.NewLuaDslLoader(in, out, "./lua")
		err = l.RunDsl(luaPath, &t, make(chan *trigr.Trigger))
		if err != nil {
			log.Errorf("Failed ot run dsl %s", err)
		}

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

			/*l.SetGlobal("log", luar.New(l, tl))
			l.DoString("handle_log(log)")
			log.Debugf("Log text is is: %s", tl.Text)*/
		}
	}
}

func main() {
	messageCount = 0
	log.SetLevel(log.DebugLevel)
	log.Debug("Hmmmmm starting trigrd")
	setupHandlers()
	//handleHttp()
}
