package reverseshell

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// SYSTEM_HANDLE_TABLE_ENTRY_INFO  x64
// 此结构用于接收系统中所有handle信息
type SYSTEM_HANDLE_TABLE_ENTRY_INFO struct {
	UniqueProcessId uint16
	Unused1         [4]byte
	HandleValue     uint16 //6  句柄值, 在进程中唯一   uint16取值范围 0~65536
	Unused2         [16]byte
}
type SYSTEM_HANDLE_INFORMATION struct {
	NumberOfHandles uint32
	HandleList      [1]SYSTEM_HANDLE_TABLE_ENTRY_INFO
}

const (
	SystemHandleInformation     = 16
	STATUS_INFO_LENGTH_MISMATCH = 0xc0000004
)

type SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX struct {
	Object                uintptr
	UniqueProcessId       uint32
	HandleValue           uint32
	GrantedAccess         uint32
	CreatorBackTraceIndex uint16
	ObjectTypeIndex       uint16
	HandleAttributes      uint32
	Reserved              uint32
}

type SYSTEM_HANDLE_INFORMATION_EX struct {
	NumberOfHandles uint64
	Reserved        uint64
	Handles         []SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX
}

var currentHandle windows.Handle = windows.InvalidHandle

type PidHandleMap map[uintptr][]windows.Handle

func transHandle(origin PidHandleMap) PidHandleMap {
	newMap := PidHandleMap{}

	for pid, handles := range origin {
		newHandle := []windows.Handle{}
		node, _ := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_DUP_HANDLE, false, uint32(pid))
		for i := 0; i < len(handles); i++ {
			dupHandle := windows.InvalidHandle
			err := windows.DuplicateHandle(node, handles[i], currentHandle,
				&dupHandle, 0, false, windows.DUPLICATE_SAME_ACCESS)
			if err != nil {
				continue
			}
			t, err := windows.GetFileType(dupHandle)
			if err != nil {
				continue
			}
			if t != 2 && t != 3 {
				continue
			}
			newHandle = append(newHandle, dupHandle)
		}
		newMap[pid] = newHandle
	}
	return newMap
}

func getHandlesForProcess(filter map[uintptr]bool) PidHandleMap {
	maps := PidHandleMap{}

	var bufSize uint32 = 0
	var buf *byte
	var err error = windows.STATUS_INFO_LENGTH_MISMATCH

	for i := 0; err == windows.STATUS_INFO_LENGTH_MISMATCH && i < 10; i++ {
		if bufSize > 0 {
			buf = &make([]byte, bufSize)[0]
		}
		err = windows.NtQuerySystemInformation(windows.SystemHandleInformation, unsafe.Pointer(buf), bufSize, &bufSize)
	}
	if err != nil {
		fmt.Printf("NtQuerySystemInformation failed, err:%v, getLastError:%v\n", err, windows.GetLastError())
		return nil
	}

	pSystemHandleInfo := (*SYSTEM_HANDLE_INFORMATION)(unsafe.Pointer(buf))

	var handleInfo = uintptr(unsafe.Pointer(buf)) + 8
	for i := 0; i < int(pSystemHandleInfo.NumberOfHandles); i++ {

		pHandleInfo := (*SYSTEM_HANDLE_TABLE_ENTRY_INFO)(unsafe.Pointer(handleInfo))

		// Compare uniPid
		pid := uintptr(pHandleInfo.UniqueProcessId)
		if filter[pid] {
			maps[pid] = append(maps[pid], windows.Handle(pHandleInfo.HandleValue))
		}
		handleInfo = handleInfo + 24 // sizeof(SYSTEM_HANDLE_TABLE_ENTRY_INFO) = 24
	}
	return maps
}
