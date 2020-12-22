package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

	sc "github.com/aronfan/shmcore"
	kv "github.com/aronfan/shmkv"
	"github.com/aronfan/xerrors"
)

func (pc *pipecmd) load() error {
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

	if err = sc.Exist(shmkey); err == nil {
		return xerrors.Wrap(fmt.Errorf("shm %d already exist", shmkey)).WithInt(-2)
	}

	s, ok = pc.params["cfg"]
	if !ok {
		return xerrors.Wrap(fmt.Errorf("cfg not exist")).WithInt(-2)
	}

	var cfg sampleCfg
	if err = kv.ParseXML(s, &cfg); err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}
	opt := &sc.SegmentOption{}
	for i := 0; i < len(cfg.ShmKV.Buckets); i++ {
		c := cfg.ShmKV.Buckets[i]
		opt.AddBucket(c.Count, c.Bytes)
	}

	db := kv.NewDB()
	if err = db.Init(shmkey, opt); err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}
	defer db.Fini()

	s, ok = pc.params["file"]
	if !ok {
		return xerrors.Wrap(fmt.Errorf("file not exist")).WithInt(-2)
	}

	// load
	if err = load(s, db); err != nil {
		return xerrors.Wrap(err).WithInt(-2)
	}

	return nil
}

func load(fileName string, db *kv.DB) error {
	f, err := os.Open(fileName)
	if err != nil {
		return xerrors.Wrap(err)
	}
	defer f.Close()

	// read the head and validate
	headLen := int64(len(tagHead))
	head := make([]byte, headLen)
	_, err = f.Read(head)
	if err != nil {
		return xerrors.Wrap(err)
	}
	if string(head[:]) != tagHead {
		return xerrors.Wrap(fmt.Errorf("head=%s, not equal %s", string(head[:]), tagHead)).WithInt(-2)
	}

	// read the foot and validate
	footLen := int64(len(tagFoot))
	foot := make([]byte, footLen)
	offsetEnd, err := f.Seek(-footLen, os.SEEK_END)
	if err != nil {
		return xerrors.Wrap(err)
	}
	_, err = f.Read(foot)
	if string(foot[:]) != tagFoot {
		return xerrors.Wrap(fmt.Errorf("foot=%s, not equal %s", string(foot[:]), tagFoot)).WithInt(-2)
	}

	// read the units
	f.Seek(headLen, os.SEEK_SET)
	offsetBegin, err := f.Seek(0, os.SEEK_CUR)

	curr := offsetBegin
	kvLen := make([]byte, 4)
	for {
		if curr+4 >= offsetEnd {
			break
		}
		f.Read(kvLen)
		unitLen := binary.BigEndian.Uint32(kvLen)
		if curr+int64(unitLen) >= offsetEnd {
			break
		}
		if unitLen == 0 {
			continue
		}
		unit := make([]byte, unitLen)
		_, err = f.Read(unit)
		if err != nil {
			return xerrors.Wrap(err)
		}
		db.SetWithSlice(unit)
	}

	return nil
}

type sampleCfg struct {
	ShmKV kv.DBCfg `xml:"shmkv"`
}
