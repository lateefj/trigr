package main

import (
	"sync"
	"sync/atomic"
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

var (
	ClientsConnected *Connected
)

func init() {
	ClientsConnected = NewConnected()
}
