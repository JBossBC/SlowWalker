package utils

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetLocalIP() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return getLocalIPWindows()
	case "linux":
		return getLocalIPLinux()
	case "darwin":
		return getLocalIPMac()
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

// 适用于 Windows 平台的方法
func getLocalIPWindows() (string, error) {
	ipconfigCmd := exec.Command("cmd", "/c", "ipconfig")
	output, err := ipconfigCmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IPv4") && strings.Contains(line, "192.168.") {
			ip := strings.Fields(line)[1]
			return ip, nil
		}
	}

	return "", fmt.Errorf("local IP not found")
}

// 适用于 Linux 平台的方法
func getLocalIPLinux() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil && strings.HasPrefix(ipNet.IP.String(), "192.168.") {
					return ipNet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("local IP not found")
}

// 适用于 macOS 平台的方法
func getLocalIPMac() (string, error) {
	ifconfigCmd := exec.Command("ifconfig")
	output, err := ifconfigCmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "inet ") && strings.Contains(line, "192.168.") {
			ip := strings.Fields(line)[1]
			return ip, nil
		}
	}

	return "", fmt.Errorf("local IP not found")
}
