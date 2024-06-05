package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	officialInstallURL   = "https://get.k3s.io"
	officialCNInstallURL = "https://rancher-mirror.rancher.cn/k3s/k3s-install.sh"
	aliyunInstallURL     = "https://rancher-mirror.oss-cn-beijing.aliyuncs.com/k3s/k3s-install.sh"
	defaultDataDir       = "/var/lib/rancher/k3s"
	defaultBinDir        = "/usr/local/bin"
	defaultVersion       = ""
)

// k3s安装函数 -- 除了环境变量和双短横线参数还考虑了安装源因素
func k3sInstall(installURL string, envArgs, cmdArgs []string) error {
	fmt.Printf("Downloading K3s install script from %s...\n", installURL)

	resp, err := http.Get(installURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Installing K3s... "
	s.Start()
	defer s.Stop()

	cmd := exec.Command("/bin/sh")
	cmd.Stdin = resp.Body
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 判断是否为国内源或阿里源，添加 INSTALL_K3S_MIRROR=cn 环境变量
	if installURL == officialCNInstallURL || installURL == aliyunInstallURL {
		envArgs = append(envArgs, "INSTALL_K3S_MIRROR=cn")
	}

	// 设置环境变量
	cmd.Env = append(os.Environ(), envArgs...)
	//fmt.Println("cmd.Env...", cmd.Env)

	// 执行shell命令,传递双短横线参数
	cmd.Args = append(cmd.Args, "-c", fmt.Sprintf("sh -s - %s", strings.Join(cmdArgs, " ")))
	//fmt.Println("cmd.Args...", cmd.Args)

	// 打印环境变量和命令行参数
	fmt.Println("Environment variables:", cmd.Env)
	fmt.Println("Command arguments:", cmd.Args)

	fmt.Println("Starting to execute K3s install script...")

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
