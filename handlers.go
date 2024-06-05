package main

import (
	"fmt"
	"github.com/fatih/color"
)

// 安装k3s
func installK3s() error {
	fmt.Println("安装k3s...")
	return autoInstallK3sByLocation([]string{}, []string{})
}

// 部署APP
func deployApp() error {
	fmt.Println("部署APP...")
	return installComponent("app")
}

// 查看服务状态
func viewServiceStatus() error {
	fmt.Println("查看服务状态...")
	return getPodStatus()
}

// 查看APP登录信息
func viewAppLoginInfo() error {
	fmt.Println("查看APP登录信息...")
	return getAppLoginInfo()
}

// 交互式安装k3s
func interactiveInstallK3s() error {
	fmt.Println("交互式安装k3s...")
	return k3sInstallInteractive()
}

// 自定义k3s安装目录
func customizeK3sInstallDir() error {
	fmt.Println("自定义k3s安装目录...")
	return customizeInstallK3sDir()
}

// 自定义k3s安装目录并禁用traefik
func customizeK3sInstallDirAndDisableTraefik() error {
	fmt.Println("自定义k3s安装目录并禁用traefik...")
	return customizeInstallK3sDirAndDisableTraefik()
}

// 自定义k3s安装目录并跳过selinux
func customizeK3sInstallDirAndSkipSelinux() error {
	fmt.Println("自定义k3s安装目录并跳过selinux...")
	return customizeInstallK3sDirAndSkipSelinux()
}

// 安装最新版k3s
func installLatestK3s() error {
	fmt.Println("安装最新版k3s(latest channel install k3s)...")
	return customizeInstallLatestK3s()
}

// 使用官方源安装k3s
func installK3sFromOfficialSource() error {
	fmt.Println("使用官方源安装k3s(curl -sfL https://get.k3s.io | sh -)...")
	return k3sInstall(officialInstallURL, []string{}, []string{})
}

// 程序自动安装源进行k3s安装
func autoInstallK3sFromBestSource() error {
	fmt.Println("程序自动安装源进行k3s安装...")
	return autoInstallK3sByLocation([]string{}, []string{})
}

// 查看k3s子节点安装命令
func viewK3sSubNodeInstallCommand() error {
	fmt.Println("查看k3s子节点安装命令...")
	return getNodeToken()
}

// 测试公共网络连接性
func testPublicNetworkConnectivity() error {
	fmt.Println("正在测试公共网络连接性...")
	if isInternetReachable("http://www.baidu.com") {
		fmt.Println(color.GreenString("OK 公共网络连接性: 可以访问互联网"))
		return nil
	} else {
		fmt.Println(color.YellowString("WARN 公共网络连接性: 无法访问互联网"))
		return fmt.Errorf("无法访问互联网")
	}
}

// 查看公网IP信息
func viewPublicIPInfo() error {
	fmt.Println("查看公网IP信息...")
	return getPublicIPInfo()
}

// 查看CPU信息
func viewCPUInfo() error {
	fmt.Println("查看CPU信息...")
	return getCPUInfo()
}

// 查看内存信息
func viewMemoryInfo() error {
	fmt.Println("查看内存信息...")
	return getMemoryInfo()
}

// 查看磁盘信息
func viewDiskInfo() error {
	fmt.Println("查看磁盘信息...")
	return getDiskInfo()
}

// 查看主机信息
func viewHostInfo() error {
	fmt.Println("查看主机信息...")
	return getHostInfo()
}

// 测试主机带宽
func testHostBandwidth() error {
	fmt.Println("测试主机带宽...")
	return speedTestBandwidth()
}

// 测试磁盘读写
func testDiskReadWrite() error {
	fmt.Println("测试磁盘读写...")
	return diskRWTest()
}

// 安装kubectl自动补全
func installKubectlAutocompletion() error {
	fmt.Println("安装kubectl自动补全...")
	return installKubectlCompletion()
}

// 安装curl
func installCurl() error {
	fmt.Println("安装curl...")
	return checkAndInstallCurl()
}
