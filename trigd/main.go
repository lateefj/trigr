package main

import (
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
	"net/http"
	"os"

	//"github.com/Shopify/go-lua"
	log "github.com/Sirupsen/logrus"
	"github.com/lateefj/trigr"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
	"golang.org/x/net/websocket"
)

type TriggerExec struct {
	trigr.Trigger
	Input  io.ReadCloser
	Output io.WriteCloser
}

var messageCount int

func PublishTrigger(ws *websocket.Conn) {
	defer ws.Close()
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
	_, out, err := os.Pipe()
	if err != nil {
		log.Errorf("Failed to get a pipe %s", err)
		return
	}
	luaPath := fmt.Sprintf("%s.lua", t.Type)
	if _, err := os.Stat(luaPath); err == nil {

		l := lua.NewState()
		defer l.Close()
		l.SetGlobal("trigr", luar.New(l, t))
		l.SetGlobal("out", luar.New(l, out))
		log.Debugf("luaPath is %s", luaPath)
		if err := l.DoFile(luaPath); err != nil {
			log.Errorf("Failed to run lua file %s", err)
			return
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
			messageCount++
			err = json.Unmarshal([]byte(reply), &tl)
			if err != nil {
				log.Errorf("Failed to unmarshal log %s", err)
				continue
			}

			l.SetGlobal("log", luar.New(l, tl))
			l.DoString("handle_log(log)")

			log.Debugf("Log text is is: %s", tl.Text)

		}
	}
}

func main() {
	messageCount = 0
	log.SetLevel(log.DebugLevel)
	log.Debug("Hmmmmm starting trigrd")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello all goood total messages %d", messageCount)))
	})
	http.Handle("/ws", websocket.Handler(PublishTrigger))
	//http.ListenAndServe(":7771", nil)
	http.ListenAndServe(":8080", nil)
}
