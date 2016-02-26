package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/lateefj/httphacks"
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
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

func statusHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Happy happy joy joy!"}`))
}

// Someday maybe move to a faster router
func handleHttp() {
	router := httprouter.New()
	router.ServeFiles("/", ui)
	router.ServeFiles("/bower", bower)
	router.GET("/api/status", statusHandler)
	router.GET("/ws", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		//fn := websocket.Handler(PublishTrigger)
		//fn(w, r)
	})
	http.ListenAndServe(":8080", router)
}

func setupHandlers() {

	http.Handle("/", http.FileServer(ui))
	http.Handle("/bower_components", http.FileServer(bower))
	http.HandleFunc("/api/status", httphacks.TimeWrap(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"Happy happy joy joy!"}`))
	}))
	http.Handle("/ws", websocket.Handler(PublishTrigger))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
	http.ListenAndServe(":8080", nil)
}
