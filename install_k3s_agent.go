package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tokenFileName = "token" // 令牌文件名
)

// 查看k3s子节点安装命令并返回可能的错误
func getNodeToken() error {
	installURL, err := getInstallURL()
	if err != nil {
		return fmt.Errorf("获取安装脚本URL失败：%v", err)
	}

	localIP, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("获取本机IP地址失败：%v", err)
	}

	// 默认数据目录
	dataDir := defaultDataDir

	// 检查是否存在数据目录
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		fmt.Printf("数据目录 %s 不存在\n", dataDir)

		// 提示用户输入自定义的安装目录
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入自定义的k3s数据目录：")
		dataDirInput, _ := reader.ReadString('\n')
		dataDir = strings.TrimSpace(dataDirInput)

		// 尝试从自定义的安装目录读取令牌文件
		tokenFilePath := filepath.Join(dataDir, "server", tokenFileName)
		tokenBytes, err := os.ReadFile(tokenFilePath)
		if err != nil {
			return fmt.Errorf("从自定义的数据目录 %s 读取令牌文件失败：%v", dataDir, err)
		}

		installToken := strings.TrimSpace(string(tokenBytes))

		// 构造安装命令
		installCommand := fmt.Sprintf("curl -sfL %s | K3S_URL=https://%s:6443 K3S_TOKEN=%s sh -", installURL, localIP, installToken)

		fmt.Println("安装k3s子节点的命令:")
		fmt.Println(installCommand)

		return nil
	}

	// 读取令牌文件
	tokenFilePath := filepath.Join(dataDir, "server", tokenFileName)
	tokenBytes, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return fmt.Errorf("读取令牌文件失败：%v", err)
	}

	installToken := strings.TrimSpace(string(tokenBytes))

	// 构造安装命令
	installCommand := fmt.Sprintf("curl -sfL %s | K3S_URL=https://%s:6443 K3S_TOKEN=%s sh -", installURL, localIP, installToken)

	fmt.Println("安装k3s子节点的命令:")
	fmt.Println(installCommand)

	return nil
}
