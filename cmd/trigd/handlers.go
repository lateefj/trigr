package main

import (
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/lateefj/httphacks"
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/websocket"
)

var (
	ui http.FileSystem
	//bower http.FileSystem
	upgrader = websocket.Upgrader{} // use default options
)

func init() {
	ui = rice.MustFindBox("ui").HTTPBox()
	// Kill UI for now
	//bower = rice.MustFindBox("ui/bower_components").HTTPBox()
}

func ReadMessages(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error converting to websocket:", err)
		return
	}
	clientId, send := ClientsConnected.New()
	defer ClientsConnected.Remove(clientId)
	defer ws.Close()
	for m := range send {
		err = ws.WriteMessage(websocket.BinaryMessage, m)
		if err != nil {
			log.Println("Error WriteMessage :", err, m)
			return
		}
	}
}

/*func PublishTrigger(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error converting to websocket:", err)
		return
	}
	clientId, send := ClientsConnected.New()
	defer ClientsConnected.Remove(clientId)
	defer ws.Close()
	go func() {
		for m := range send {
			err = ws.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.Println("Error WriteMessage :", err, m)
			}
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
}*/

func setupHandlers() {

	http.Handle("/", http.FileServer(ui))
	//http.Handle("/bower_components", http.FileServer(bower))
	http.HandleFunc("/api/status", httphacks.TimeWrap(func(w http.ResponseWriter, r *http.Request) {
		ClientsConnected.Send([]byte("Status being checked"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"Happy happy joy joy!"}`))
	}))
	//http.HandleFunc("/ws/trigger", PublishTrigger)
	http.HandleFunc("/ws", ReadMessages)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
}
