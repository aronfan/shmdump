package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTms(t *testing.T) {
	shmkey := uint32(255)
	base := time.Now().Format("20060102_150405")
	file := fmt.Sprintf("%s.SHM%d", base, shmkey)
	t.Logf("%s", file)
}
