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

package database

import (
	"gorm.io/gorm"
)

var (
	// globalDB 全局数据库实例
	globalDB *gorm.DB
)

// InitDB 初始化数据库实例
func InitDB(db *gorm.DB) {
	globalDB = db
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return globalDB
}
