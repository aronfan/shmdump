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

func TestTms(t *testing.T) {
	shmkey := uint32(255)
	base := time.Now().Format("20060102_150405")
	file := fmt.Sprintf("%s.SHM%d", base, shmkey)
	t.Logf("%s", file)
}

func TestStat(t *testing.T) {
	cmd := newPipeCommand("shmkey=255&op=stat")
	t.Logf("%+v", cmd.params)
	if err := cmd.dispatch(); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
}

func TestSave(t *testing.T) {
	shmkey := uint32(255)
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
	c := 0

	// write the head
	n, err := w.WriteString(tagHead)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	c += n

	failed := false
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
				_, err := w.WriteString(string(kvLen[:]))
				if err != nil {
					t.Errorf("%s", err.Error())
					failed = true
					return
				}
				_, err = w.WriteString(string(unit.GetReadSlice()[:]))
				if err != nil {
					t.Errorf("%s", err.Error())
					failed = true
					return
				}
			}
		},
	)

	if !failed {
		return
	}

	// write the foot
	n, err = w.WriteString(tagFoot)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	c += n

	w.Flush()
	t.Logf("%d", c)
}
