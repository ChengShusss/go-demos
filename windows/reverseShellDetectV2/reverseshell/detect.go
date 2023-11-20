package reverseshell

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

var (
	modKernelBase            = windows.NewLazySystemDLL("Kernelbase.dll")
	procCompareObjectHandles = modKernelBase.NewProc("CompareObjectHandles")

	shellList = map[string]bool{
		"cmd.exe":        true,
		"powershell.exe": true,
	}
)

type ProcessExt struct {
	Process     *process.Process
	Name        string
	Path        string
	Connections []net.ConnectionStat
}

func DetectProcesses(processes []*process.Process) bool {
	for _, process := range processes {
		DetectSingleProcess(process)
	}
	return false
}

func DetectByPid(pid int32) bool {
	process, err := process.NewProcess(pid)
	if err != nil {
		return false
	}
	return DetectSingleProcess(process)
}

func DetectSingleProcess(p *process.Process) bool {
	exe, err := p.Exe()
	if err != nil {
		// fmt.Printf("Failed to get process %v, err: %v\n", process.Pid, err)
		return false
	}
	if !isShell(exe) {
		// fmt.Printf("Process %v is not shell, continue\n", process.Pid)
		return false
	}

	pNew := ProcessExt{
		Process: p,
		Path:    exe,
		Name:    filepath.Base(exe),
	}

	return checkProcessTree(&pNew)
}

func checkProcessTree(current *ProcessExt) bool {

	var nodeList []*process.Process
	nodeFilter := map[uintptr]bool{uintptr(current.Process.Pid): true}

	parent := current.Process
	var err error
	for {
		parent, err = parent.Parent()
		if err != nil {
			break
		}
		conns, err := parent.Connections()
		if err != nil {
			continue
		}
		if len(conns) > 0 {

			if inWhite(parent, conns) {
				// Should have debug
				continue
			}
			nodeList = append(nodeList, parent)
			nodeFilter[uintptr(parent.Pid)] = true
		}
	}

	fmt.Printf("Fileter: %+v\n", nodeFilter)

	pidHandleMap := transHandle(getHandlesForProcess(nodeFilter))

	for _, node := range nodeList {
		cond1 := hasSharedHandle(
			pidHandleMap[uintptr(node.Pid)],
			pidHandleMap[uintptr(current.Process.Pid)],
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
func hasSharedHandle(h1, h2 []windows.Handle) bool {

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

			if compareObjectHandles(handle1, handle2) {
				return true
			}
		}
	}

	return false
}

func compareObjectHandles(handle1, handle2 windows.Handle) bool {
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
	return shellList[strings.ToLower(exe)]
}

func inWhite(p *process.Process, conns []net.ConnectionStat) bool {

	for _, conn := range conns {
		if conn.Raddr.IP == "10.231.2.248" {
			return true
		}
	}

	return false
}
