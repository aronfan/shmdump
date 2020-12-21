package main

import (
	"fmt"
	"strings"
)

type pipecmd struct {
	params map[string]string
}

func newPipeCommand(cmd string) *pipecmd {
	ss := strings.Split(cmd, "&")
	params := make(map[string]string)
	for i := 0; i < len(ss); i++ {
		s := ss[i]
		kv := strings.Split(s, "=")
		if len(kv) >= 2 {
			params[kv[0]] = kv[1]
		} else {
			params[kv[0]] = ""
		}
	}
	return &pipecmd{params: params}
}

func (pc *pipecmd) dispatch() error {
	op, ok := pc.params["op"]
	if !ok {
		return fmt.Errorf("op not exist")
	}

	switch op {
	case "stat":
		return pc.stat()
	case "save":
		return pc.save()
	case "del":
		return pc.del()
	case "load":
		return pc.load()
	}
	return nil
}

func pipe(cmd string) error {
	return newPipeCommand(cmd).dispatch()
}
