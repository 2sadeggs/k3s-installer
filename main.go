package main

import (
	"fmt"
	"os"
)

func main() {
	// 检查是否为root用户
	if os.Geteuid() != 0 {
		fmt.Println("请以root用户身份运行此程序。")
		return
	}

	// 关闭交换空间
	err := disableSwap()
	if err != nil {
		fmt.Println("Error disabling swap space:", err)
	} else {
		fmt.Println("Swap space disabled successfully.")
	}

	// 关闭防火墙
	err = disableLinuxFirewall()
	if err != nil {
		fmt.Println("Error disabling firewall:", err)
	} else {
		fmt.Println("Firewall disabled successfully.")
	}

	// 检查其他系统要求
	runSystemChecks()

	// 主菜单
	mainMenu()
}
