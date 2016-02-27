// Some basic contracts for the project

package trigr

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	Timestamp int64                  `json:"timestamp"` // Timestamp
	Type      string                 `json:"type"`      // Type of trigger event
	Data      map[string]interface{} `json:"data"`      // Additional data associated
	Logs      chan *Log              `json:"-"`
}

func NewTrigger(t string, data map[string]interface{}) Trigger {
	return Trigger{Timestamp: milli(), Type: t, Data: data}
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

// Process represents a step in the entire pipeline
type Process struct {
	Name string
}

// Pipeline represents the entire development process
type Pipeline struct {
	Configuration map[string]Process
	Prepare       map[string]Process
	Build         map[string]Process
	Package       map[string]Process
	Deploy        map[string]Process
	Running       []Process
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Configuration: make(map[string]Process),
		Prepare:       make(map[string]Process),
		Build:         make(map[string]Process),
		Package:       make(map[string]Process),
		Deploy:        make(map[string]Process),
		Running:       make([]Process, 0),
	}
}

// Project
type Project struct {
	Name     string    `json:"name"`   // Project name
	Source   *Source   `json:"source"` // Handles dealing with source code ect
	Path     string    `json:"path"`   // Directory path
	Output   chan *Log // Stream of output
	Pipeline Pipeline  // Pipeline models the entire process for building the project
}

func LoadProject(path string) (Project, error) {
	var p Project
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open file %s with error %s", path, err)
		return p, err
	}
	bits, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Failed to read json bytes %s", err)
		return p, err
	}
	err = json.Unmarshal(bits, &p)
	return p, err

}

// Get a paginated list of the history of triggers for the project
func History(offset, size int) []Trigger {
	return make([]Trigger, 0)
}
