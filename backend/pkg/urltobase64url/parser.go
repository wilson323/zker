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
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type FileData struct {
	Base64Url string
	MimeType  string
}

// validateURL 验证URL安全性，防止SSRF攻击
func validateURL(rawURL string) error {
	// 检查空URL
	if rawURL == "" {
		return fmt.Errorf("empty URL")
	}

	// 解析URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// 只允许 http 和 https 协议
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("only http and https schemes are allowed, got: %s", u.Scheme)
	}

	// 获取主机名
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("empty hostname")
	}

	// 检查是否是云元数据服务（优先检查，因为169.254.169.254在169.254.0.0/16范围内）
	if isMetadataEndpoint(host) {
		return fmt.Errorf("access to metadata endpoint is not allowed: %s", host)
	}

	// 检查是否是localhost或本地回环地址
	if isLocalhost(host) {
		return fmt.Errorf("access to localhost is not allowed: %s", host)
	}

	// 检查是否是私有IP地址
	if isPrivateIP(host) {
		return fmt.Errorf("access to private IP is not allowed: %s", host)
	}

	return nil
}

// isLocalhost 检查是否是localhost
func isLocalhost(host string) bool {
	// 检查常见localhost表示
	localhostNames := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
	}

	lowerHost := strings.ToLower(host)
	for _, name := range localhostNames {
		if lowerHost == name {
			return true
		}
	}

	// 检查127.0.0.0/8网段（回环地址）
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return true
	}

	return false
}

// isPrivateIP 检查是否是私有IP地址
func isPrivateIP(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	// 检查私有IP地址段
	privateIPBlocks := []string{
		"10.0.0.0/8",     // RFC1918 - Class A私有网络
		"172.16.0.0/12",  // RFC1918 - Class B私有网络
		"192.168.0.0/16", // RFC1918 - Class C私有网络
		"169.254.0.0/16", // RFC3927 - 链路本地地址
		"fc00::/7",       // RFC4193 - 唯一本地地址（IPv6）
		"fe80::/10",      // RFC4291 - 链路本地地址（IPv6）
	}

	for _, block := range privateIPBlocks {
		_, cidr, _ := net.ParseCIDR(block)
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}

// isMetadataEndpoint 检查是否是云服务商元数据端点
func isMetadataEndpoint(host string) bool {
	metadataEndpoints := []string{
		"169.254.169.254", // AWS/GCP/Azure元数据服务
		"100.100.100.200", // 阿里云元数据服务
	}

	for _, endpoint := range metadataEndpoints {
		if host == endpoint {
			return true
		}
	}

	return false
}

func URLToBase64(rawURL string) (*FileData, error) {
	// ✅ 安全验证：防止SSRF攻击
	if err := validateURL(rawURL); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("http get error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code error: %d", resp.StatusCode)
	}

	fileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read file content error: %v", err)
	}

	var mimeType string

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err == nil && mediaType != "" {
			mimeType = mediaType
		}
	}

	if mimeType == "" {
		detectedType := http.DetectContentType(fileContent)
		if detectedType != "application/octet-stream" {
			mimeType = detectedType
		}
	}

	if mimeType == "" || mimeType == "application/octet-stream" {
		urlPath := rawURL
		if idx := strings.Index(urlPath, "?"); idx != -1 {
			urlPath = urlPath[:idx]
		}
		if idx := strings.Index(urlPath, "#"); idx != -1 {
			urlPath = urlPath[:idx]
		}

		ext := filepath.Ext(urlPath)
		if ext != "" {
			extMimeType := mime.TypeByExtension(ext)
			if extMimeType != "" {
				mimeType = extMimeType
			}
		}
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	base64Str := base64.StdEncoding.EncodeToString(fileContent)

	return &FileData{
		Base64Url: "data:" + mimeType + ";base64," + base64Str,
		MimeType:  mimeType,
	}, nil
}
