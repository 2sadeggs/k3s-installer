package main

import (
	"fmt"
	"github.com/fatih/color"
)

const (
	appPort     = "32000"
	kuboardPort = "30080"
	kuboardUser = "admin"
)

// getAppLoginInfo displays the application login information.
func getAppLoginInfo() error {
	// 获取内网IP和外网IP
	internalIP, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("获取内网IP失败：%v", err)
	}

	externalIP, err := getPublicIP()
	if err != nil {
		return fmt.Errorf("获取外网IP失败：%v", err)
	}

	// 获取 Kuboard 密码
	kuboardPassword, err := getKuboardAdminDefaultPassword()
	if err != nil {
		return fmt.Errorf("获取Kuboard密码失败：%v", err)
	}

	// 使用 color 包来设置终端输出颜色
	green := color.New(color.FgGreen).SprintFunc()

	// 内网APP登录信息绿色显示
	fmt.Println(green("==================== APP登录信息 ===================="))
	fmt.Println(green(fmt.Sprintf("登录地址:        %s:%s", internalIP, appPort)))
	fmt.Println(green(fmt.Sprintf("Kuboard登录地址: %s:%s", internalIP, kuboardPort)))
	fmt.Println(green(fmt.Sprintf("Kuboard登录用户: %s", kuboardUser)))
	fmt.Println(green(fmt.Sprintf("Kuboard登录密码: %s", kuboardPassword)))
	fmt.Println(green("======================================================"))

	// 外网APP登录信息正常显示
	fmt.Println("================= APP外网登录信息 =====================")
	fmt.Printf("登录地址:        %s:%s\n", externalIP, appPort)
	fmt.Printf("Kuboard登录地址: %s:%s\n", externalIP, kuboardPort)
	fmt.Printf("Kuboard登录用户: %s\n", kuboardUser)
	fmt.Printf("Kuboard登录密码: %s\n", kuboardPassword)
	fmt.Println("(注意：外网IP是程序根据部署服务器上网的公网IP自动生成，不一定准确，仅供参考)")
	fmt.Println("======================================================")

	return nil
}
