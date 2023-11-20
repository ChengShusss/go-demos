package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

var (
	modKernelBase            = windows.NewLazySystemDLL("Kernelbase.dll")
	procCompareObjectHandles = modKernelBase.NewProc("CompareObjectHandles")

	ShellList = map[string]bool{
		"cmd.exe":        true,
		"powershell.exe": true,
	}
)

func main() {
	// Demo()

	start := time.Now()
	var pids []int32
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			pid, err := strconv.ParseInt(arg, 10, 32)
			if err != nil {
				fmt.Printf("Invalid pid: %v\n", arg)
				os.Exit(1)
			}
			pids = append(pids, int32(pid))
		}
	}

	processes, err := LoadProcess(pids)
	if err != nil {
		fmt.Printf("Failed to get process, err: %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Len of process: %v\n", processes)
	DetectReverseShell(processes)

	fmt.Printf("consumed: %v ms\n", time.Since(start).Milliseconds())
}

func DetectReverseShell(processes []*process.Process) {
	for _, process := range processes {
		exe, err := process.Exe()
		if err != nil {
			// fmt.Printf("Failed to get process %v, err: %v\n", process.Pid, err)
			continue
		}
		if !isShell(exe) {
			// fmt.Printf("Process %v is not shell, continue\n", process.Pid)
			continue
		}

		isReverse := checkProcessTree(process)
		if isReverse {
			fmt.Printf("%v, - %v\n", process.Pid, isReverse)
		}
	}
}

func checkProcessTree(current *process.Process) bool {

	var nodeList []*process.Process
	nodeFilter := map[uintptr]bool{uintptr(current.Pid): true}

	parent := current
	var err error
	for {
		parent, err = parent.Parent()
		if err != nil {
			break
		}
		nodeList = append(nodeList, parent)
		nodeFilter[uintptr(parent.Pid)] = true
	}

	pidHandleMap := TransHandle(GetHandlesForProcess(nodeFilter))

	for _, node := range nodeList {
		cond1 := DetectHandle(
			pidHandleMap[uintptr(node.Pid)],
			pidHandleMap[uintptr(current.Pid)],
		)
		conns, err := node.Connections()
		if err != nil {
			continue
		}
		cond2 := len(conns) > 0
		if cond1 && cond2 {
			for _, c := range conns {
				fmt.Printf("%+v\n", c)
			}
			return true
		}
	}

	return false
}

// 检查两个进程是否有共用的pipe或socket对象
func DetectHandle(h1, h2 []windows.Handle) bool {

	for _, handle1 := range h1 {

		type1, err := windows.GetFileType(handle1)
		if err != nil {
			continue
		}
		if type1 != 3 && type1 != 2 {
			continue
		}

		for _, handle2 := range h2 {

			type2, err := windows.GetFileType(handle2)
			if err != nil {
				continue
			}
			if type2 != 3 && type2 != 2 {
				continue
			}

			if CompareObjectHandles(handle1, handle2) {
				return true
			}
		}
	}

	return false
}

func CompareObjectHandles(handle1, handle2 windows.Handle) bool {
	result, _, _ := syscall.SyscallN(procCompareObjectHandles.Addr(), uintptr(handle1), uintptr(handle2))
	return result == 1
}

func LoadProcess(pids []int32) ([]*process.Process, error) {

	if len(pids) == 0 {
		processes, err := process.Processes()
		return processes, err
	}

	var processes []*process.Process
	for _, pid := range pids {
		process, err := process.NewProcess(pid)
		if err != nil {
			return nil, err
		}
		processes = append(processes, process)
	}

	return processes, nil
}

func isShell(path string) bool {
	exe := filepath.Base(path)
	return ShellList[strings.ToLower(exe)]
}
