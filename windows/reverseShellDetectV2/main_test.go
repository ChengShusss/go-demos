package main

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

func TestCompate(t *testing.T) {
	// 获取进程句柄
	handle1, err := windows.OpenProcess(windows.PROCESS_DUP_HANDLE, false, 21456)
	if err != nil {
		fmt.Println("Error opening process handle 1:", err)
		return
	}
	defer windows.CloseHandle(handle1)

	handle2, err := windows.OpenProcess(windows.PROCESS_DUP_HANDLE, false, 18332)
	if err != nil {
		fmt.Println("Error opening process handle 2:", err)
		return
	}
	defer windows.CloseHandle(handle2)

	// 比较句柄
	equal := CompareObjectHandles(handle1, handle2)

	if equal {
		fmt.Println("The two processes have the same handle.")
	} else {
		fmt.Println("The two processes do not have the same handle.")
	}
}

func TestOpenFiles(t *testing.T) {
	process, _ := process.NewProcess(21456)
	fds, _ := process.OpenFiles()
	name, _ := process.Name()
	process.Rlimit()
	fmt.Printf("Names: %v\n", name)
	for _, fd := range fds {
		fmt.Printf("Name: %v, fd: %+v\n", name, fd)
	}
}

func TestConnections(t *testing.T) {
	process, _ := process.NewProcess(21456)
	conns, _ := process.Connections()
	name, _ := process.Name()
	fmt.Printf("Names: %v\n", name)
	for _, c := range conns {
		fmt.Printf("  connetions: %+v\n", c)
	}
}

func TestDectect(t *testing.T) {
	var pid1, pid2 int32 = 21456, 18332
	p1, _ := process.NewProcess(pid1)
	p2, _ := process.NewProcess(pid2)
	nodeFilter := map[uintptr]bool{uintptr(pid1): true, uintptr(pid2): true}

	pidHandleMap := TransHandle(GetHandlesForProcess(nodeFilter))

	res := DetectHandle(
		pidHandleMap[uintptr(p1.Pid)],
		pidHandleMap[uintptr(p2.Pid)],
	)

	t.Logf("Res: %v\n", res)
}
