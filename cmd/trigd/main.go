package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
)

var (
	// Track total number of messages
	messageCount int64 = 0
	confFile           = flag.String("conf", "", "Path to configuration file")
)

func init() {
	home := os.Getenv("HOME")
	if *confFile == "" && home != "" {
		defaultFile := fmt.Sprintf("%s/.trigr/config.json", home)
		confFile = &defaultFile
	}
}

func setOutput() string {

	cmd := exec.Command("bash", "-c", "set")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("FAILED WITH ERROR: %s\n", err)
	}
	return string(output)
}

func handleTrigger(env []string, p *Project, t *trigr.Trigger) {
	in := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	if p.LocalSource != nil {
		path := p.LocalSource.Path
		luaPath := fmt.Sprintf("%s/.trigr/%s.lua", path, t.Type)
		if _, err := os.Stat(luaPath); err == nil {
			l := ext.NewTrigSL(in, out, "./lsl/lua")
			l.SetGlobalVar("exec", func(cmd, directory string) string {
				split := strings.Split(cmd, " ")
				c := split[0]
				args := make([]string, 0)
				if len(split) > 1 {
					args = split[1:]
				}
				p := exec.Command(c, args...)
				p.Dir = directory
				p.Env = env
				t.Logs <- trigr.NewLog(fmt.Sprintf("running: %s ", cmd))
				output, err := p.CombinedOutput()
				if err != nil {
					t.Logs <- trigr.NewLog(err.Error())
				}
				return string(output)
			})
			// Add the trig event to the context
			l.SetGlobalVar("trig", t)
			err = l.RunFile(luaPath, t, make(chan *trigr.Trigger))
			if err != nil {
				msgErr := fmt.Sprintf("Failed to run dsl %s\n", err)
				t.Logs <- trigr.NewLog(msgErr)
				log.Print(msgErr)
			}
		}
	}
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func setAppend(items []string, e string) []string {
	if !contains(items, e) {
		items = append(items, e)
	}
	return items
}

func main() {
	flag.Parse()
	log.Printf("Starting trigd\n")
	if _, err := os.Stat(*confFile); err == nil {
		log.Printf("Config file path %s\n", *confFile)
		bits, err := ioutil.ReadFile(*confFile)
		if err != nil {
			log.Printf("ERROR: Reading file %s error: %s\n", *confFile, err)
		}
		pm, err := LoadProjectManager(bits)
		if err != nil {
			log.Printf("ERROR: Loading file %s error: %s\n", *confFile, err)
		}
		Manager = pm
		// Append the current env to the global saved one
		for _, e := range os.Environ() {
			Manager.GlobalEnv = setAppend(Manager.GlobalEnv, e)
		}
	} else {
		log.Printf("ERROR: Configuration file does not exist: %s\n", *confFile)
	}
	// Default watch current directory
	//dw := NewDirectoryWatcher("./", TriggerChannel, true)
	//go dw.Watch()
	setupHandlers()
	log.Printf("Handled %d messages\n", messageCount)
}
