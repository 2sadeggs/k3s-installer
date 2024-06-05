package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
)

func checkAndInstallCurl() error {
	// 检查 curl 是否已安装
	if _, err := exec.LookPath("curl"); err == nil {
		fmt.Println("Curl is already installed.")
		return nil // curl 已安装,无需操作
	}

	// 启动加载动画
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Installing curl... "
	s.Start()
	defer s.Stop()

	// 根据操作系统安装 curl
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		fmt.Println("Installing curl on Linux...")
		cmd = exec.Command("sh", "-c", "sudo apt install -y curl || sudo yum install -y curl || sudo pacman -S --noconfirm curl")
	case "windows":
		fmt.Println("Installing curl on Windows...")
		cmd = exec.Command("choco", "install", "-y", "curl")
	case "darwin":
		fmt.Println("Installing curl on MacOS...")
		cmd = exec.Command("brew", "install", "curl")
	default:
		return fmt.Errorf("unsupported operating system for curl installation")
	}

	// 绑定标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令并直接返回结果
	return cmd.Run()
}
