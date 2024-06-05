package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"os/exec"
	"runtime"
	"strings"
)

func runSystemChecks() {
	checkFunctions := []struct {
		name     string
		checkFn  func() error
		okFormat string
	}{
		{"交换空间检查", checkSwap, "交换空间已关闭"},
		{"防火墙状态检查", checkFirewall, "防火墙已关闭"},
		{"互联网连接检查", checkInternet, "成功连接到互联网"},
		{"CPU核数检查", checkCPU, "CPU核数达到要求"},
		{"系统内存检查", checkMemory, "系统内存满足要求"},
		{"系统最大磁盘空间检查", checkDiskSpace, "系统最大磁盘空间满足要求"},
	}

	for _, check := range checkFunctions {
		fmt.Printf("正在执行 %s...\n", check.name)
		err := check.checkFn()
		if err != nil {
			fmt.Printf(color.YellowString("WARN %s: %v\n"), check.name, err)
		} else {
			fmt.Printf(color.GreenString("OK %s: %s\n"), check.name, check.okFormat)
		}
	}
}

func checkSwap() error {
	//fmt.Println("检查交换空间...")

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return fmt.Errorf("获取交换空间信息失败: %v", err)
	}

	if swapInfo.Total > 0 {
		return errors.New("发现交换空间，请考虑关闭交换空间以获得更好的性能。")
	}

	return nil
}

func checkFirewall() error {
	//fmt.Println("检查防火墙状态...")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		firewallCmd, err := getLinuxFirewallCmd()
		if err != nil {
			return err
		}
		cmd = exec.Command("sh", "-c", firewallCmd)
	case "windows":
		cmd = exec.Command("cmd", "/c", "NetSh Advfirewall show allprofiles")
	case "darwin":
		return nil // macOS默认没有启用防火墙，直接返回OK
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("无法检查防火墙状态: %v", err)
	}

	firewallStatus := strings.TrimSpace(string(output))
	if strings.Contains(firewallStatus, "inactive") {
		return nil // 防火墙已关闭
	} else if strings.Contains(firewallStatus, "active") {
		return errors.New("防火墙未关闭")
	} else {
		return fmt.Errorf("无法识别防火墙状态")
	}
}

// 根据 Linux 发行版识别防火墙管理工具命令
func getLinuxFirewallCmd() (string, error) {
	distroCmd := exec.Command("cat", "/etc/os-release")
	output, err := distroCmd.Output()
	if err != nil {
		return "", fmt.Errorf("无法获取Linux发行版信息: %v", err)
	}
	distroInfo := string(output)

	switch {
	case strings.Contains(distroInfo, "Ubuntu"), strings.Contains(distroInfo, "Debian"):
		return "sudo ufw status", nil
	case strings.Contains(distroInfo, "CentOS"), strings.Contains(distroInfo, "Red Hat"), strings.Contains(distroInfo, "Fedora"):
		return "sudo systemctl status firewalld", nil
	case strings.Contains(distroInfo, "Arch Linux"):
		return "sudo systemctl status iptables", nil
	default:
		return "", fmt.Errorf("未知的Linux发行版, 无法检查防火墙状态")
	}
}

func checkInternet() error {
	//fmt.Println("检查互联网连接...")

	if isInternetReachable() {
		return nil // 成功连接到互联网
	} else {
		return fmt.Errorf("无法连接到互联网")
	}
}

func checkCPU() error {
	//fmt.Println("CPU核数检查...")

	cpuInfo, err := cpu.Counts(true) // 获取逻辑 CPU 核心数
	if err != nil {
		return fmt.Errorf("无法获取CPU信息: %v", err)
	}

	requiredCores := 8
	if cpuInfo >= requiredCores {
		return nil // CPU核数达到要求
	} else {
		return fmt.Errorf("当前系统CPU核心数为 %d，需要至少 %d 颗", cpuInfo, requiredCores)
	}
}

func checkMemory() error {
	//fmt.Println("系统内存大小检查...")

	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("无法获取内存信息: %v", err)
	}

	// 将内存大小从字节转换为GB
	GB := uint64(1024 * 1024 * 1024)
	totalMemoryGB := memoryInfo.Total / GB

	requiredMemoryGB := uint64(32) // 需要的内存大小，单位为GB

	if totalMemoryGB >= requiredMemoryGB {
		return nil // 系统内存满足要求
	} else {
		return fmt.Errorf("当前系统内存总量为 %.2f GB，需要至少 %.2f GB", float64(totalMemoryGB), float64(requiredMemoryGB))
	}
}

func checkDiskSpace() error {
	//fmt.Println("系统最大磁盘空间检查...")

	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("无法获取磁盘分区信息: %v", err)
	}

	var maxPartitionSize float64 = 0
	for _, partition := range partitions {
		if partition.Mountpoint != "" {
			usageStat, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				continue // 忽略获取使用情况失败的分区
			}

			// 更新最大分区的大小
			partitionSizeGB := float64(usageStat.Total) / float64(1024*1024*1024)
			if partitionSizeGB > maxPartitionSize {
				maxPartitionSize = partitionSizeGB
			}
		}
	}

	requiredSpaceGB := float64(500) // 需要的最大磁盘空间大小，单位为GB

	if maxPartitionSize >= requiredSpaceGB {
		return nil // 最大分区的磁盘空间满足要求
	} else {
		return fmt.Errorf("当前系统最大分区磁盘空间总量为 %.2f GB，需要至少 %.2f GB", maxPartitionSize, requiredSpaceGB)
	}
}

// DisableSwap disables swap space on the system.
func disableSwap() error {
	_, err := exec.Command("swapoff", "-a").Output()
	if err != nil {
		return fmt.Errorf("failed to disable swap space: %v", err)
	}
	return nil
}

// DisableLinuxFirewall disables the firewall on the system.
func disableLinuxFirewall() error {
	firewallCmd, err := getLinuxFirewallCommand()
	if err != nil {
		return err
	}

	switch firewallCmd {
	case "ufw":
		_, err := exec.Command("ufw", "disable").Output()
		if err != nil {
			return fmt.Errorf("failed to disable ufw firewall: %v", err)
		}
	case "firewalld":
		_, err := exec.Command("systemctl", "stop", "firewalld").Output()
		if err != nil {
			return fmt.Errorf("failed to stop firewalld: %v", err)
		}
		_, err = exec.Command("systemctl", "disable", "firewalld").Output()
		if err != nil {
			return fmt.Errorf("failed to disable firewalld: %v", err)
		}
	case "SuSEfirewall2":
		_, err := exec.Command("rcSuSEfirewall2", "stop").Output()
		if err != nil {
			return fmt.Errorf("failed to stop SuSEfirewall2: %v", err)
		}
	default:
		return fmt.Errorf("unsupported firewall management tool: %s", firewallCmd)
	}

	return nil
}

// getLinuxFirewallCommand determines the firewall management tool based on the Linux distribution.
func getLinuxFirewallCommand() (string, error) {
	distroCmd := exec.Command("lsb_release", "-si")
	output, err := distroCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to determine Linux distribution: %v", err)
	}

	distro := strings.TrimSpace(string(output))

	switch distro {
	case "Ubuntu", "Debian":
		return "ufw", nil
	case "CentOS", "Fedora", "RedHatEnterpriseServer":
		return "firewalld", nil
	case "SUSE":
		return "SuSEfirewall2", nil
	default:
		return "", fmt.Errorf("unsupported Linux distribution: %s", distro)
	}
}
