package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/lateefj/trigr"
)

var (
	Manager = NewProjectManager()
)

type Connected struct {
	outbound  map[int32]chan []byte
	Inbound   chan string
	Lock      sync.RWMutex
	currentId int32
}

func NewConnected() *Connected {
	return &Connected{make(map[int32]chan []byte, 0), make(chan string), sync.RWMutex{}, 1}
}

func (c *Connected) New() (int32, chan []byte) {

	out := make(chan []byte, 10)

	c.Lock.Lock()
	id := atomic.AddInt32(&c.currentId, 1)
	c.outbound[id] = out
	c.Lock.Unlock()

	return id, out

}

func (c *Connected) Remove(id int32) {
	c.Lock.Lock()
	delete(c.outbound, id)
	c.Lock.Unlock()
}

func (c *Connected) Send(m []byte) {
	c.Lock.RLock()
	for _, out := range c.outbound {
		out <- m
	}
	c.Lock.RUnlock()
}

type Project struct {
	trigr.Project
	Connected *Connected
}

func NewProject(id string) *Project {
	return &Project{Project: trigr.Project{Id: id, Triggers: make(chan *trigr.Trigger), Logs: make(chan *trigr.Log)}, Connected: NewConnected()}
}

func (p *Project) Send(m []byte) {
	p.Connected.Send(m)
}

func (p *Project) MonitorDirectory(path string) {

	dw := NewDirectoryWatcher(path, p.Triggers, true)
	go dw.Watch()

	for t := range p.Triggers {

		// First send the trigger out clients
		b, err := json.Marshal(t)
		if err != nil {
			log.Printf("Failed to marshal trigger %s\n", err)
			return
		}
		p.Connected.Send(b)
		// Next send logs to clients
		go func() {
			for l := range t.Logs {
				b, err := json.Marshal(l)
				if err != nil {
					log.Printf("Failed to trigger log %v error:  %s\n", l, err)
					continue
				}
				p.Connected.Send(b)
			}
		}()
		handleTrigger(t)
	}
}

type ProjectManager struct {
	Projects map[string]*Project
	mutex    sync.RWMutex
}

func NewProjectManager() *ProjectManager {
	return &ProjectManager{Projects: make(map[string]*Project), mutex: sync.RWMutex{}}
}

func (pm *ProjectManager) Add(p *Project) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if _, exists := pm.Projects[p.Id]; exists {
		return errors.New(fmt.Sprintf("Project %s already exists", p.Id))
	}
	pm.Projects[p.Id] = p
	return nil
}

func (pm *ProjectManager) Remove(id string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if _, exists := pm.Projects[id]; !exists {
		return errors.New(fmt.Sprintf("Project %s doesn't exist", id))
	}
	delete(pm.Projects, id)
	return nil
}

func (pm *ProjectManager) Get(id string) (*Project, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if p, exists := pm.Projects[id]; !exists {
		return nil, errors.New(fmt.Sprintf("Project %s doesn't exist", id))
	} else {
		return p, nil
	}
}
