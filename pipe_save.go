package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sc "github.com/aronfan/shmcore"
)

func (pc *pipecmd) save() error {
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

	s, ok = pc.params["key"]
	if !ok {
		return fmt.Errorf("key not exist")
	}
	if s == "" {
		return fmt.Errorf("key is empty")
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
