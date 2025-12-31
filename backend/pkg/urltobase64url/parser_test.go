/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package urltobase64url

import (
	"strings"
	"testing"
)

// TestValidateURL 测试URL验证功能（SSRF防护）
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errMsg  string
	}{
		// ✅ 合法URL
		{
			name:    "valid http URL",
			url:     "http://example.com/file.txt",
			wantErr: false,
		},
		{
			name:    "valid https URL",
			url:     "https://example.com/file.txt",
			wantErr: false,
		},
		{
			name:    "valid http URL with port",
			url:     "http://example.com:8080/file.txt",
			wantErr: false,
		},
		{
			name:    "valid https URL with path",
			url:     "https://api.example.com/v1/files/document.pdf",
			wantErr: false,
		},

		// ❌ 非法协议
		{
			name:    "invalid ftp URL",
			url:     "ftp://example.com/file.txt",
			wantErr: true,
			errMsg:  "only http and https schemes are allowed",
		},
		{
			name:    "invalid file URL",
			url:     "file:///etc/passwd",
			wantErr: true,
			errMsg:  "only http and https schemes are allowed",
		},
		{
			name:    "invalid javascript URL",
			url:     "javascript:alert(1)",
			wantErr: true,
			errMsg:  "only http and https schemes are allowed",
		},

		// ❌ localhost访问
		{
			name:    "localhost access denied",
			url:     "http://localhost:8080/file.txt",
			wantErr: true,
			errMsg:  "access to localhost is not allowed",
		},
		{
			name:    "127.0.0.1 access denied",
			url:     "http://127.0.0.1:8080/file.txt",
			wantErr: true,
			errMsg:  "access to localhost is not allowed",
		},
		{
			name:    "0.0.0.0 access denied",
			url:     "http://0.0.0.0:8080/file.txt",
			wantErr: true,
			errMsg:  "access to localhost is not allowed",
		},
		{
			name:    "IPv6 localhost denied",
			url:     "http://::1/file.txt",
			wantErr: true,
			errMsg:  "access to localhost is not allowed",
		},

		// ❌ 私有IP访问
		{
			name:    "private IP 10.0.0.0/8 denied",
			url:     "http://10.0.0.1/file.txt",
			wantErr: true,
			errMsg:  "access to private IP is not allowed",
		},
		{
			name:    "private IP 172.16.0.0/12 denied",
			url:     "http://172.16.0.1/file.txt",
			wantErr: true,
			errMsg:  "access to private IP is not allowed",
		},
		{
			name:    "private IP 192.168.0.0/16 denied",
			url:     "http://192.168.1.1/file.txt",
			wantErr: true,
			errMsg:  "access to private IP is not allowed",
		},
		{
			name:    "link-local IP 169.254.0.0/16 denied",
			url:     "http://169.254.1.1/file.txt",
			wantErr: true,
			errMsg:  "access to private IP is not allowed",
		},

		// ❌ 云元数据服务访问
		{
			name:    "AWS metadata endpoint denied",
			url:     "http://169.254.169.254/latest/meta-data/",
			wantErr: true,
			errMsg:  "access to metadata endpoint is not allowed",
		},
		{
			name:    "阿里云元数据服务denied",
			url:     "http://100.100.100.200/latest/meta-data/",
			wantErr: true,
			errMsg:  "access to metadata endpoint is not allowed",
		},

		// ❌ 无效URL
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
			errMsg:  "invalid URL",
		},
		{
			name:    "malformed URL",
			url:     "http://",
			wantErr: true,
			errMsg:  "invalid URL",
		},
		{
			name:    "URL without hostname",
			url:     "http:///file.txt",
			wantErr: true,
			errMsg:  "empty hostname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				errStr := err.Error()
				// 使用strings.Contains而不是字符串切片比较，避免越界
				if !strings.Contains(errStr, tt.errMsg) {
					t.Errorf("validateURL() error message = %v, want to contain %v", errStr, tt.errMsg)
				}
			}
		})
	}
}

// TestIsLocalhost 测试localhost检测
func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{"localhost", "localhost", true},
		{"LOCALHOST", "LOCALHOST", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"0.0.0.0", "0.0.0.0", true},
		{"::1", "::1", true},
		{"127.0.0.2", "127.0.0.2", true}, // 回环地址其他IP
		{"example.com", "example.com", false},
		{"8.8.8.8", "8.8.8.8", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLocalhost(tt.host); got != tt.want {
				t.Errorf("isLocalhost() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsPrivateIP 测试私有IP检测
func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		// 私有IP
		{"10.0.0.1", "10.0.0.1", true},
		{"10.255.255.255", "10.255.255.255", true},
		{"172.16.0.1", "172.16.0.1", true},
		{"172.31.255.255", "172.31.255.255", true},
		{"192.168.0.1", "192.168.0.1", true},
		{"192.168.255.255", "192.168.255.255", true},
		{"169.254.1.1", "169.254.1.1", true}, // 链路本地

		// 公网IP
		{"8.8.8.8", "8.8.8.8", false},
		{"1.1.1.1", "1.1.1.1", false},
		{"114.114.114.114", "114.114.114.114", false},

		// 域名（非IP）
		{"example.com", "example.com", false},
		{"localhost", "localhost", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrivateIP(tt.host); got != tt.want {
				t.Errorf("isPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsMetadataEndpoint 测试云元数据端点检测
func TestIsMetadataEndpoint(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{"AWS metadata", "169.254.169.254", true},
		{"阿里云元数据", "100.100.100.200", true},
		{"public IP", "8.8.8.8", false},
		{"localhost", "localhost", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMetadataEndpoint(tt.host); got != tt.want {
				t.Errorf("isMetadataEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestURLToBase64_InvalidURL 测试无效URL应被拒绝
func TestURLToBase64_InvalidURL(t *testing.T) {
	invalidURLs := []string{
		"http://localhost/file.txt",
		"http://127.0.0.1/file.txt",
		"http://10.0.0.1/file.txt",
		"http://169.254.169.254/file.txt",
		"ftp://example.com/file.txt",
		"file:///etc/passwd",
	}

	for _, url := range invalidURLs {
		t.Run(url, func(t *testing.T) {
			_, err := URLToBase64(url)
			if err == nil {
				t.Errorf("URLToBase64(%s) should return error, got nil", url)
			}
		})
	}
}
