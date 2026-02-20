package monitor

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"win-cleaner/internal/model"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// IP 信息缓存（避免每次刷新都请求外网）
var (
	cachedIP       string
	cachedLocation string
	cachedOperator string
	ipCacheTime    time.Time
	ipCacheMu      sync.Mutex
)

// GetSystemInfo 获取系统概览
func GetSystemInfo() (*model.SystemInfo, error) {
	hostname, _ := os.Hostname()

	// CPU
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	cpuUsage := 0.0
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// 内存
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 磁盘（C 盘）
	diskStat, err := disk.Usage("C:\\")
	if err != nil {
		return nil, err
	}

	// 公网 IP（缓存 5 分钟）
	ip, location, operator := getPublicIPInfo()

	return &model.SystemInfo{
		OS:          runtime.GOOS + " " + runtime.GOARCH,
		Hostname:    hostname,
		CPUUsage:    cpuUsage,
		MemTotal:    memStat.Total,
		MemUsed:     memStat.Used,
		MemPercent:  memStat.UsedPercent,
		DiskTotal:   diskStat.Total,
		DiskUsed:    diskStat.Used,
		DiskPercent: diskStat.UsedPercent,
		PublicIP:    ip,
		IPLocation:  location,
		IPOperator:  operator,
	}, nil
}

// getPublicIPInfo 获取公网 IP 信息（带缓存）
func getPublicIPInfo() (string, string, string) {
	ipCacheMu.Lock()
	defer ipCacheMu.Unlock()

	if cachedIP != "" && time.Since(ipCacheTime) < 5*time.Minute {
		return cachedIP, cachedLocation, cachedOperator
	}

	// 使用 ip-api.com JSON 接口（免费、稳定、支持中文）
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/?lang=zh-CN&fields=query,country,regionName,city,isp")
	if err != nil {
		return cachedIP, cachedLocation, cachedOperator
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cachedIP, cachedLocation, cachedOperator
	}

	var result struct {
		Query      string `json:"query"`
		Country    string `json:"country"`
		RegionName string `json:"regionName"`
		City       string `json:"city"`
		ISP        string `json:"isp"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return cachedIP, cachedLocation, cachedOperator
	}

	if result.Query != "" {
		cachedIP = result.Query
		cachedLocation = result.Country + " " + result.RegionName + " " + result.City
		cachedOperator = result.ISP
		ipCacheTime = time.Now()
	}

	return cachedIP, cachedLocation, cachedOperator
}

// GetMemoryInfo 获取内存信息
func GetMemoryInfo() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}
