package main

import (
	"fmt"

	"github.com/lateefj/trigr"
)

func Handle(event *trigr.Trigger) {
	fmt.Printf("Trigger event is %v\n", event)
}
