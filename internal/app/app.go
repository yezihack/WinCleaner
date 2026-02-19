package app

import (
	"context"
	"time"

	"win-cleaner/internal/cleaner"
	"win-cleaner/internal/memory"
	"win-cleaner/internal/model"
	"win-cleaner/internal/monitor"
)

// App Wails 应用主结构
type App struct {
	ctx         context.Context
	scanResults []model.ScanResult
	stopSampler chan struct{}
}

func NewApp() *App {
	return &App{
		stopSampler: make(chan struct{}),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	// 启动流量采样协程（每 30 秒采样一次）
	go a.netSamplerLoop()
}

func (a *App) Shutdown(ctx context.Context) {
	close(a.stopSampler)
}

func (a *App) netSamplerLoop() {
	// 初始化第一次采样基准
	monitor.RecordNetTrafficSample()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			monitor.RecordNetTrafficSample()
		case <-a.stopSampler:
			return
		}
	}
}

// GetSystemInfo 获取系统信息
func (a *App) GetSystemInfo() (*model.SystemInfo, error) {
	return monitor.GetSystemInfo()
}

// ScanJunk 扫描垃圾文件
func (a *App) ScanJunk() []model.ScanResult {
	categories := cleaner.DefaultCategories()
	a.scanResults = cleaner.Scan(categories)
	return a.scanResults
}

// CleanJunk 清理垃圾文件（传入要清理的分类名列表）
func (a *App) CleanJunk(categoryNames []string) model.CleanResult {
	// 收集选中分类的所有文件
	var items []model.JunkItem
	nameSet := make(map[string]bool)
	for _, name := range categoryNames {
		nameSet[name] = true
	}

	for _, result := range a.scanResults {
		if nameSet[result.Category] {
			items = append(items, result.Items...)
		}
	}

	result := cleaner.Clean(items)

	// 记录清理历史
	_ = cleaner.RecordClean(result)

	return result
}

// GetCleanHistory 获取清理历史统计
func (a *App) GetCleanHistory() (*model.CleanHistoryStats, error) {
	return cleaner.GetCleanHistoryStats()
}

// OptimizeMemory 执行内存优化
func (a *App) OptimizeMemory() (*model.MemoryOptResult, error) {
	result, err := memory.Optimize()
	if err != nil {
		return nil, err
	}
	// 记录优化历史
	_ = memory.RecordOptimize(result)
	return result, nil
}

// GetProcessList 获取进程列表
func (a *App) GetProcessList() ([]model.ProcessInfo, error) {
	return monitor.GetProcessList()
}

// KillProcess 结束进程
func (a *App) KillProcess(pid int32) error {
	return monitor.KillProcess(pid)
}

// GetGPUInfo 获取显卡信息
func (a *App) GetGPUInfo() (*model.GPUResult, error) {
	return monitor.GetGPUInfo()
}

// GetRealtimeStats 获取实时状态（CPU/内存/网速）
func (a *App) GetRealtimeStats() (*model.RealtimeStats, error) {
	return monitor.GetRealtimeStats()
}

// GetNetTraffic 获取网络流量详情
func (a *App) GetNetTraffic() (*model.NetTrafficResult, error) {
	return monitor.GetNetTraffic()
}

// GetNetTrafficStats 获取流量历史统计
func (a *App) GetNetTrafficStats() (*model.NetTrafficStats, error) {
	return monitor.GetNetTrafficStats()
}

// GetDiskList 获取所有磁盘分区信息
func (a *App) GetDiskList() ([]model.DiskInfo, error) {
	return monitor.GetDiskList()
}

// ScanLargeFiles 扫描大文件
func (a *App) ScanLargeFiles(root string, minSizeMB int64) *model.DiskScanResult {
	files := monitor.ScanLargeFiles(root, minSizeMB, 100)
	return &model.DiskScanResult{
		Files: files,
		Count: len(files),
	}
}

// GetMemOptStats 获取内存优化历史统计
func (a *App) GetMemOptStats() (*model.MemOptStats, error) {
	return memory.GetMemOptStats()
}
