package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// 查看CPU信息
func getCPUInfo() error {
	// 获取所有逻辑CPU信息
	cpuInfos, err := cpu.Info()
	if err != nil {
		return fmt.Errorf("无法获取CPU信息: %v", err)
	}

	// 遍历每个逻辑CPU信息并打印
	for idx, cpuInfo := range cpuInfos {
		fmt.Printf("CPU %d:\n", idx+1)
		fmt.Printf("  CPU编号: %d\n", cpuInfo.CPU)
		fmt.Printf("  厂商ID: %s\n", cpuInfo.VendorID)
		fmt.Printf("  CPU系列: %s\n", cpuInfo.Family)
		fmt.Printf("  型号: %s\n", cpuInfo.Model)
		fmt.Printf("  步进: %d\n", cpuInfo.Stepping)
		fmt.Printf("  物理ID: %s\n", cpuInfo.PhysicalID)
		fmt.Printf("  核心ID: %s\n", cpuInfo.CoreID)
		fmt.Printf("  核心数: %d\n", cpuInfo.Cores)
		fmt.Printf("  型号名称: %s\n", cpuInfo.ModelName)
		fmt.Printf("  频率: %.2f MHz\n", cpuInfo.Mhz)
		fmt.Printf("  缓存大小: %d KB\n", cpuInfo.CacheSize)
		fmt.Printf("  特性标志: %v\n", cpuInfo.Flags)
		fmt.Printf("  微码版本: %s\n", cpuInfo.Microcode)
		fmt.Println("-------------------------")
	}

	return nil
}

// 查看内存信息
func getMemoryInfo() error {
	// 获取内存信息
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("无法获取内存信息: %v", err)
	}

	// 输出内存信息
	fmt.Println("系统内存信息:")
	fmt.Printf("  总内存: %.2f GB\n", toGiB(memoryInfo.Total))
	fmt.Printf("  可用内存: %.2f GB\n", toGiB(memoryInfo.Available))
	fmt.Printf("  已使用内存: %.2f GB\n", toGiB(memoryInfo.Used))
	fmt.Printf("  使用率: %.2f%%\n", memoryInfo.UsedPercent)

	return nil
}

// 查看磁盘信息
func getDiskInfo() error {
	// 获取磁盘分区信息
	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("无法获取磁盘分区信息: %v", err)
	}

	// 输出磁盘信息
	fmt.Println("系统磁盘信息:")
	for _, partition := range partitions {
		if partition.Mountpoint != "" {
			fmt.Printf("  分区挂载点: %s\n", partition.Mountpoint)
			fmt.Printf("  文件系统类型: %s\n", partition.Fstype)

			// 获取分区使用情况
			usageStat, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				fmt.Printf("    无法获取分区使用情况: %v\n", err)
			} else {
				fmt.Printf("  总容量: %.2f GB\n", toGiB(usageStat.Total))
				fmt.Printf("  可用空间: %.2f GB\n", toGiB(usageStat.Free))
				fmt.Printf("  使用率: %.2f%%\n", usageStat.UsedPercent)
			}

			fmt.Println() // 输出空行分隔各个分区信息
		}
	}

	return nil
}

// 查看主机信息
func getHostInfo() error {
	info, err := host.Info()
	if err != nil {
		return fmt.Errorf("无法获取主机信息: %v", err)
	}

	fmt.Println("主机信息:")
	fmt.Printf("  Hostname:             %s\n", info.Hostname)
	fmt.Printf("  Uptime:               %d seconds\n", info.Uptime)
	fmt.Printf("  BootTime:             %d (unixtime)\n", info.BootTime)
	fmt.Printf("  Number of processes:  %d\n", info.Procs)
	fmt.Printf("  OS:                   %s\n", info.OS)
	fmt.Printf("  Platform:             %s\n", info.Platform)
	fmt.Printf("  PlatformFamily:       %s\n", info.PlatformFamily)
	fmt.Printf("  PlatformVersion:      %s\n", info.PlatformVersion)
	fmt.Printf("  KernelVersion:        %s\n", info.KernelVersion)
	fmt.Printf("  KernelArch:           %s\n", info.KernelArch)
	fmt.Printf("  VirtualizationSystem: %s\n", info.VirtualizationSystem)
	fmt.Printf("  VirtualizationRole:   %s\n", info.VirtualizationRole)
	fmt.Printf("  HostID:               %s\n", info.HostID)

	return nil
}

// 将字节数转换为 GiB 单位
func toGiB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024 / 1024
}
