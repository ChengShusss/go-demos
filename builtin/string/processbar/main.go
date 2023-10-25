package main

import (
	"fmt"
	"strings"
	"time"
)

type ProcessBar struct {
	length int
}

func NewProcessBar(n int) *ProcessBar {
	return &ProcessBar{
		length: n,
	}
}

func (bar *ProcessBar) ShowBar(p int) {
	if p < 0 || p > 100 {
		return
	}

	backs := strings.Repeat("\b", bar.length)
	process := strings.Repeat("=", (bar.length-7)*p/100) + ">"
	format := fmt.Sprintf("%%s[%%-%ds%%3d%%%%]", bar.length-6)
	fmt.Printf(format, backs, process, p)
}

func main() {
	bar := NewProcessBar(20)
	for i := 0; i < 100; i += 2 {
		bar.ShowBar(i)
		time.Sleep(time.Second / 20)
	}
	bar.ShowBar(100)
}
