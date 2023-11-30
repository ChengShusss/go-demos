package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chengshusss/go-demos/windows/reverseShellDetectV2/reverseshell"
)

func main() {
	// Demo()

	start := time.Now()
	defer func() { fmt.Printf("consumed: %v ms\n", time.Since(start).Milliseconds()) }()
	var pids []int32
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			pid, err := strconv.ParseInt(arg, 10, 32)
			if err != nil {
				fmt.Printf("Invalid pid: %v\n", arg)
				os.Exit(1)
			}
			if reverseshell.DetectByPid(int32(pid)) {
				fmt.Printf("%d is reverse shell\n", pid)
			}
		}
		return
	}

	processes, err := reverseshell.LoadProcess(pids)
	if err != nil {
		fmt.Printf("Failed to get process, err: %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Len of process: %v\n", processes)
	reverseshell.DetectProcesses(processes)
}
