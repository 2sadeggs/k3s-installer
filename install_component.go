package main

import (
	"embed"
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"os/exec"
	"strings"
	"time"
)

//go:embed yaml
var yamlFS embed.FS

// Constants
const (
	yamlDirectory = "yaml"
)

// installComponent installs a Kubernetes component using kubectl apply.
func installComponent(componentName string) error {
	// 显示等待旋转动画
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("Deploying %s...", componentName)
	s.Start()
	defer s.Stop()

	yamlPath := fmt.Sprintf("%s/%s.yaml", yamlDirectory, componentName)
	yamlBytes, err := yamlFS.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("读取文件错误: %v", err)
	}
	yamlConfig := string(yamlBytes)

	cmd := exec.Command("sh", "-c", fmt.Sprintf("sudo kubectl apply -f -"))
	cmd.Stdin = strings.NewReader(yamlConfig)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
