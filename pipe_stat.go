package main

import (
	"fmt"
	"os"
	"strconv"

	sc "github.com/aronfan/shmcore"
	"github.com/aronfan/xerrors"
)

type bstat struct {
	bytes uint32
	count uint32
	frees uint32
}

func (pc *pipecmd) stat() error {
	ok := sc.ResumeEnabled()
	if !ok {
		return xerrors.Wrap(fmt.Errorf("resume not enabled")).WithInt(-2)
	}
	s, ok := pc.params["shmkey"]
	if !ok {
		return xerrors.Wrap(fmt.Errorf("shmkey not exist")).WithInt(-2)
	}
	if s == "" {
		return xerrors.Wrap(fmt.Errorf("shmkey is empty")).WithInt(-2)
	}
	k, err := strconv.Atoi(s)
	if err != nil {
		return xerrors.Wrap(err)
	}

	shmkey := uint32(k)
	err = sc.IsShmExist(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}

	seg, err := sc.NewSegment(shmkey, 0)
	if err != nil {
		return xerrors.Wrap(err)
	}

	err = seg.Attach()
	if err != nil {
		return xerrors.Wrap(err)
	}
	defer seg.Detach()

	m := make(map[uint16]*bstat)
	seg.Observe(
		func(shead *sc.SegmentHead) {
		},
		func(index uint16, bhead *sc.BucketHead) {
			m[index] = &bstat{bytes: bhead.GetBytes(), count: bhead.GetCount(), frees: 0}
		},
		func(hindex uint16, uindex uint32, unit *sc.BucketUnit) {
			if unit.GetLen() == 0 {
				stat, ok := m[hindex]
				if ok {
					stat.frees++
				}
			}
		},
	)

	fmt.Fprintf(os.Stderr, "shmkey=%d\n", shmkey)
	fmt.Fprintln(os.Stderr, "bytes\tfree\ttotal")
	for _, v := range m {
		fmt.Fprintf(os.Stdout, "%d\t%d\t%d\n", v.bytes, v.frees, v.count)
	}

	return nil
}
