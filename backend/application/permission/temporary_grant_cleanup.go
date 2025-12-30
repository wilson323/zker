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

package permission

import (
	"context"
	"log"
	"time"
)

// TemporaryGrantCleanupJob 临时授权清理任务
// 定期清理过期的临时授权和历史记录
type TemporaryGrantCleanupJob struct {
	temporaryGrantService *TemporaryGrantService
	cleanupInterval       time.Duration // 清理间隔
	historyRetentionDays  int           // 历史记录保留天数
}

// NewTemporaryGrantCleanupJob 创建清理任务实例
func NewTemporaryGrantCleanupJob(
	temporaryGrantService *TemporaryGrantService,
	cleanupInterval time.Duration,
	historyRetentionDays int,
) *TemporaryGrantCleanupJob {
	if cleanupInterval == 0 {
		cleanupInterval = 1 * time.Hour // 默认1小时清理一次
	}
	if historyRetentionDays == 0 {
		historyRetentionDays = 90 // 默认保留90天
	}

	return &TemporaryGrantCleanupJob{
		temporaryGrantService: temporaryGrantService,
		cleanupInterval:       cleanupInterval,
		historyRetentionDays:  historyRetentionDays,
	}
}

// Start 启动清理任务
func (j *TemporaryGrantCleanupJob) Start(ctx context.Context) {
	log.Printf("[临时授权清理任务] 启动，清理间隔: %v, 历史保留天数: %d", j.cleanupInterval, j.historyRetentionDays)

	ticker := time.NewTicker(j.cleanupInterval)
	defer ticker.Stop()

	// 立即执行一次清理
	j.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[临时授权清理任务] 停止")
			return
		case <-ticker.C:
			j.cleanup(ctx)
		}
	}
}

// cleanup 执行清理
func (j *TemporaryGrantCleanupJob) cleanup(ctx context.Context) {
	startTime := time.Now()
	log.Println("[临时授权清理任务] 开始执行清理...")

	// 1. 清理过期的临时授权
	count, err := j.temporaryGrantService.CleanupExpiredGrants(ctx)
	if err != nil {
		log.Printf("[临时授权清理任务] 清理过期授权失败: %v", err)
	} else {
		log.Printf("[临时授权清理任务] 清理了 %d 个过期临时授权", count)
	}

	// 2. 清理过期的历史记录（保留指定天数）
	historyBefore := time.Now().AddDate(0, 0, -j.historyRetentionDays).UnixMilli()
	historyCount, err := j.temporaryGrantService.CleanupExpiredHistory(ctx, historyBefore)
	if err != nil {
		log.Printf("[临时授权清理任务] 清理过期历史失败: %v", err)
	} else {
		log.Printf("[临时授权清理任务] 清理了 %d 条过期历史记录", historyCount)
	}

	duration := time.Since(startTime)
	log.Printf("[临时授权清理任务] 清理完成，耗时: %v", duration)
}
