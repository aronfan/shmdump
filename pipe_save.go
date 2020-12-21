package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	sc "github.com/aronfan/shmcore"
	kv "github.com/aronfan/shmkv"
	"github.com/aronfan/xerrors"
)

const (
	tagHead = "MF20"
	tagFoot = "MF20"
	tagLine = "\r\n"
)

func (pc *pipecmd) save() error {
	ok := sc.CanResume()
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

	// ensure the Natt is 0
	ds, err := sc.GetShmDsByKey(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}
	if ds.Natt > 0 {
		return xerrors.Wrap(fmt.Errorf("natt=%d, should shutdown app attached", ds.Natt)).WithInt(-2)
	}

	db := kv.NewDB()
	err = db.Resume(shmkey)
	if err != nil {
		return xerrors.Wrap(err)
	}

	// collect the units with length
	var units unitArray = nil
	db.Visit(func(_ string, h uint16, u uint32, unit *sc.BucketUnit) bool {
		ul := unitLen{H: h, U: u, L: unit.GetLen()}
		units = append(units, ul)
		return true
	})
	sort.Sort(units)
	if true {
		for i := 0; i < len(units); i++ {
			unit := units[i]
			fmt.Fprintf(os.Stderr, "hindex=%d, uindex=%d, length=%d\n", unit.H, unit.U, unit.L)
		}
	}

	// save
	file, ok := pc.params["file"]
	if !ok {
		base := time.Now().Format("20060102_150405")
		file = fmt.Sprintf("%s.SHM%d", base, shmkey)
	}

	fmt.Fprintf(os.Stdout, "%s\n", file)

	return nil
}

type unitLen struct {
	H uint16
	U uint32
	L uint32
}

type unitArray []unitLen

func (s unitArray) Len() int           { return len(s) }
func (s unitArray) Less(i, j int) bool { return s[i].L < s[j].L }
func (s unitArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
