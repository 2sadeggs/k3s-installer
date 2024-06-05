package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type IPInfo struct {
	Query        string  `json:"query"`
	Country      string  `json:"country"`
	Region       string  `json:"regionName"`
	City         string  `json:"city"`
	ISP          string  `json:"isp"`
	Organization string  `json:"org"`
	AS           string  `json:"as"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lon"`
}

// 获取上网公网IP信息
func getPublicIPInfo() error {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return fmt.Errorf("无法获取公网IP信息: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("获取公网IP信息失败: HTTP 状态码 %d", resp.StatusCode)
	}

	var ipInfo IPInfo
	err = json.NewDecoder(resp.Body).Decode(&ipInfo)
	if err != nil {
		return fmt.Errorf("解析公网IP信息失败: %v", err)
	}

	fmt.Printf("公网IP地址:        %s\n", ipInfo.Query)
	fmt.Printf("国家:             %s\n", ipInfo.Country)
	fmt.Printf("地区:             %s\n", ipInfo.Region)
	fmt.Printf("城市:             %s\n", ipInfo.City)
	fmt.Printf("ISP提供商:        %s\n", ipInfo.ISP)
	fmt.Printf("组织:             %s\n", ipInfo.Organization)
	fmt.Printf("自治系统编号:       %s\n", ipInfo.AS)
	fmt.Printf("纬度:             %.6f\n", ipInfo.Latitude)
	fmt.Printf("经度:             %.6f\n", ipInfo.Longitude)

	return nil
}

// 获取主机上网公网IP raw格式
func getPublicIP() (string, error) {
	response, err := http.Get("https://checkip.amazonaws.com")
	//response, err := http.Get("https://www.trackip.net/ip")
	if err != nil {
		return "", errors.New("external IP fetch failed, detail:" + err.Error())
	}
	// 有很多类似网站提供这种服务，这是我知道且当前能用的
	/*
		curl https://myexternalip.com/raw
		curl ifconfig.me
		curl ip.sb
		curl ipinfo.io
		curl ip.cn
		curl cip.cc
		curl myip.ipip.net
		curl ifconfig.io
		curl httpbin.org/ip
		curl members.3322.org/dyndns/getip
		curl icanhazip.com
		curl www.trackip.net/ip
		curl checkip.amazonaws.com
	*/

	defer response.Body.Close()
	res := ""

	// 类似的API应当返回一个纯净的IP地址
	for {
		tmp := make([]byte, 32)
		n, err := response.Body.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return "", errors.New("external IP fetch failed, detail:" + err.Error())
			}
			res += string(tmp[:n])
			break
		}
		res += string(tmp[:n])
	}
	return strings.TrimSpace(res), nil
}

// 获取本地内网地址
func getLocalIP() (string, error) {
	// 获取本地IP地址
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	// 去除端口号
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip, nil
}
