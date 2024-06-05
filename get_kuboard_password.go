package main

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

//go:embed yaml/app.yaml
var embeddedFiles embed.FS

type ConfigMap struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

func getKuboardAdminDefaultPassword() (string, error) {
	// 读取嵌入式文件内容
	fileData, err := embeddedFiles.ReadFile("yaml/app.yaml")
	if err != nil {
		return "", err
	}

	// 按照 "---" 分隔符拆分 YAML 文件内容
	yamlParts := strings.Split(string(fileData), "---")

	// 遍历每个 YAML 部分
	for _, yamlPart := range yamlParts {
		if strings.TrimSpace(yamlPart) == "" {
			continue
		}

		// 解析单个 YAML 部分为 map[string]interface{}
		var data map[string]interface{}
		if err := yaml.Unmarshal([]byte(yamlPart), &data); err != nil {
			return "", err
		}

		// 提取 Kind 和 Metadata 字段
		kind, kindOk := data["kind"].(string)
		metadata, metadataOk := data["metadata"].(map[string]interface{})

		// 检查是否为 ConfigMap 并且名称为 kuboard-v3-config
		if kindOk && kind == "ConfigMap" && metadataOk {
			name, nameOk := metadata["name"].(string)
			if nameOk && name == "kuboard-v3-config" {
				// 提取数据字段
				dataMap, dataOk := data["data"].(map[string]interface{})
				if dataOk {
					// 提取密码字段
					if password, ok := dataMap["KUBOARD_ADMIN_DEFAULT_PASSWORD"].(string); ok {
						return strings.TrimSpace(password), nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("KUBOARD_ADMIN_DEFAULT_PASSWORD not found in kuboard-v3-config")
}
