package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/lateefj/trigr"
)

var (
	// Ignore directories that are source code
	scmDirs = []string{".git", ".hg", ".cvs", ".svn"}
)

// duplicateLimiter ... rate limit the number of changes for a single file
type duplicateLimiter struct {
	files       sync.Map
	changeDelay time.Duration
}

// Add file if not duplicate
func (dl *duplicateLimiter) add(name string) bool {
	// Check to see if it exists
	expire, exists := dl.files.Load(name)

	// If has a key but it has expired
	if exists && expire.(time.Time).Before(time.Now()) {
		exists = false
		dl.files.Delete(name)
	}
	if !exists {
		// Doesn't exist then add it
		dl.files.Store(name, time.Now().Add(dl.changeDelay))
		return true
	}
	return false
}

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
	limiter        *duplicateLimiter
	stopChannel    chan bool
}

// NewDirectoryWatcher ... Watches a directory publish file change events
func NewDirectoryWatcher(path string, trigChan chan *trigr.Trigger, excludeSCM bool) *DirectoryWatcher {
	return &DirectoryWatcher{Path: path, TriggerChannel: trigChan, ExcludeSCM: excludeSCM, limiter: &duplicateLimiter{changeDelay: 1000 * time.Millisecond}}
}

func (dw *DirectoryWatcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("ERROR: notify file %s\n", err)
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
			if dw.limiter.add(ev.Name) {
				dw.TriggerChannel <- t
			}
		case err := <-watcher.Errors:
			log.Printf("ERROR: %s\n", err)
		case <-dw.stopChannel:
			log.Printf("Stopping to watch %s\n", dw.Path)
			return nil
		}
	}
}

func (dw *DirectoryWatcher) Stop() {
	dw.stopChannel <- true
}
