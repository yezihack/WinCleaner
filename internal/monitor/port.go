package monitor

import (
	"strconv"
	"strings"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/winapi"
)

func GetPortListSimple(port uint16) ([]model.PortInfo, error) {
	psScript := `$conns = Get-NetTCPConnection -LocalPort ` + strconv.FormatUint(uint64(port), 10) + ` -ErrorAction SilentlyContinue
foreach ($c in $conns) {
  $proc = Get-Process -Id $c.OwningProcess -ErrorAction SilentlyContinue
  $name = if ($proc) { $proc.ProcessName } else { "Unknown" }
  "$($c.LocalAddress):$($c.LocalPort)|$($c.RemoteAddress):$($c.RemotePort)|$($c.State)|$($c.OwningProcess)|$name"
}`

	output, err := winapi.HiddenCmd("powershell", "-NoProfile", "-Command", psScript).Output()
	if err != nil {
		return nil, err
	}

	var result []model.PortInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}

		localAddr := strings.TrimSpace(parts[0])
		state := strings.TrimSpace(parts[2])
		pid, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
		procName := strings.TrimSpace(parts[4])

		result = append(result, model.PortInfo{
			ListenAddr:  localAddr,
			Port:        port,
			Proto:       "tcp",
			PID:         int32(pid),
			ProcessName: procName,
			Status:      state,
		})
	}

	return result, nil
}

func KillProcessesByPort(port uint16) (int, error) {
	ports, err := GetPortListSimple(port)
	if err != nil {
		return 0, err
	}

	killed := 0
	for _, p := range ports {
		if err := KillProcess(p.PID); err == nil {
			killed++
		}
	}

	return killed, nil
}

func GetPortList(port uint16) ([]model.PortInfo, error) {
	return GetPortListSimple(port)
}

func GetListeningPorts(minPort uint16) ([]model.PortInfo, error) {
	psScript := `$conns = Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue | Where-Object { $_.LocalPort -ge ` + strconv.FormatUint(uint64(minPort), 10) + ` }
foreach ($c in $conns) {
  $proc = Get-Process -Id $c.OwningProcess -ErrorAction SilentlyContinue
  $name = if ($proc) { $proc.ProcessName } else { "Unknown" }
  "$($c.LocalAddress):$($c.LocalPort)|$($c.RemoteAddress):$($c.RemotePort)|$($c.State)|$($c.OwningProcess)|$name"
}`

	output, err := winapi.HiddenCmd("powershell", "-NoProfile", "-Command", psScript).Output()
	if err != nil {
		return nil, err
	}

	var result []model.PortInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}

		addrParts := strings.SplitN(parts[0], ":", 2)
		localPort, _ := strconv.ParseUint(strings.TrimSpace(addrParts[1]), 10, 64)
		state := strings.TrimSpace(parts[2])
		pid, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
		procName := strings.TrimSpace(parts[4])

		result = append(result, model.PortInfo{
			ListenAddr:  strings.TrimSpace(parts[0]),
			Port:        uint16(localPort),
			Proto:       "tcp",
			PID:         int32(pid),
			ProcessName: procName,
			Status:      state,
		})
	}

	return result, nil
}