package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sc "github.com/aronfan/shmcore"
	"github.com/aronfan/xerrors"
)

func (pc *pipecmd) save() error {
	ok := sc.CanObserve()
	if !ok {
		return xerrors.Wrap(fmt.Errorf("observe not enabled")).WithInt(-2)
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
	err = sc.Exist(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}

	file, ok := pc.params["file"]
	if !ok {
		base := time.Now().Format("20060102_150405")
		file = fmt.Sprintf("%s.SHM%d", base, shmkey)
	}

	err = sc.Exist(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}

	// ensure the Natt is 0
	ds, err := sc.GetShmDsByKey(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}
	if ds.Natt > 0 {
		return xerrors.Wrap(fmt.Errorf("natt=%d, should shutdown app attached", ds.Natt)).WithInt(-2)
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

	// save

	fmt.Fprintf(os.Stdout, "%s\n", file)

	return nil
}
