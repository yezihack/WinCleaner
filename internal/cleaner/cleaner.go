package cleaner

import (
	"os"
	"path/filepath"
)

// 垃圾分类定义
type JunkCategory struct {
	Name  string
	Paths []string // 支持环境变量
	Glob  string   // 文件匹配模式，空则匹配所有
}

// DefaultCategories 默认扫描分类
func DefaultCategories() []JunkCategory {
	temp := os.Getenv("TEMP")
	localAppData := os.Getenv("LOCALAPPDATA")
	winDir := os.Getenv("WINDIR")

	return []JunkCategory{
		{
			Name:  "系统临时文件",
			Paths: []string{temp, filepath.Join(winDir, "Temp")},
		},
		{
			Name:  "Windows Update 缓存",
			Paths: []string{filepath.Join(winDir, "SoftwareDistribution", "Download")},
		},
		{
			Name:  "缩略图缓存",
			Paths: []string{filepath.Join(localAppData, "Microsoft", "Windows", "Explorer")},
			Glob:  "thumbcache_*.db",
		},
		{
			Name:  "系统日志",
			Paths: []string{filepath.Join(winDir, "Logs")},
			Glob:  "*.log",
		},
		{
			Name: "浏览器缓存",
			Paths: []string{
				filepath.Join(localAppData, "Google", "Chrome", "User Data", "Default", "Cache"),
				filepath.Join(localAppData, "Microsoft", "Edge", "User Data", "Default", "Cache"),
			},
		},
		{
			Name:  "回收站",
			Paths: []string{"C:\\$Recycle.Bin"},
		},
		{
			Name:  "Windows 预读取",
			Paths: []string{filepath.Join(winDir, "Prefetch")},
			Glob:  "*.pf",
		},
	}
}
