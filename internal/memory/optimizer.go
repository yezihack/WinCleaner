package memory

import (
	"win-cleaner/internal/model"
	"win-cleaner/pkg/winapi"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows"
)

// Optimize 执行内存优化：遍历所有进程，收缩工作集
func Optimize() (*model.MemoryOptResult, error) {
	// 优化前内存状态
	beforeMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 获取所有进程
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	for _, p := range procs {
		handle, err := windows.OpenProcess(
			windows.PROCESS_SET_QUOTA|windows.PROCESS_QUERY_INFORMATION,
			false,
			uint32(p.Pid),
		)
		if err != nil {
			continue // 无权限的进程跳过
		}

		_ = winapi.EmptyWorkingSet(handle)
		_ = winapi.TrimProcessMemory(handle)
		windows.CloseHandle(handle)
	}

	// 优化后内存状态
	afterMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	freedBytes := int64(beforeMem.Used) - int64(afterMem.Used)
	freedMB := float64(freedBytes) / 1024 / 1024
	if freedMB < 0 {
		freedMB = 0
	}

	return &model.MemoryOptResult{
		BeforeUsed:    beforeMem.Used,
		AfterUsed:     afterMem.Used,
		FreedMB:       freedMB,
		BeforePercent: beforeMem.UsedPercent,
		AfterPercent:  afterMem.UsedPercent,
	}, nil
}
