package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 全局配置结构
type Config struct {
	AuthType    string `json:"auth_type"`     // api-key 或 token
	AuthValue   string `json:"auth_value"`    // API Key 或 JWT Token
	APIEndpoint string `json:"api_endpoint"`  // API端点
	CurrentBot  string `json:"current_bot"`   // 当前Bot ID
}

// globalConfig 全局配置实例
var globalConfig *Config

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".zker")
	configFile := filepath.Join(configDir, "config.json")

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		globalConfig = &Config{
			APIEndpoint: "http://localhost:8888",
		}
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	globalConfig = cfg
	return nil
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	if globalConfig == nil {
		return &Config{
			APIEndpoint: "http://localhost:8888",
		}
	}
	return globalConfig
}

// Save 保存配置到文件
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".zker")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsAuthenticated 检查是否已认证
func (c *Config) IsAuthenticated() bool {
	return c.AuthType != "" && c.AuthValue != ""
}

// GetAuthToken 获取认证Token
func (c *Config) GetAuthToken() string {
	return c.AuthValue
}

// GetAPIEndpoint 获取API端点
func (c *Config) GetAPIEndpoint() string {
	if c.APIEndpoint == "" {
		return "http://localhost:8888"
	}
	return c.APIEndpoint
}
