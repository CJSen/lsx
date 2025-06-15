package config

import (
	_ "embed"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// 配置结构体，包含数据目录和远程数据源
// DataDir: 本地数据存储目录
// RemoteBaseUrl: 远程命令数据源地址
type Config struct {
	DataDir       string `yaml:"dataDir"`
	RemoteBaseUrl string `yaml:"remoteBaseUrl"`
	UseShow       bool   `yaml:"useshow"`
}

var GlobalConfig *Config // 全局配置实例

// 返回默认配置
func defaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		DataDir:       filepath.Join(homeDir, ".lsx"), // 默认路径 ~/.lsx
		RemoteBaseUrl: "https://unpkg.com/linux-command@latest",
		UseShow:       false,
	}
}

// 加载配置，优先环境变量 LSXCONFIG
func loadConfig() (*Config, string) {
	cfg := defaultConfig()
	configFile := os.Getenv("LSXCONFIG")
	if configFile == "" {
		configFile = "default config"
	} else {
		if _, err := os.Stat(configFile); err == nil {
			content, err := os.ReadFile(configFile)
			if err == nil {
				yaml.Unmarshal(content, cfg)
			}
		}
	}
	return cfg, configFile
}

// 解析配置并设置全局变量
func ParseConfig() *Config {
	cfg, _ := loadConfig()
	GlobalConfig = cfg
	return cfg
}
