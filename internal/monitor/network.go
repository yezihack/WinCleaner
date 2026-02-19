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
	// 使用 TAB 作为分隔符，PowerShell 中 [char]9 表示 TAB
	psScript := "$sep = [char]9\n" +
		"$pids = Get-NetTCPConnection -State Established,Listen -ErrorAction SilentlyContinue |\n" +
		"  Select-Object -ExpandProperty OwningProcess -Unique\n" +
		"foreach ($p_id in $pids) {\n" +
		"  try {\n" +
		"    $p = Get-CimInstance Win32_Process -Filter \"ProcessId=$p_id\" -ErrorAction SilentlyContinue\n" +
		"    if ($p -and $p.Name -ne 'powershell.exe' -and $p.Name -ne 'conhost.exe') {\n" +
		"      $n = $p.Name -replace '\\.exe$',''\n" +
		"      $w = if ($p.WriteTransferCount) { $p.WriteTransferCount } else { 0 }\n" +
		"      $r = if ($p.ReadTransferCount) { $p.ReadTransferCount } else { 0 }\n" +
		"      \"$n$sep$w$sep$r\"\n" +
		"    }\n" +
		"  } catch {}\n" +
		"}"

	encoded := toUTF16LEBase64(psScript)
	cmd := winapi.HiddenCmd("powershell", "-NoProfile", "-EncodedCommand", encoded)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	appMap := make(map[string]*model.ProcessNetInfo)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		sent, _ := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		recv, _ := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)

		if name == "" {
			continue
		}

		if existing, ok := appMap[name]; ok {
			existing.Sent += sent
			existing.Recv += recv
			existing.Count++
		} else {
			appMap[name] = &model.ProcessNetInfo{
				Name:  name,
				Count: 1,
				Sent:  sent,
				Recv:  recv,
			}
		}
	}

	var result []model.ProcessNetInfo
	for _, p := range appMap {
		result = append(result, *p)
	}
	return result
}
