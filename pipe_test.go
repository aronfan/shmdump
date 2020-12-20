package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"

	sc "github.com/aronfan/shmcore"
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
	shmkey := testkey
	base := time.Now().Format("20060102_150405")
	name := fmt.Sprintf("%s.SHM%d", base, shmkey)

	err := sc.Exist(shmkey)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	seg, err := sc.NewSegment(shmkey, 0)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	err = seg.Attach()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	defer seg.Detach()

	f, err := os.Create(name)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// write the head
	_, err = w.WriteString(tagHead + tagLine)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	err = w.Flush()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	ok := true
	kvLen := make([]byte, 4)
	seg.Observe(
		func(_ *sc.SegmentHead) {
		},
		func(_ uint16, _ *sc.BucketHead) {
		},
		func(_ uint16, _ uint32, unit *sc.BucketUnit) {
			uLen := unit.GetLen()
			if uLen > 0 {
				binary.BigEndian.PutUint32(kvLen, uLen)
				_, err = w.WriteString(string(kvLen[:]))
				if err != nil {
					t.Errorf("%s", err.Error())
					ok = false
					return
				}
				_, err = w.WriteString(string(unit.GetReadSlice()[:]))
				if err != nil {
					t.Errorf("%s", err.Error())
					ok = false
					return
				}
			}
		},
	)

	if !ok {
		return
	}

	err = w.Flush()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	// write the foot
	_, err = w.WriteString(tagLine + tagFoot)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	err = w.Flush()
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestLoad(t *testing.T) {
	name := "sample.SHM256"
	t.Logf("%s", name)
}
