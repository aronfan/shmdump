package main

import (
	"fmt"
	"testing"
	"time"
)

const testkey = uint32(256)

func TestTms(t *testing.T) {
	shmkey := testkey
	base := time.Now().Format("20060102_150405")
	file := fmt.Sprintf("%s.SHM%d", base, shmkey)
	t.Logf("%s", file)
}

func TestStat(t *testing.T) {
	cmd := newPipeCommand("shmkey=256&op=stat")
	t.Logf("%+v", cmd.params)
	if err := cmd.dispatch(); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestSave(t *testing.T) {
	cmd := newPipeCommand("shmkey=256&op=save&verbose=1&file=sample.shm")
	t.Logf("%+v", cmd.params)
	if err := cmd.dispatch(); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestLoad(t *testing.T) {
	name := "sample.shm"
	t.Logf("%s", name)
}
