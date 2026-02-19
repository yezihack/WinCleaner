package model

// SystemInfo 系统概览信息
type SystemInfo struct {
	OS          string  `json:"os"`
	Hostname    string  `json:"hostname"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemTotal    uint64  `json:"mem_total"`
	MemUsed     uint64  `json:"mem_used"`
	MemPercent  float64 `json:"mem_percent"`
	DiskTotal   uint64  `json:"disk_total"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskPercent float64 `json:"disk_percent"`
	PublicIP    string  `json:"public_ip"`
	IPLocation  string  `json:"ip_location"`
	IPOperator  string  `json:"ip_operator"`
}

// ScanResult 扫描结果
type ScanResult struct {
	Category string     `json:"category"`
	Items    []JunkItem `json:"items"`
	Size     int64      `json:"size"`
	Count    int        `json:"count"`
}

// JunkItem 垃圾文件条目
type JunkItem struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Category string `json:"category"`
}

// CleanResult 清理结果
type CleanResult struct {
	FreedSize    int64 `json:"freed_size"`
	CleanedCount int   `json:"cleaned_count"`
	FailedCount  int   `json:"failed_count"`
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	Pid        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemRSS     uint64  `json:"mem_rss"`
	MemPercent float32 `json:"mem_percent"`
	Status     string  `json:"status"`
	Username   string  `json:"username"`
}

// MemoryOptResult 内存优化结果
type MemoryOptResult struct {
	BeforeUsed    uint64  `json:"before_used"`
	AfterUsed     uint64  `json:"after_used"`
	FreedMB       float64 `json:"freed_mb"`
	BeforePercent float64 `json:"before_percent"`
	AfterPercent  float64 `json:"after_percent"`
}

