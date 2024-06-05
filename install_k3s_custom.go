package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 通用的自定义k3s安装函数
func customizeInstallK3s(dataDirDefault, binDirDefault string, cmdArgs []string, envArgs []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("请输入自定义的k3s数据目录（默认为 %s）：", dataDirDefault)
	dataDirInput, _ := reader.ReadString('\n')
	dataDir := strings.TrimSpace(dataDirInput)
	if dataDir == "" {
		dataDir = dataDirDefault
	}

	fmt.Printf("请输入自定义的k3s二进制文件目录（默认为 %s）：", binDirDefault)
	binDirInput, _ := reader.ReadString('\n')
	binDir := strings.TrimSpace(binDirInput)
	if binDir == "" {
		binDir = binDirDefault
	}

	fmt.Printf("自定义k3s数据目录：%s\n", dataDir)
	fmt.Printf("自定义k3s二进制文件目录：%s\n", binDir)

	// 构造安装命令参数
	cmdArgs = append(cmdArgs, "--data-dir="+dataDir)

	// 设置环境变量
	envArgs = append(envArgs, "INSTALL_K3S_BIN_DIR="+binDir)

	// 调用安装函数
	return autoInstallK3sByLocation(envArgs, cmdArgs)
}

// 自定义安装k3s安装目录函数
func customizeInstallK3sDir() error {
	cmdArgs := []string{}
	envArgs := []string{}
	return customizeInstallK3s(defaultDataDir, defaultBinDir, cmdArgs, envArgs)
}

// 自定义k3s安装目录并禁用traefik函数
func customizeInstallK3sDirAndDisableTraefik() error {
	cmdArgs := []string{"--disable=traefik"}
	envArgs := []string{}
	return customizeInstallK3s(defaultDataDir, defaultBinDir, cmdArgs, envArgs)
}

// 自定义k3s安装目录并跳过selinux函数
func customizeInstallK3sDirAndSkipSelinux() error {
	cmdArgs := []string{}
	envArgs := []string{
		"INSTALL_K3S_SKIP_SELINUX_RPM=true",
		"INSTALL_K3S_SELINUX_WARN=true",
	}
	return customizeInstallK3s(defaultDataDir, defaultBinDir, cmdArgs, envArgs)
}

// 自定义安装最新版k3s函数
func customizeInstallLatestK3s() error {
	cmdArgs := []string{}
	envArgs := []string{"INSTALL_K3S_CHANNEL=latest"}
	return customizeInstallK3s(defaultDataDir, defaultBinDir, cmdArgs, envArgs)
}
