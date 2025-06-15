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
	UseLess       bool   `yaml:"useLess"`
}

var GlobalConfig *Config // 全局配置实例

// 返回默认配置
func defaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		DataDir:       filepath.Join(homeDir, ".lsx"), // dataDir: "/Users/css/.lsx" # 默认为"~/.lsx",自定义请写完整目录路径（自定义时不要使用～，结尾不要有/）
		RemoteBaseUrl: "https://unpkg.com/linux-command@latest",
		UseLess:       false,
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
