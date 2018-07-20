// Some basic contracts for the project

package trigr

import (
	"context"
	"encoding/json"
	"time"
)

type SourceType int

const (
	LogBufferSize = 10

	Directory SourceType = 1 << iota
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
	Timestamp int64                  `json:"timestamp"` // Timestamp
	Type      string                 `json:"type"`      // Type of trigger event
	Data      map[string]interface{} `json:"data"`      // Additional data associated
	Logs      chan *Log              `json:"-"`
	Context   context.Context        `json:"-"`
}

func NewTrigger(t string, data map[string]interface{}) *Trigger {
	// Create a new trigger with unbuffered log channel
	return &Trigger{Timestamp: milli(), Type: t, Data: data, Context: context.Background(), Logs: make(chan *Log, LogBufferSize)}
}

func (t *Trigger) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

type Source struct {
}

type LocalSource struct {
	Type SourceType `json:"type"`
	Path string     `json:"path"`
}
type RemoteSource struct {
	Type SourceType `json:"type"`
	Url  string     `json:"url"`
}

// Build
type Build struct {
	Id     string  `json:"id"`
	Commit *Commit `json:"commit"`
}

// Commit
type Commit struct {
	Hash   string  `json:"hash"`
	Branch *Branch `json:"branch"`
}

// Branch
type Branch struct {
	Name    string   `json:"name"`
	Project *Project `json:"branch"`
}

// Project
type Project struct {
	Id           string        `json:"id"`                      // Project name
	LocalSource  *LocalSource  `json:"local_source,omitempty"`  // Local source project configuration
	RemoteSource *RemoteSource `json:"remote_source,omitempty"` // Remote source project configuration
	Triggers     chan *Trigger `json:"-"`
	Logs         chan *Log     `json:"-"` // Stream of output
}

func NewProject(id string) *Project {
	return &Project{Id: id, Triggers: make(chan *Trigger), Logs: make(chan *Log)}
}

func (p *Project) AssignLocalSource(t SourceType, path string) {
	p.LocalSource = &LocalSource{Type: t, Path: path}
}
func (p *Project) AssignRemoteSource(t SourceType, url string) {
	p.RemoteSource = &RemoteSource{Type: t, Url: url}
}
