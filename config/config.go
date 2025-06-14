package config

import (
	_ "embed"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DataDir       string `yaml:"dataDir"`
	RemoteBaseUrl string `yaml:"remoteBaseUrl"`
}

var GlobalConfig *Config

func defaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		DataDir:       filepath.Join(homeDir, ".lsx"), // 默认路径 ~/.lsx
		RemoteBaseUrl: "https://unpkg.com/linux-command@latest",
	}
}

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

func ParseConfig() *Config {
	cfg, _ := loadConfig()
	GlobalConfig = cfg
	return cfg
}
