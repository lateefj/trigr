package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/lateefj/trigr"
)

var (
	// Ignore directories
	scmDirs = []string{".git", ".hg", ".cvs", ".svn"}
)

// Helper to match any source control paths
func isSCMPath(path string) bool {
	for _, d := range scmDirs {
		if strings.Contains(path, d) {
			return true
		}
	}
	return false
}

// Way to monitor file system events
type DirectoryWatcher struct {
	Path           string
	TriggerChannel chan *trigr.Trigger
	ExcludeSCM     bool
}

func NewDirectoryWatcher(path string, trigChan chan *trigr.Trigger, excludeSCM bool) *DirectoryWatcher {
	return &DirectoryWatcher{Path: path, TriggerChannel: trigChan, ExcludeSCM: excludeSCM}
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
		// Exclude source control management
		if dw.ExcludeSCM && isSCMPath(newPath) {
			return nil
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
				"path": fmt.Sprintf("%s", ev.Name),
			}
			if ev.Op == fsnotify.Write {
				d["op"] = "write"
			} else if ev.Op == fsnotify.Create {
				d["op"] = "create"
			} else if ev.Op == fsnotify.Remove {
				d["op"] = "remove"
			} else if ev.Op == fsnotify.Chmod {
				d["op"] = "chmod"
			} else if ev.Op == fsnotify.Rename {
				d["op"] = "rename"
			}
			t := trigr.NewTrigger("file", d)
			dw.TriggerChannel <- t
			//log.Printf("Write event: %v\n", ev)
		case err := <-watcher.Errors:
			log.Printf("ERROR: %s\n", err)
		}
	}
	return nil
}
