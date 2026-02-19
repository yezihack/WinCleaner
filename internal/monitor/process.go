package monitor

import (
	"strings"
	"time"

	"win-cleaner/internal/model"

	"github.com/shirou/gopsutil/v3/process"
)

// GetProcessList 获取进程列表（含 CPU 和内存占用）
func GetProcessList() ([]model.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	// 先采样一次 CPU，等待短暂间隔后再取值才有意义
	for _, p := range procs {
		_, _ = p.CPUPercent()
	}
	time.Sleep(500 * time.Millisecond)

	var result []model.ProcessInfo
	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		cpuPct, _ := p.CPUPercent()
		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}
		memPct, _ := p.MemoryPercent()

		username, _ := p.Username()
		// 简化用户名：去掉域前缀
		if idx := strings.LastIndex(username, "\\"); idx >= 0 {
			username = username[idx+1:]
		}

		statusSlice, _ := p.Status()
		status := ""
		if len(statusSlice) > 0 {
			status = statusSlice[0]
		}

		result = append(result, model.ProcessInfo{
			Pid:        p.Pid,
			Name:       name,
			CPUPercent: cpuPct,
			MemRSS:     memInfo.RSS,
			MemPercent: memPct,
			Status:     status,
			Username:   username,
		})
	}

	return result, nil
}

// KillProcess 结束指定进程
func KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}
