package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitGlobalConfig(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "zker-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 备份原配置
	homeDir, _ := os.UserHomeDir()
	originalConfigDir := filepath.Join(homeDir, ".zker")
	backupConfigDir := filepath.Join(tmpDir, "backup")

	if _, err := os.Stat(originalConfigDir); err == nil {
		os.Rename(originalConfigDir, backupConfigDir)
		defer os.Rename(backupConfigDir, originalConfigDir)
	}

	// 测试初始化
	err = InitGlobalConfig()
	if err != nil {
		t.Errorf("InitGlobalConfig failed: %v", err)
	}

	cfg := GetGlobalConfig()
	if cfg == nil {
		t.Error("GetGlobalConfig should not return nil")
	}

	if cfg.APIEndpoint == "" {
		t.Error("Default APIEndpoint should not be empty")
	}
}

func TestConfigSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zker-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		AuthType:    "api-key",
		AuthValue:   "sk-test-key",
		APIEndpoint: "http://localhost:8888",
	}

	// 测试保存（会保存到真实目录，但会被覆盖）
	err = cfg.Save()
	if err != nil {
		t.Errorf("Config.Save failed: %v", err)
	}
}

func TestConfigIsAuthenticated(t *testing.T) {
	tests := []struct {
		name      string
		authType  string
		authValue string
		want      bool
	}{
		{
			name:      "authenticated with api-key",
			authType:  "api-key",
			authValue: "sk-test",
			want:      true,
		},
		{
			name:      "authenticated with token",
			authType:  "token",
			authValue: "eyJhbGci",
			want:      true,
		},
		{
			name:      "not authenticated",
			authType:  "",
			authValue: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				AuthType:  tt.authType,
				AuthValue: tt.authValue,
			}
			if got := cfg.IsAuthenticated(); got != tt.want {
				t.Errorf("Config.IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigGetAPIEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "custom endpoint",
			endpoint: "https://api.example.com",
			want:     "https://api.example.com",
		},
		{
			name:     "empty endpoint returns default",
			endpoint: "",
			want:     "http://localhost:8888",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				APIEndpoint: tt.endpoint,
			}
			if got := cfg.GetAPIEndpoint(); got != tt.want {
				t.Errorf("Config.GetAPIEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
