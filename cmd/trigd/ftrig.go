package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/lateefj/trigr"
)

type DirectoryWatcher struct {
	Path           string
	TriggerChannel chan *trigr.Trigger
}

func (dw *DirectoryWatcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("ERROR: from notify file %s\n", err)
		return err
	}
	filepath.Walk(dw.Path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = watcher.Add(newPath)
			if err != nil {
				log.Printf("ERROR: adding path %s %s", newPath, err)
				return err
			}
			fmt.Printf("Watching path %s\n", newPath)
		}
		return nil
	})
	for {
		select {
		case ev := <-watcher.Events:
			d := map[string]interface{}{
				"path": fmt.Sprintf("%s%s", dw.Path, ev.Name),
			}
			if ev.Op == fsnotify.Write {
				d["op"] = "write"
			} else if ev.Op == fsnotify.Create {
				d["op"] = "create"
			} else if ev.Op == fsnotify.Remove {
				d["op"] = "remove"
			}
			t := trigr.NewTrigger("file", d)
			dw.TriggerChannel <- t
			log.Printf("Write event: %v\n", ev)
		case err := <-watcher.Errors:
			log.Printf("ERROR: %s\n", err)
		}
	}
	return nil
}
