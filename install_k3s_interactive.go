package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 交互式安装函数
func k3sInstallInteractive() error {
	// 选择安装源
	installURL, err := selectInstallSource()
	if err != nil {
		return err
	}

	// 确认安装的频道
	installChannel := confirmOption("Choose install channel (default: stable)", "stable")

	// 确认是否指定安装版本
	installVersion := confirmOption("Specify installation version e.g., v1.28.7+k3s1 (default: latest stable)", defaultVersion)

	// 确认安装数据目录
	installDataDir := confirmOption("Specify installation data directory (default: /var/lib/rancher/k3s)", defaultDataDir)

	// 确认安装二进制文件目录
	installBinDir := confirmOption("Specify installation binary directory (default: /usr/local/bin)", defaultBinDir)

	// 是否禁用traefik
	disableTraefik := confirmOption("Disable Traefik? (yes/no, default: no)", "no")

	// 是否禁用SELinux
	disableSelinux := confirmOption("Disable SELinux? (yes/no, default: no)", "no")

	// 是否使用docker作为容器运行时
	useDocker := confirmOption("Use Docker as container runtime? (yes/no, default: no)", "no")

	// 是否使用其他安装变量
	otherEnvArgs := confirmOption("Specify other installation environment variables \n(separated by comma, e.g. K3S_TOKEN=12345,INSTALL_K3S_SKIP_DOWNLOAD=false, default: none)", "")

	// 是否使用其他安装参数
	otherCmdArgs := confirmOption("Specify other installation command arguments \n(separated by comma, e.g. --flannel-backend=none,--token=12345, default: none)", "")

	// 构建环境变量列表
	envArgs := []string{
		fmt.Sprintf("INSTALL_K3S_CHANNEL=%s", installChannel),
		fmt.Sprintf("INSTALL_K3S_VERSION=%s", installVersion),
		fmt.Sprintf("INSTALL_K3S_DATA_DIR=%s", installDataDir),
		fmt.Sprintf("INSTALL_K3S_BIN_DIR=%s", installBinDir),
	}

	// 处理禁用traefik选项
	if disableTraefik == "yes" {
		otherCmdArgs += ",--disable=traefik"
	}

	// 处理禁用SELinux选项
	if disableSelinux == "yes" {
		envArgs = append(envArgs, "INSTALL_K3S_SKIP_SELINUX_RPM=true", "INSTALL_K3S_SELINUX_WARN=true")
	}

	// 处理使用docker选项
	if useDocker == "yes" {
		otherCmdArgs += ",--docker"
	}

	// 添加其他安装变量和参数
	envArgs = append(envArgs, strings.Split(otherEnvArgs, ",")...)
	cmdArgs := strings.Split(otherCmdArgs, ",")

	// 执行安装
	if err := k3sInstall(installURL, envArgs, cmdArgs); err != nil {
		return err
	}

	return nil
}

// 选择安装源
func selectInstallSource() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	// 提示用户选择安装源
	fmt.Println("Select install source:")
	fmt.Println("1. Official (", officialInstallURL, ")")
	fmt.Println("2. Official China (", officialCNInstallURL, ")")
	fmt.Println("3. Alibaba Cloud (", aliyunInstallURL, ")")

	// 读取用户输入
	fmt.Print("Enter your choice (1/2/3, or press Enter for default - Alibaba Cloud): ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	choice = strings.TrimSpace(choice)

	// 如果用户输入为空，则使用默认的阿里源
	if choice == "" {
		fmt.Println("Default: Alibaba Cloud")
		return aliyunInstallURL, nil
	}

	// 根据用户选择返回对应的安装源
	switch choice {
	case "1":
		return officialInstallURL, nil
	case "2":
		return officialCNInstallURL, nil
	case "3":
		return aliyunInstallURL, nil
	default:
		return "", fmt.Errorf("invalid choice")
	}
}

// 确认选项
func confirmOption(prompt, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}
