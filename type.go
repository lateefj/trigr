// Some basic contracts for the project

package trigr

import (
	"encoding/json"
	"time"
)

const (
	Directory = 1 << iota
	Git
	Mercurial
)

func milli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type Log struct {
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

func NewLog(text string) *Log {
	return &Log{Timestamp: milli(), Text: text}
}

func (l *Log) Marshal() ([]byte, error) {
	return json.Marshal(l)
}

type Trigger struct {
	Timestamp int64                   `json:"timestamp"` // Timestamp
	Type      string                  `json:"type"`      // Type of trigger event
	Data      *map[string]interface{} `json:"data"`      // Additional data associated
	Logs      chan *Log               `json:"-"`
}

func NewTrigger(t string, data *map[string]interface{}) *Trigger {
	return &Trigger{Timestamp: milli(), Type: t, Data: data}
}

func (t *Trigger) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

type Source struct {
	Type int
	Path string
	Url  *string
}

func NewSource(t int, path string, url *string) *Source {
	return &Source{Type: t, Path: path, Url: url}
}

// Project
type Project struct {
	Name   string    // Project name
	Source *Source   // Handles dealing with source code ect
	Path   string    // Directory path
	Output chan *Log // Stream of output
}

// Get a paginated list of the history of triggers for the project
func History(offset, size int) []Trigger {
	return make([]Trigger, 0)
}
