package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func installKubectlCompletion() error {
	// 确定kubectl是否已安装
	if !isKubectlInstalled() {
		return fmt.Errorf("kubectl未安装，请先安装kubectl")
	}

	// 确定当前shell类型
	shell := detectShell()
	if shell == "" {
		return fmt.Errorf("无法检测到当前shell类型")
	}

	// 执行kubectl自动补全安装命令
	switch shell {
	case "bash":
		cmd := exec.Command("kubectl", "completion", "bash")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("安装kubectl自动补全失败: %v", err)
		}
		if !isCompletionInBashrc() {
			addCompletionToBashrc()
			fmt.Println("已将kubectl自动补全添加到 ~/.bashrc 文件中。")
		} else {
			fmt.Println("kubectl自动补全已存在于 ~/.bashrc 文件中，无需重复添加。")
		}
	default:
		return fmt.Errorf("不支持的shell类型：%s", shell)
	}

	fmt.Println("kubectl自动补全安装成功。")
	return nil
}

// 检查kubectl是否已安装
func isKubectlInstalled() bool {
	_, err := exec.LookPath("kubectl")
	return err == nil
}

// 检测当前shell类型
func detectShell() string {
	shell := os.Getenv("SHELL")
	switch {
	case strings.Contains(shell, "bash"):
		return "bash"
	case strings.Contains(shell, "zsh"):
		return "zsh"
	default:
		return ""
	}
}

// 检查是否已将kubectl自动补全添加到~/.bashrc
func isCompletionInBashrc() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	bashrcPath := homeDir + "/.bashrc"
	content, err := os.ReadFile(bashrcPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "kubectl completion bash")
}

// 将kubectl自动补全命令添加到~/.bashrc
func addCompletionToBashrc() error {
	cmd := exec.Command("echo", "'source <(kubectl completion bash)' >> ~/.bashrc")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("将kubectl自动补全添加到 ~/.bashrc 文件失败: %v", err)
	}
	return nil
}
