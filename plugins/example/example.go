package main

import (
	"fmt"

	"github.com/lateefj/trigr"
)

var V int

func F() { fmt.Printf("Hello, number %d\n", V) }

func Handle(event *trigr.Trigger) {
	fm.Printf("Trigger event is %v\n", event)
}
