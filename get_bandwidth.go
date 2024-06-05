package main

import (
	"fmt"
	"github.com/showwin/speedtest-go/speedtest"
)

// 检测主机上下行带宽并返回错误信息
func speedTestBandwidth() error {
	var speedtestClient = speedtest.New()

	// 获取服务器列表
	fmt.Println("Fetching server list...")
	serverList, err := speedtestClient.FetchServers()
	if err != nil {
		return fmt.Errorf("Error fetching server list: %v", err)
	}
	fmt.Println("Available servers: ", serverList.Available())

	// 选择测试目标
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return fmt.Errorf("Error finding server: %v", err)
	}

	for i, s := range targets {
		// 进行延迟测试
		fmt.Printf("[%d/%d] Testing latency for %s (%s)... ", i+1, len(targets), s.Sponsor, s.Name)
		err := s.PingTest(nil)
		if err != nil {
			return fmt.Errorf("Error testing latency: %v", err)
		}
		fmt.Printf("Done. Latency: %s\n", s.Latency)

		// 下载测试
		fmt.Printf("[%d/%d] Testing download speed for %s (%s)... ", i+1, len(targets), s.Sponsor, s.Name)
		err = s.DownloadTest()
		if err != nil {
			return fmt.Errorf("Error testing download speed: %v", err)
		}
		fmt.Printf("Done. Download: %f Mbps\n", s.DLSpeed)

		// 上传测试
		fmt.Printf("[%d/%d] Testing upload speed for %s (%s)... ", i+1, len(targets), s.Sponsor, s.Name)
		err = s.UploadTest()
		if err != nil {
			return fmt.Errorf("Error testing upload speed: %v", err)
		}
		fmt.Printf("Done. Upload: %f Mbps\n", s.ULSpeed)

		// 输出结果
		fmt.Printf("[%d/%d] Results for %s (%s):\n", i+1, len(targets), s.Sponsor, s.Name)
		fmt.Printf("Latency: %s, Download: %f Mbps, Upload: %f Mbps\n\n", s.Latency, s.DLSpeed, s.ULSpeed)

		// 重置计数器
		s.Context.Reset()
	}

	return nil
}
