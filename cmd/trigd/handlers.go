package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/lateefj/httphacks"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

var (
	ui    http.FileSystem
	bower http.FileSystem
)

func init() {
	ui = rice.MustFindBox("ui").HTTPBox()
	bower = rice.MustFindBox("ui/bower_components").HTTPBox()
}

func setupHandlers() {

	http.Handle("/", http.FileServer(ui))
	http.Handle("/bower_components", http.FileServer(bower))
	http.HandleFunc("/api/status", httphacks.TimeWrap(func(w http.ResponseWriter, r *http.Request) {
		ClientsConnected.Send("Status being checked")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"Happy happy joy joy!"}`))
	}))
	http.Handle("/ws/trigger", websocket.Handler(PublishTrigger))
	http.Handle("/ws", websocket.Handler(ReadMessages))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
}
