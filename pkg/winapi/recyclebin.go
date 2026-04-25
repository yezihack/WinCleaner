package winapi

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type RecycleBinInfo struct {
	ItemCount int64
	SizeBytes int64
}

func GetRecycleBinInfo() (RecycleBinInfo, error) {
	script := `
		$shell = New-Object -ComObject Shell.Application
		$recycleBin = $shell.NameSpace(10)
		$items = $recycleBin.Items()
		$count = $items.Count
		$size = 0
		foreach ($item in $items) {
			$size += $item.Size
		}
		Write-Output "$count|$size"
	`
	cmd := HiddenCmd("powershell", "-NoProfile", "-Command", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return RecycleBinInfo{}, fmt.Errorf("获取回收站信息失败: %w, stderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return RecycleBinInfo{ItemCount: 0, SizeBytes: 0}, nil
	}

	parts := strings.Split(output, "|")
	if len(parts) != 2 {
		return RecycleBinInfo{}, fmt.Errorf("回收站信息格式错误: %s", output)
	}

	count, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return RecycleBinInfo{}, fmt.Errorf("解析回收站数量失败: %w", err)
	}

	size, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return RecycleBinInfo{}, fmt.Errorf("解析回收站大小失败: %w", err)
	}

	return RecycleBinInfo{
		ItemCount: count,
		SizeBytes: size,
	}, nil
}

func EmptyRecycleBin() error {
	script := `Clear-RecycleBin -Force -ErrorAction SilentlyContinue`
	cmd := HiddenCmd("powershell", "-NoProfile", "-Command", script)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("清空回收站失败: %w, stderr: %s", err, stderr.String())
	}
	return nil
}
