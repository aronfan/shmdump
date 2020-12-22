package main

import (
	"fmt"
	"testing"

	sc "github.com/aronfan/shmcore"
)

const (
	testkey  = uint32(256)
	loadfile = "sample.shm"
	savefile = "sample2.shm"
)

func TestPipe(t *testing.T) {
	cmd := newPipeCommand(fmt.Sprintf("shmkey=%d&op=del", testkey))
	t.Logf("%+v", cmd.params)
	if err := sc.Exist(testkey); err == nil {
		if err := cmd.dispatch(); err != nil {
			t.Errorf("%s", err.Error())
			return
		}
	}

	cmd = newPipeCommand(fmt.Sprintf("shmkey=%d&op=load&file=%s&cfg=sample.xml", testkey, loadfile))
	t.Logf("%+v", cmd.params)
	if err := cmd.dispatch(); err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	cmd = newPipeCommand(fmt.Sprintf("shmkey=%d&op=save&verbose=1&file=%s", testkey, savefile))
	t.Logf("%+v", cmd.params)
	if err := cmd.dispatch(); err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	cmd = newPipeCommand(fmt.Sprintf("shmkey=%d&op=stat", testkey))
	t.Logf("%+v", cmd.params)
	if err := cmd.dispatch(); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}
