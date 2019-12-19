package main

import (
	"fmt"
	"time"

	"log"

	"github.com/lateefj/trigr"
	"golang.org/x/net/websocket"
)

func main() {
	//url := "ws://localhost:7771/ws"
	url := "ws://localhost:8080/ws/trigger"
	for {
		conn, err := websocket.Dial(url, "", "http://localhost")
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
		t := trigr.NewTrigger("test", nil)
		b, err := t.Marshal()
		if err != nil {
			log.Fatalf("Failed to marshal log %s", err)

		}
		_, err = fmt.Fprint(conn, string(b))
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < 1; i++ {
			l := trigr.NewLog(fmt.Sprintf("Hello %d\n", i))
			b, err := l.Marshal()
			if err != nil {
				log.Printf("Failed to marshal log %s", err)
				continue
			}
			_, err = fmt.Fprint(conn, string(b))
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			time.Sleep(1 * time.Second)
		}

	}
}
