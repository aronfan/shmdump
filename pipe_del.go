package main

import (
	"fmt"
	"strconv"

	sc "github.com/aronfan/shmcore"
	"github.com/aronfan/xerrors"
)

func (pc *pipecmd) del() error {
	ok := sc.CanObserve()
	if !ok {
		return nil
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

	// ensure the Natt is 0
	ds, err := sc.GetShmDsByKey(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}
	if ds.Natt > 0 {
		return xerrors.Wrap(fmt.Errorf("natt=%d, should shutdown app attached", ds.Natt)).WithInt(-2)
	}

	err = sc.DelShmByKey(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}

	return nil
}
