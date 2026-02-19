package monitor

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
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

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://www.cip.cc", nil)
	if err != nil {
		return cachedIP, cachedLocation, cachedOperator
	}
	req.Header.Set("User-Agent", "curl/7.0")

	resp, err := client.Do(req)
	if err != nil {
		return cachedIP, cachedLocation, cachedOperator
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cachedIP, cachedLocation, cachedOperator
	}

	// 解析 cip.cc 返回格式
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "IP") {
			cachedIP = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		} else if strings.HasPrefix(line, "地址") {
			cachedLocation = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		} else if strings.HasPrefix(line, "运营商") {
			cachedOperator = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}
	}
	ipCacheTime = time.Now()

	return cachedIP, cachedLocation, cachedOperator
}

// GetMemoryInfo 获取内存信息
func GetMemoryInfo() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}
