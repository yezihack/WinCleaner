package winapi

import (
	"os/exec"
	"syscall"
)

// HiddenCmd 创建一个隐藏窗口的 exec.Cmd（不弹出 PowerShell 黑窗口）
func HiddenCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	return cmd
}
