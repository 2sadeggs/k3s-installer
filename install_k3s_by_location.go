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
	if isInMainlandChina() { // 如果是中国大陆返回 aliyunInstallURL 其他返回 officialInstallURL
		//return officialCNInstallURL, nil
		return aliyunInstallURL, nil
	}
	return officialInstallURL, nil
}

// 判断是否是国内环境
func isInMainlandChina() bool {
	// 如果无法访问互联网 也就是程序无法访问到百度 假定为国内环境
	if !isInternetReachable("http://www.baidu.com") {
		return true
	}
	// 如果程序能访问到谷歌 结果为 true 再取反为 false 表示不是中国大陆
	// 相反 如果不能访问谷歌 很容易理解 就是中国大陆
	return !isInternetReachable("http://www.google.com")
}

// isInternetReachable 检测当前系统是否可以访问互联网
func isInternetReachable(url string) bool {
	client := http.Client{
		Timeout: 3 * time.Second, // 设置 HTTP 请求超时时间
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}
