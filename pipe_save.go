package main

import (
	"bufio"
	"encoding/binary"
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
	tagHead = "MF20\r\n"
	tagFoot = "\r\nMF20"
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
	db.Visit(func(k string, h uint16, u uint32, unit *sc.BucketUnit) bool {
		ul := unitLen{K: k, H: h, U: u, L: unit.GetLen()}
		units = append(units, ul)
		return true
	})
	sort.Sort(units)

	// save
	file, ok := pc.params["file"]
	if !ok {
		base := time.Now().Format("20060102_150405")
		file = fmt.Sprintf("%s.shm%d", base, shmkey)
	}
	_, ok = pc.params["verbose"]
	if ok {
		return save(file, units, db, true)
	} else {
		return save(file, units, db, false)
	}
}

func save(file string, units unitArray, db *kv.DB, verbose bool) error {
	fmt.Fprintf(os.Stdout, "%s\n", file)
	f, err := os.Create(file)
	if err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// write the head
	_, err = w.WriteString(tagHead)
	if err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}

	seg := db.GetSegment()
	kvLen := make([]byte, 4)
	for i := 0; i < len(units); i++ {
		t := units[i]
		unit, err := seg.GetBucketUnit(t.H, t.U)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "get from db failed: %s\n", err.Error())
			}
			continue
		}
		s := unit.GetReadSlice()
		length := uint32(len(s))
		binary.BigEndian.PutUint32(kvLen, length)
		w.WriteString(string(kvLen[:]))
		_, err = w.WriteString(string(s[:]))
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "write failed: %s\n", err.Error())
			}
			return xerrors.Wrap(err).WithInt(-2)
		} else {
			if verbose {
				fmt.Fprintf(os.Stderr, "key=%s h=%d u=%d len=%d write ok\n", t.K, t.H, t.U, length)
			}
		}
	}

	// write the foot
	_, err = w.WriteString(tagFoot)
	if err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}

	err = w.Flush()
	if err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}

	return nil
}

type unitLen struct {
	K string
	H uint16
	U uint32
	L uint32
}

type unitArray []unitLen

func (s unitArray) Len() int           { return len(s) }
func (s unitArray) Less(i, j int) bool { return s[i].L < s[j].L }
func (s unitArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
