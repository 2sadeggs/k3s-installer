package main

import (
	"os"
	"os/exec"
)

// viewServiceStatus displays the status of all pods in the cluster.
func getPodStatus() error {
	cmd := exec.Command("kubectl", "get", "po", "-A", "-o", "wide")

	// 将命令输出连接到标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	//err := cmd.Run()
	//if err != nil {
	//	return fmt.Errorf("执行kubectl命令失败: %v", err)
	//}
	//
	//return nil
	return cmd.Run()
}
