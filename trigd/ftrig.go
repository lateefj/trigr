// TODO: Recursive watch needs doing

package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lateefj/trigr"
	"golang.org/x/exp/inotify"
)

type DirectoryWatcher struct {
	Path           string
	TriggerChannel chan trigr.Trigger
}

func (dw *DirectoryWatcher) Watch() error {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Error(err)
		return err
	}
	err = watcher.Watch(dw.Path)
	if err != nil {
		log.Error(err)
		return err
	}
	for {
		select {
		case ev := <-watcher.Event:
			if ev.Mask == inotify.IN_CLOSE_WRITE {
				d := map[string]interface{}{
					"path": fmt.Sprintf("%s%s", dw.Path, ev.Name),
				}
				t := trigr.NewTrigger("file", d)

				dw.TriggerChannel <- t
				log.Debugf("Write event: %v", ev)
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}

	return nil

}
