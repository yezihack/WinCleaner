package monitor

import (
	"fmt"
	"strconv"
	"strings"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/winapi"
)

// 已知独显厂商关键词
var discreteKeywords = []string{
	"nvidia", "geforce", "rtx", "gtx", "quadro", "tesla",
	"radeon", "rx ", "rx5", "rx6", "rx7",
	"arc ", "arc a",
}

// 已知核显关键词
var integratedKeywords = []string{
	"intel hd", "intel uhd", "intel iris", "intel(r) hd", "intel(r) uhd", "intel(r) iris",
	"vega", "radeon graphics", "radeon(tm) graphics",
	"microsoft basic", "remote desktop",
}

// GetGPUInfo 通过 WMI 查询显卡信息
func GetGPUInfo() (*model.GPUResult, error) {
	// AdapterRAM 是 uint32，超过 4GB 会溢出
	// 优先从注册表 HardwareInformation.qwMemorySize 获取真实显存（uint64）
	// 回退到 AdapterRAM
	psScript := `
$adapters = Get-CimInstance Win32_VideoController
$regPaths = Get-ChildItem 'HKLM:\SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}' -ErrorAction SilentlyContinue |
  Where-Object { $_.GetValue('DriverDesc') }
foreach ($a in $adapters) {
  $vram = [uint64]$a.AdapterRAM
  foreach ($r in $regPaths) {
    if ($r.GetValue('DriverDesc') -eq $a.Name) {
      $qw = $r.GetValue('HardwareInformation.qwMemorySize')
      if ($qw) { $vram = [uint64]$qw; break }
    }
  }
  "$($a.Name)|$vram|$($a.DriverVersion)|$($a.VideoModeDescription)"
}`
	cmd := winapi.HiddenCmd("powershell", "-NoProfile", "-Command", psScript)
	output, err := cmd.Output()
	if err != nil {
		return &model.GPUResult{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var gpus []model.GPUInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 1 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}

		var vram uint64
		if len(parts) > 1 {
			v, _ := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
			vram = v
		}

		driverVer := ""
		if len(parts) > 2 {
			driverVer = strings.TrimSpace(parts[2])
		}

		resolution := ""
		if len(parts) > 3 {
			res := strings.TrimSpace(parts[3])
			// 格式通常是 "1920 x 1080 x 4294967296 colors"
			if idx := strings.Index(res, " x "); idx > 0 {
				// 提取宽高
				resParts := strings.Split(res, " x ")
				if len(resParts) >= 2 {
					resolution = fmt.Sprintf("%s x %s", resParts[0], resParts[1])
				}
			}
		}

		gpuType, typeLabel := classifyGPU(name)

		gpus = append(gpus, model.GPUInfo{
			Name:       name,
			Type:       gpuType,
			TypeLabel:  typeLabel,
			VRAM:       vram,
			DriverVer:  driverVer,
			Resolution: resolution,
		})
	}

	if len(gpus) == 0 {
		gpus = append(gpus, model.GPUInfo{
			Name:      "未检测到显卡",
			Type:      "none",
			TypeLabel: "无显卡",
		})
	}

	return &model.GPUResult{GPUs: gpus}, nil
}

// classifyGPU 判断显卡类型
func classifyGPU(name string) (string, string) {
	lower := strings.ToLower(name)

	// 先检查核显（因为有些名字同时包含厂商名）
	for _, kw := range integratedKeywords {
		if strings.Contains(lower, kw) {
			return "integrated", "核显"
		}
	}

	// 再检查独显
	for _, kw := range discreteKeywords {
		if strings.Contains(lower, kw) {
			return "discrete", "独显"
		}
	}

	// 有显存大于 512MB 的大概率是独显
	// 这里无法拿到 vram 参数，默认归为核显
	return "integrated", "核显"
}
