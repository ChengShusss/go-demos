package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		str := strings.Repeat("\b", 105)
		process := strings.Repeat("=", i) + ">"
		fmt.Printf("%s[%-100s%2d%%]", str, process, i)
		time.Sleep(time.Second / 20)
	}
}
