package main

import (
	"fmt"
	"os"
	"strconv"

	sc "github.com/aronfan/shmcore"
)

type bstat struct {
	bytes uint32
	count uint32
	frees uint32
}

func (pc *pipecmd) stat() error {
	ok := sc.ResumeEnabled()
	if !ok {
		return fmt.Errorf("resume not enabled")
	}
	s, ok := pc.params["shmkey"]
	if !ok {
		return fmt.Errorf("shmkey not exist")
	}
	if s == "" {
		return fmt.Errorf("shmkey is empty")
	}
	k, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	shmkey := uint32(k)
	err = sc.IsShmExist(shmkey)
	if err != nil {
		return err
	}

	seg, err := sc.NewSegment(shmkey, 0)
	if err != nil {
		return err
	}

	err = seg.Attach()
	if err != nil {
		return err
	}

	m := make(map[uint16]*bstat)
	seg.Observe(
		func(shead *sc.SegmentHead) {
		},
		func(index uint16, bhead *sc.BucketHead) {
			m[index] = &bstat{bytes: bhead.GetBytes(), count: bhead.GetCount(), frees: 0}
		},
		func(hindex uint16, uindex uint32, unit *sc.BucketUnit) {
			l := unit.GetLen()
			if l == 0 {
				stat, ok := m[hindex]
				if ok {
					stat.frees++
				}
			}
		},
	)

	seg.Detach()

	fmt.Fprintf(os.Stderr, "shmkey=%d\n", shmkey)
	fmt.Fprintln(os.Stderr, "bytes\tfree\ttotal")
	for _, v := range m {
		fmt.Fprintf(os.Stdout, "%d\t%d\t%d\n", v.bytes, v.frees, v.count)
	}

	return nil
}
