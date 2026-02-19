package winapi

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var (
	modKernel32              = windows.NewLazySystemDLL("kernel32.dll")
	modPsapi                 = windows.NewLazySystemDLL("psapi.dll")
	procSetProcessWorkingSet = modKernel32.NewProc("SetProcessWorkingSetSize")
	procEmptyWorkingSet      = modPsapi.NewProc("EmptyWorkingSet")
)

// EmptyWorkingSet 清空指定进程的工作集（释放物理内存到页面文件）
func EmptyWorkingSet(handle windows.Handle) error {
	ret, _, err := procEmptyWorkingSet.Call(uintptr(handle))
	if ret == 0 {
		return fmt.Errorf("EmptyWorkingSet failed: %w", err)
	}
	return nil
}

// TrimProcessMemory 收缩进程工作集
func TrimProcessMemory(handle windows.Handle) error {
	// -1 表示让系统自动决定最小/最大工作集
	ret, _, err := procSetProcessWorkingSet.Call(
		uintptr(handle),
		uintptr(0xFFFFFFFFFFFFFFFF),
		uintptr(0xFFFFFFFFFFFFFFFF),
	)
	if ret == 0 {
		return fmt.Errorf("SetProcessWorkingSetSize failed: %w", err)
	}
	return nil
}
