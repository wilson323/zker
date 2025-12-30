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

package health

// ==================== Response DTOs ====================

// HealthResponse 基本健康检查响应
type HealthResponse struct {
	Status    struct {
		Database       bool `json:"database"`
		Redis          bool `json:"redis"`
		Elasticsearch  bool `json:"elasticsearch"`
		MinIO          bool `json:"minio"`
	} `json:"status"`
	Timestamp int64 `json:"timestamp"`
}

// LiveResponse 存活检查响应
type LiveResponse struct {
	Status  string `json:"status"` // always "ok"
	Message string `json:"message,omitempty"`
}

// ReadyResponse 就绪检查响应
type ReadyResponse struct {
	Status string `json:"status"` // "ready" or "not_ready"`
	Checks struct {
		Database       bool `json:"database"`
		Redis          bool `json:"redis"`
		Elasticsearch  bool `json:"elasticsearch"`
		MinIO          bool `json:"minio"`
	} `json:"checks"`
}

// ComponentStatus 组件状态
type ComponentStatus struct {
	Status     string  `json:"status"`     // "healthy", "degraded", "unhealthy"
	LatencyMs  int64   `json:"latency_ms"`
	Message    string  `json:"message,omitempty"`
	Error      string  `json:"error,omitempty"`
}

// DetailedHealthResponse 详细健康检查响应
type DetailedHealthResponse struct {
	Status     string                    `json:"status"`     // "healthy", "degraded", "unhealthy"
	Timestamp  int64                     `json:"timestamp"`
	Version    string                    `json:"version,omitempty"`
	Components map[string]ComponentStatus `json:"components"`
}
