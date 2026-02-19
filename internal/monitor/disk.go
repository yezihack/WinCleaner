package monitor

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/winapi"

	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskList 获取所有磁盘分区信息
func GetDiskList() ([]model.DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []model.DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		disks = append(disks, model.DiskInfo{
			Device:      p.Device,
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}
	return disks, nil
}

// ScanLargeFiles 扫描指定分区的大文件（通过 PowerShell，返回前 topN 个）
func ScanLargeFiles(root string, minSizeMB int64, topN int) []model.LargeFileInfo {
	if minSizeMB <= 0 {
		minSizeMB = 50
	}
	if topN <= 0 {
		topN = 100
	}

	minBytes := minSizeMB * 1024 * 1024

	script := `Get-ChildItem -Path '` + root + `' -Recurse -File -ErrorAction SilentlyContinue | ` +
		`Where-Object { $_.Length -ge ` + strconv.FormatInt(minBytes, 10) + ` } | ` +
		`Sort-Object Length -Descending | ` +
		`Select-Object -First ` + strconv.Itoa(topN) + ` | ` +
		`ForEach-Object { "$($_.Length)|$($_.FullName)" }`

	cmd := winapi.HiddenCmd("powershell", "-NoProfile", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var files []model.LargeFileInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) < 2 {
			continue
		}
		size, _ := strconv.ParseInt(parts[0], 10, 64)
		path := parts[1]
		if size == 0 || path == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(path))
		files = append(files, model.LargeFileInfo{
			Path: path,
			Size: size,
			Ext:  ext,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	return files
}
