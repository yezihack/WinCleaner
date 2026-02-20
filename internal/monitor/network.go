package monitor

import (
	"encoding/base64"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf16"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/winapi"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var (
	lastNetSent uint64
	lastNetRecv uint64
	lastNetTime time.Time
	netMu       sync.Mutex
)

// GetRealtimeStats 获取实时状态（CPU、内存、网速）
func GetRealtimeStats() (*model.RealtimeStats, error) {
	cpuPct, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	cpuUsage := 0.0
	if len(cpuPct) > 0 {
		cpuUsage = cpuPct[0]
	}

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	upSpeed, downSpeed := getNetSpeed()

	return &model.RealtimeStats{
		CPUPercent:   cpuUsage,
		MemPercent:   memStat.UsedPercent,
		NetUpSpeed:   upSpeed,
		NetDownSpeed: downSpeed,
	}, nil
}

// getNetSpeed 计算网络速率
func getNetSpeed() (uint64, uint64) {
	netMu.Lock()
	defer netMu.Unlock()

	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return 0, 0
	}

	currentSent := counters[0].BytesSent
	currentRecv := counters[0].BytesRecv
	now := time.Now()

	var upSpeed, downSpeed uint64
	if !lastNetTime.IsZero() {
		elapsed := now.Sub(lastNetTime).Seconds()
		if elapsed > 0 {
			upSpeed = uint64(float64(currentSent-lastNetSent) / elapsed)
			downSpeed = uint64(float64(currentRecv-lastNetRecv) / elapsed)
		}
	}

	lastNetSent = currentSent
	lastNetRecv = currentRecv
	lastNetTime = now

	return upSpeed, downSpeed
}

// GetNetTraffic 获取网络流量详情（总览 + 按应用）
func GetNetTraffic() (*model.NetTrafficResult, error) {
	counters, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	overview := model.NetTrafficInfo{}
	if len(counters) > 0 {
		overview.TotalSent = counters[0].BytesSent
		overview.TotalRecv = counters[0].BytesRecv
	}
	overview.UpSpeed, overview.DownSpeed = getNetSpeed()

	processes := getProcessNetUsage()

	return &model.NetTrafficResult{
		Overview:  overview,
		Processes: processes,
	}, nil
}

// toUTF16LEBase64 将 PowerShell 脚本编码为 -EncodedCommand 所需的 Base64 格式
func toUTF16LEBase64(s string) string {
	runes := utf16.Encode([]rune(s))
	bytes := make([]byte, len(runes)*2)
	for i, r := range runes {
		bytes[i*2] = byte(r)
		bytes[i*2+1] = byte(r >> 8)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// getProcessNetUsage 获取应用级网络流量（按应用名合并多个子进程）
func getProcessNetUsage() []model.ProcessNetInfo {
	// 批量查询：先获取有网络连接的 PID，再一次性 WMI 查询所有进程 IO 计数器
	psScript := `$sep = [char]9
$pids = @(Get-NetTCPConnection -State Established,Listen -ErrorAction SilentlyContinue |
  Select-Object -ExpandProperty OwningProcess -Unique |
  Where-Object { $_ -ne 0 })
if ($pids.Count -eq 0) { exit }
$filter = ($pids | ForEach-Object { "ProcessId=$_" }) -join ' OR '
$procs = Get-CimInstance Win32_Process -Filter $filter -ErrorAction SilentlyContinue |
  Where-Object { $_.Name -notin @('powershell.exe','conhost.exe','System') }
$grouped = $procs | Group-Object Name
foreach ($g in $grouped) {
  $name = $g.Name -replace '\.exe$',''
  $cnt = $g.Count
  $w = [uint64]0
  $r = [uint64]0
  foreach ($p in $g.Group) {
    if ($p.WriteTransferCount) { $w += [uint64]$p.WriteTransferCount }
    if ($p.ReadTransferCount) { $r += [uint64]$p.ReadTransferCount }
  }
  "$name$sep$cnt$sep$w$sep$r"
}`

	encoded := toUTF16LEBase64(psScript)
	cmd := winapi.HiddenCmd("powershell", "-NoProfile", "-EncodedCommand", encoded)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var result []model.ProcessNetInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 4 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		count, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		sent, _ := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
		recv, _ := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 64)

		if name == "" || (sent == 0 && recv == 0) {
			continue
		}

		if count == 0 {
			count = 1
		}

		result = append(result, model.ProcessNetInfo{
			Name:  name,
			Count: count,
			Sent:  sent,
			Recv:  recv,
		})
	}
	return result
}
