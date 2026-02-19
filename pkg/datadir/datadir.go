package datadir

import (
	"os"
	"path/filepath"
)

const appDir = ".wincleaner"

// Get 返回数据目录路径 (~/.wincleaner)，自动创建
func Get() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return appDir
	}
	dir := filepath.Join(home, appDir)
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// FilePath 返回数据目录下指定文件的完整路径
func FilePath(filename string) string {
	return filepath.Join(Get(), filename)
}
