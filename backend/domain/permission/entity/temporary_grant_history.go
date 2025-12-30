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

package entity

import (
	"database/sql/driver"
	"encoding/json"
)

// ActionType 操作类型
type ActionType string

const (
	ActionTypeCreated ActionType = "created" // 创建
	ActionTypeUsed    ActionType = "used"    // 使用
	ActionTypeExpired ActionType = "expired" // 过期
	ActionTypeRevoked ActionType = "revoked" // 撤销
)

// TemporaryGrantHistory 临时授权历史实体
type TemporaryGrantHistory struct {
	HistoryID int64 `json:"history_id" gorm:"primaryKey;autoIncrement"`

	// 关联授权
	GrantID   int64  `json:"grant_id" gorm:"not null;index:idx_grant_id"`
	GrantCode string `json:"grant_code" gorm:"size:100;not null;index:idx_grant_code"`

	// 操作信息
	ActionType  ActionType `json:"action_type" gorm:"type:enum('created','used','expired','revoked');not null;index:idx_action_type"`
	OperatorID  string     `json:"operator_id" gorm:"size:36;not null;index:idx_operator_id"` // 操作人ID
	OperatorName string     `json:"operator_name" gorm:"size:100"`                            // 操作人姓名

	// 快照数据
	SnapshotData SnapshotData `json:"snapshot_data" gorm:"type:json"` // 操作时的权限数据快照

	// 元数据
	Reason    string `json:"reason,omitempty" gorm:"size:500"`   // 操作原因
	RequestID string `json:"request_id,omitempty" gorm:"size:36"` // 关联请求ID

	// 时间戳
	CreatedAt int64 `json:"created_at" gorm:"not null;default:0;index:idx_created_at"`
}

// TableName 指定表名
func (TemporaryGrantHistory) TableName() string {
	return "temporary_grant_history"
}

// SnapshotData 快照数据（支持JSON序列化）
type SnapshotData map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (sd *SnapshotData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, sd)
}

// Value 实现 driver.Valuer 接口
func (sd SnapshotData) Value() (driver.Value, error) {
	if len(sd) == 0 {
		return "{}", nil
	}

	return json.Marshal(sd)
}

// GetActionTypeText 获取操作类型的中文描述
func (th *TemporaryGrantHistory) GetActionTypeText() string {
	switch th.ActionType {
	case ActionTypeCreated:
		return "创建授权"
	case ActionTypeUsed:
		return "使用授权"
	case ActionTypeExpired:
		return "授权过期"
	case ActionTypeRevoked:
		return "撤销授权"
	default:
		return "未知操作"
	}
}
