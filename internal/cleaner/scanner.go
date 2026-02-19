package cleaner

import (
	"os"
	"path/filepath"
	"win-cleaner/internal/model"
)

// Scan 扫描所有分类的垃圾文件
func Scan(categories []JunkCategory) []model.ScanResult {
	var results []model.ScanResult

	for _, cat := range categories {
		result := model.ScanResult{
			Category: cat.Name,
		}

		for _, dir := range cat.Paths {
			if dir == "" {
				continue
			}
			items := scanDir(dir, cat.Glob, cat.Name)
			result.Items = append(result.Items, items...)
		}

		for _, item := range result.Items {
			result.Size += item.Size
		}
		result.Count = len(result.Items)
		results = append(results, result)
	}

	return results
}

// scanDir 扫描单个目录
func scanDir(dir, glob, category string) []model.JunkItem {
	var items []model.JunkItem

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return items
	}

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无权限的文件
		}
		if info.IsDir() {
			return nil
		}

		// 如果指定了 glob 模式，进行匹配
		if glob != "" {
			matched, _ := filepath.Match(glob, info.Name())
			if !matched {
				return nil
			}
		}

		items = append(items, model.JunkItem{
			Path:     path,
			Size:     info.Size(),
			Category: category,
		})
		return nil
	})

	return items
}

// Clean 清理指定的垃圾文件
func Clean(items []model.JunkItem) model.CleanResult {
	var result model.CleanResult

	for _, item := range items {
		if err := os.Remove(item.Path); err != nil {
			result.FailedCount++
			continue
		}
		result.FreedSize += item.Size
		result.CleanedCount++
	}

	return result
}
