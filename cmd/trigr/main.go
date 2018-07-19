package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lateefj/trigr"
)

var (
	addr       = flag.String("addr", "localhost:8080", "http service address")
	connection *websocket.Conn
	exit       = false
)

func consumeMessages(u url.URL) {

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Failed to connect with error:", err)
		// Breath before trying to reconnect
		time.Sleep(1 * time.Second)
		return
	}
	connection = c
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		//log.Printf("msg recieved: %s", message)
		// TODO: This is a crappy way to send multiple messages need a better message format
		if strings.Contains(string(message), "\"type\":") {
			var t trigr.Trigger
			err = json.Unmarshal(message, &t)
			if err != nil {
				log.Printf("Failed to unmarshal %s\n", err)
				continue
			}
			et := time.Unix(t.Timestamp, 0)
			ds, err := json.Marshal(t.Data)
			if err != nil {
				log.Printf("Failed to unmarshal data %s\n", err)
				continue
			}
			fmt.Printf("%s ➜ Trigger event type %s with data %s\n", et, t.Type, ds)
		} else {
			var l trigr.Log
			err = json.Unmarshal(message, &l)
			if err != nil {
				log.Printf("Failed to unmarshal %s\n", err)
				continue
			}
			et := time.Unix(l.Timestamp, 0)
			fmt.Printf("%s ➜ %s\n", et, l.Text)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	defer close(done)
	// Handle exit signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		for sig := range sigs {
			fmt.Printf("Exiting for sig %s\n", sig.String())
			connection.Close()
			exit = true
		}
	}()

	for {
		consumeMessages(u)
		if exit {
			os.Exit(0)
			return
		}
	}
}
