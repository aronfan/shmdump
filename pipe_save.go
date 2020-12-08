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

	s, ok = pc.params["key"]
	if !ok {
		return xerrors.Wrap(fmt.Errorf("key not exist")).WithInt(-2)
	}
	if s == "" {
		return xerrors.Wrap(fmt.Errorf("key is empty")).WithInt(-2)
	}

	file, ok := pc.params["file"]
	if !ok {
		base := time.Now().Format("20060102_150405")
		file = fmt.Sprintf("%s.SHM%d", base, shmkey)
	}

	saveall := false
	if s == "*" {
		saveall = true
	}

	if saveall {

	} else {

	}

	fmt.Fprintf(os.Stdout, "%s\n", file)

	return nil
}
