package ext

import (
	"log"
	"plugin"
	"sync"

	"github.com/lateefj/trigr"
)

type Extension struct {
	Name   string
	Symbol plugin.Symbol
}

func (e *Extension) Handle(event *trigr.Trigger) {
	e.Symbol.(func(evt *trigr.Trigger))(event)
}

type Registrar struct {
	Extensions map[string]Extension
	Mutex      sync.RWMutex
}

func (r *Registrar) Add(extName string, pluginPath string) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}
	s, err := p.Lookup("TriggerHandler")
	if err != nil {
		return err
	}
	defer r.Mutex.Unlock()
	r.Mutex.Lock()
	r.Extensions[extName] = Extension{Name: extName, Symbol: s}
}

func (r *Registrar) Publish(event *trigr.Trigger) {
	defer r.Mutex.RUnlock()
	r.Mutex.RLock()
	for k, p := range r.Extensions {
		log.Printf("Calling handle for plugin %s\n", k)
		p.Handle(event)
	}
}

var reg Registrar

func init() {
	reg = Registrar{Extensions: make(map[string]Extension)}
}
