package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type menuOption struct {
	title   string
	handler func() error
}

func mainMenu() {
	// 主菜单
	var mainMenuOptions = []menuOption{
		{"安装k3s", installK3s},
		{"部署APP", deployApp},
		{"查看服务状态", viewServiceStatus},
		{"查看APP登录信息", viewAppLoginInfo},
		{"自定义安装k3s选项", customInstallK3sMenu},
		{"工具箱", toolboxMenu},
		{"交互式安装k3s", interactiveInstallK3s},
		{"退出程序", exit},
	}
	for {
		displayMenu("主菜单", mainMenuOptions)
		choice := getUserChoice()
		if choice >= 0 && choice < len(mainMenuOptions) {
			if err := mainMenuOptions[choice].handler(); err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		} else {
			fmt.Println("无效选项，请重试。")
		}
	}
}

func customInstallK3sMenu() error {
	// 自定义安装k3s子菜单
	installK3sOptions := []menuOption{
		{"自定义k3s安装目录", customizeK3sInstallDir},
		{"自定义k3s安装目录并禁用traefik", customizeK3sInstallDirAndDisableTraefik},
		{"自定义k3s安装目录并跳过selinux(适配kylin\\openEuler\\Delix等国产操作系统)", customizeK3sInstallDirAndSkipSelinux},
		{"安装最新版k3s(latest channel install k3s)", installLatestK3s},
		{"使用官方源安装k3s(curl -sfL https://get.k3s.io | sh -)", installK3sFromOfficialSource},
		{"程序自动选择安装源进行k3s安装", autoInstallK3sFromBestSource},
		{"查看k3s子节点安装命令", viewK3sSubNodeInstallCommand},
		{"返回上一级菜单", nil},
	}
	for {
		displayMenu("自定义安装k3s", installK3sOptions)
		choice := getUserChoice()
		if choice >= 0 && choice < len(installK3sOptions) {
			if installK3sOptions[choice].handler == nil {
				return nil
			}
			if err := installK3sOptions[choice].handler(); err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		} else {
			fmt.Println("无效选项，请重试。")
		}
	}
}

func toolboxMenu() error {
	// 工具箱子菜单
	toolboxOptions := []menuOption{
		{"测试公网连通性", testPublicNetworkConnectivity},
		{"查看公网IP信息", viewPublicIPInfo},
		{"查看CPU信息", viewCPUInfo},
		{"查看内存信息", viewMemoryInfo},
		{"查看磁盘信息", viewDiskInfo},
		{"查看主机信息", viewHostInfo},
		{"测试主机带宽", testHostBandwidth},
		{"测试磁盘读写", testDiskReadWrite},
		{"安装kubectl自动补全", installKubectlAutocompletion},
		{"安装curl", installCurl},
		{"返回上一级菜单", nil},
	}
	for {
		displayMenu("工具箱", toolboxOptions)
		choice := getUserChoice()
		if choice >= 0 && choice < len(toolboxOptions) {
			if toolboxOptions[choice].handler == nil {
				return nil
			}
			if err := toolboxOptions[choice].handler(); err != nil {
				fmt.Printf("错误: %v\n", err)
			}
		} else {
			fmt.Println("无效选项，请重试。")
		}
	}
}

func exit() error {
	fmt.Println("退出程序...")
	os.Exit(0)
	return nil
}

func displayMenu(title string, options []menuOption) {
	fmt.Printf("%s:\n", title)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i, option.title)
	}
}
func getUserChoice() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请选择操作：")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue // 如果输入为空，则继续提示用户重新输入
		}
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("无效输入，请输入数字。")
			continue // 如果输入无效，则继续提示用户重新输入
		}
		return choice
	}
}