// GPUInfo 显卡信息
type GPUInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`       // "discrete"(独显) / "integrated"(核显) / "none"(无)
	TypeLabel  string `json:"type_label"` // 中文标签
	VRAM       uint64 `json:"vram"`       // 显存（字节）
	DriverVer  string `json:"driver_ver"` // 驱动版本
	Resolution string `json:"resolution"` // 当前分辨率
}

// GPUResult 显卡查询结果
type GPUResult struct {
	GPUs []GPUInfo `json:"gpus"`
}

// CleanRecord 单次清理记录
type CleanRecord struct {
	Date         string `json:"date"`          // 日期 YYYY-MM-DD
	Time         string `json:"time"`          // 时间 HH:MM:SS
	FreedSize    int64  `json:"freed_size"`    // 释放字节数
	CleanedCount int    `json:"cleaned_count"` // 清理文件数
}

// CleanHistory 清理历史
type CleanHistory struct {
	Records []CleanRecord `json:"records"`
}

// CleanHistoryStats 历史统计（返回给前端）
type CleanHistoryStats struct {
	Records       []CleanRecord `json:"records"`         // 全部记录
	DailyStats    []DailyStat   `json:"daily_stats"`     // 按天汇总（近30天）
	MonthlyStats  []MonthlyStat `json:"monthly_stats"`   // 按月汇总
	LastCleanTime string        `json:"last_clean_time"` // 上次清理时间
	LastCleanAgo  string        `json:"last_clean_ago"`  // 距上次清理多久
	TotalFreed    int64         `json:"total_freed"`     // 累计释放
	TotalCount    int           `json:"total_count"`     // 累计清理文件数
}

// DailyStat 按天统计
type DailyStat struct {
	Date      string `json:"date"`
	FreedSize int64  `json:"freed_size"`
	Count     int    `json:"count"`
}

// MonthlyStat 按月统计
type MonthlyStat struct {
	Month     string `json:"month"` // YYYY-MM
	FreedSize int64  `json:"freed_size"`
	Count     int    `json:"count"`
}

// RealtimeStats 实时状态（侧边栏用）
type RealtimeStats struct {
	CPUPercent   float64 `json:"cpu_percent"`
	MemPercent   float64 `json:"mem_percent"`
	NetUpSpeed   uint64  `json:"net_up_speed"`   // 上传速率 bytes/s
	NetDownSpeed uint64  `json:"net_down_speed"` // 下载速率 bytes/s
}

// NetTrafficInfo 网络流量总览
type NetTrafficInfo struct {
	TotalSent uint64 `json:"total_sent"` // 总发送字节
	TotalRecv uint64 `json:"total_recv"` // 总接收字节
	UpSpeed   uint64 `json:"up_speed"`   // 上传速率 bytes/s
	DownSpeed uint64 `json:"down_speed"` // 下载速率 bytes/s
}

// ProcessNetInfo 应用网络流量（按应用名合并）
type ProcessNetInfo struct {
	Name  string `json:"name"`
	Count int    `json:"count"` // 进程数量
	Sent  uint64 `json:"sent"`  // 发送字节
	Recv  uint64 `json:"recv"`  // 接收字节
}

// NetTrafficResult 网络流量结果
type NetTrafficResult struct {
	Overview  NetTrafficInfo   `json:"overview"`
	Processes []ProcessNetInfo `json:"processes"`
}

// NetTrafficRecord 流量采样记录
type NetTrafficRecord struct {
	Timestamp string `json:"timestamp"` // YYYY-MM-DD HH:MM
	Date      string `json:"date"`      // YYYY-MM-DD
	Sent      uint64 `json:"sent"`      // 该采样周期发送字节
	Recv      uint64 `json:"recv"`      // 该采样周期接收字节
}

// NetTrafficHistory 流量历史
type NetTrafficHistory struct {
	Records []NetTrafficRecord `json:"records"`
}

// NetDailyStat 按天流量统计
type NetDailyStat struct {
	Date string `json:"date"`
	Sent uint64 `json:"sent"`
	Recv uint64 `json:"recv"`
}

// NetMonthlyStat 按月流量统计
type NetMonthlyStat struct {
	Month string `json:"month"`
	Sent  uint64 `json:"sent"`
	Recv  uint64 `json:"recv"`
}

// NetYearlyStat 按年流量统计
type NetYearlyStat struct {
	Year string `json:"year"`
	Sent uint64 `json:"sent"`
	Recv uint64 `json:"recv"`
}

// NetTrafficStats 流量历史统计
type NetTrafficStats struct {
	DailyStats   []NetDailyStat   `json:"daily_stats"`
	MonthlyStats []NetMonthlyStat `json:"monthly_stats"`
	YearlyStats  []NetYearlyStat  `json:"yearly_stats"`
	TotalSent    uint64           `json:"total_sent"`
	TotalRecv    uint64           `json:"total_recv"`
}

// DiskInfo 磁盘分区信息
type DiskInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// LargeFileInfo 大文件信息
type LargeFileInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

// DiskScanResult 大文件扫描结果
type DiskScanResult struct {
	Files []LargeFileInfo `json:"files"`
	Count int             `json:"count"`
}

// MemOptRecord 单次内存优化记录
type MemOptRecord struct {
	Date          string  `json:"date"`
	Time          string  `json:"time"`
	FreedMB       float64 `json:"freed_mb"`
	BeforePercent float64 `json:"before_percent"`
	AfterPercent  float64 `json:"after_percent"`
}

// MemOptHistory 内存优化历史
type MemOptHistory struct {
	Records []MemOptRecord `json:"records"`
}

// MemOptDailyStat 按天优化统计
type MemOptDailyStat struct {
	Date    string  `json:"date"`
	FreedMB float64 `json:"freed_mb"`
	Count   int     `json:"count"`
}

// MemOptMonthlyStat 按月优化统计
type MemOptMonthlyStat struct {
	Month   string  `json:"month"`
	FreedMB float64 `json:"freed_mb"`
	Count   int     `json:"count"`
}

// MemOptStats 内存优化历史统计
type MemOptStats struct {
	RecentRecords []MemOptRecord      `json:"recent_records"`
	DailyStats    []MemOptDailyStat   `json:"daily_stats"`
	MonthlyStats  []MemOptMonthlyStat `json:"monthly_stats"`
	LastOptTime   string              `json:"last_opt_time"`
	LastOptAgo    string              `json:"last_opt_ago"`
	TotalFreedMB  float64             `json:"total_freed_mb"`
	TotalCount    int                 `json:"total_count"`
}
