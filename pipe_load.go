package main

import (
	"fmt"
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
	return nil
}

type sampleCfg struct {
	ShmKV kv.DBCfg `xml:"shmkv"`
}
