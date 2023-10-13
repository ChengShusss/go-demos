package main

import (
	"os"
	"syscall"
	"testing"
)

func TestFilePerm(t *testing.T) {

	mask := syscall.Umask(0022)
	defer syscall.Umask(mask)
	err := os.MkdirAll("test", os.ModePerm)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

}
