package main

import (
	"fmt"
	"net/http"
	"time"
)

// 自定义k3s安装 -- 根据国内环境国外环境判断使用国内源或官方国外源
func autoInstallK3sByLocation(envArgs, cmdArgs []string) error {
	// 获取安装脚本 URL
	installURL, err := getInstallURL()
	fmt.Println("installURL...", installURL)
	if err != nil {
		return err
	}

	// 执行安装
	return k3sInstall(installURL, envArgs, cmdArgs)
}

// 获取安装脚本URL -- 默认看做国内环境 使用阿里源
func getInstallURL() (string, error) {
	// 判断是否在中国大陆
	if isInMainlandChina() {
		//return officialCNInstallURL, nil
		return aliyunInstallURL, nil
	}
	return officialInstallURL, nil
}

// 判断是否是国内环境
func isInMainlandChina() bool {
	// 检测当前系统是否可以访问互联网
	if !isInternetReachable() {
		return true // 如果无法访问互联网，假定为国内环境
	}

	// 尝试访问谷歌
	client := http.Client{
		Timeout: 3 * time.Second, // 设置 HTTP 请求超时时间
	}
	resp, err := client.Get("http://www.google.com")
	if err != nil {
		return true // 如果访问谷歌失败，假定为国内环境
	}
	defer resp.Body.Close()

	return false // 如果访问谷歌成功，假定为国外环境
}

// 检测当前系统是否可以访问互联网（通过访问百度）
func isInternetReachable() bool {
	client := http.Client{
		Timeout: 3 * time.Second, // 设置 HTTP 请求超时时间
	}

	resp, err := client.Get("http://www.baidu.com")
	if err != nil {
		return false // 如果访问百度失败，假定当前系统处于离线环境
	}
	defer resp.Body.Close()

	return true // 如果访问百度成功，说明当前系统可以访问互联网
}
