package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	// Track total number of messages
	messageCount int64 = 0
	confFile           = flag.String("conf", "", "Path to configuration file")
)

// Read in configuration file if it exists
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
