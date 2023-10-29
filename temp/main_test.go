package main

import (
	"log"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestFilePerm(t *testing.T) {

	mask := syscall.Umask(0022)
	defer syscall.Umask(mask)
	err := os.MkdirAll("test", os.ModePerm)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

}

var (
	UpdateWaitTime = time.Second * 5
)

type Save struct {
	updateLock sync.Mutex
	updateCh   chan struct{}
}

func (s *Save) DelaySave() {
	// TODO Persistence
	ticker := time.NewTicker(UpdateWaitTime)

	s.updateLock.Lock()
	if s.updateCh != nil {
		close(s.updateCh)
	}
	s.updateCh = make(chan struct{})
	s.updateLock.Unlock()

	select {
	case <-ticker.C:
		// TODO update
		log.Printf("ATP - sensitiveFile: start to write feature into minio\n")
	case <-s.updateCh:
		log.Printf("ATP - sensitiveFile: receive new update, return directly\n")
		return
	}

}

func TestDelaySave(t *testing.T) {
	s := Save{}

	go s.DelaySave()
	time.Sleep(1 * time.Second)
	go s.DelaySave()
	time.Sleep(1 * time.Second)
	go s.DelaySave()

	time.Sleep(10 * time.Second)
}
