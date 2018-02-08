package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/lateefj/trigr"
	"github.com/lateefj/trigr/ext"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func main() {
	luaPath := os.Args[1]
	log.Debugf("Luapath is %s", luaPath)
	in, out, err := os.Pipe()
	if err != nil {
		log.Errorf("Failed to get a pipe %s", err)
		return
	}
	log.Debugf("Lua loading file %s", luaPath)
	t := trigr.NewTrigger("test", make(map[string]interface{}))
	if _, err := os.Stat(luaPath); err == nil {

		l := ext.NewLuaDslLoader(in, out, "./ext/lua")
		err = l.RunDslFile(luaPath, &t, make(chan *trigr.Trigger))
		if err != nil {
			log.Errorf("Failed ot run dsl %s", err)
		}
	}
}
