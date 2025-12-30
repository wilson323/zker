// backend/domain/tenant/migration/gray_release_controller.go
// 灰度发布控制器 - 4阶段灰度切换
package migration

import (
	"context"
	"fmt"
	"sync"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// GrayReleaseController 灰度发布控制器
type GrayReleaseController struct {
	db             *gorm.DB
	logger         logs.CtxLogger
	currentStage   int
	stages         []ReleaseStage
	enabledTenants *TenantWhitelist
	metrics        *GrayReleaseMetrics
	mu             sync.RWMutex
}

// ReleaseStage 发布阶段
type ReleaseStage struct {
	StageID     int        `json:"stage_id"`
	Name        string     `json:"name"`
	Percentage  int        `json:"percentage"` // 流量百分比
	Description string     `json:"description"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Status      string     `json:"status"` // pending, active, completed, rolled_back
}

// GrayReleaseMetrics 灰度发布指标
type GrayReleaseMetrics struct {
	TotalRequests    int64     `json:"total_requests"`
	NewRouteRequests int64     `json:"new_route_requests"`
	OldRouteRequests int64     `json:"old_route_requests"`
	ErrorRate        float64   `json:"error_rate"`
	LatencyP95       float64   `json:"latency_p95"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewGrayReleaseController 创建灰度发布控制器
func NewGrayReleaseController(db *gorm.DB) *GrayReleaseController {
	return &GrayReleaseController{
		db:             db,
		logger:         logs.DefaultLogger(),
		currentStage:   0,
		enabledTenants: NewTenantWhitelist(db),
		metrics:        &GrayReleaseMetrics{},
		stages: []ReleaseStage{
			{StageID: 1, Name: "5%灰度", Percentage: 5, Description: "5%流量切换到新架构", Status: "pending"},
			{StageID: 2, Name: "20%灰度", Percentage: 20, Description: "20%流量切换到新架构", Status: "pending"},
			{StageID: 3, Name: "50%灰度", Percentage: 50, Description: "50%流量切换到新架构", Status: "pending"},
			{StageID: 4, Name: "100%全量", Percentage: 100, Description: "全量流量切换到新架构", Status: "pending"},
		},
	}
}

// StartGrayRelease 启动灰度发布
func (c *GrayReleaseController) StartGrayRelease(ctx context.Context) error {
	c.logger.CtxInfof(ctx, "[GrayRelease] ========== 启动灰度发布 ==========")

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.currentStage > 0 {
		return fmt.Errorf("灰度发布已在进行中，当前阶段: %d", c.currentStage)
	}

	c.currentStage = 1
	c.logger.CtxInfof(ctx, "[GrayRelease] 启动阶段1: 5%%灰度")

	return c.startStage(ctx, 1)
}

// startStage 启动指定阶段
func (c *GrayReleaseController) startStage(ctx context.Context, stageID int) error {
	stage := &c.stages[stageID-1]

	if stage.Status == "completed" {
		c.logger.CtxWarnf(ctx, "[GrayRelease] 阶段 %d 已完成", stageID)
		return nil
	}

	stage.Status = "active"
	now := time.Now()
	stage.StartedAt = &now

	c.logger.CtxInfof(ctx, "[GrayRelease] ========== 阶段%d: %s ==========", stageID, stage.Name)
	c.logger.CtxInfof(ctx, "[GrayRelease] 流量百分比: %d%%", stage.Percentage)
	c.logger.CtxInfof(ctx, "[GrayRelease] 描述: %s", stage.Description)

	// 记录指标
	c.metrics.UpdatedAt = time.Now()

	return nil
}

// PromoteStage 晋升到下一阶段
func (c *GrayReleaseController) PromoteStage(ctx context.Context, currentStageID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.currentStage != currentStageID {
		return fmt.Errorf("当前阶段不匹配: 期望=%d, 实际=%d", currentStageID, c.currentStage)
	}

	// 完成当前阶段
	c.stages[currentStageID-1].Status = "completed"
	now := time.Now()
	c.stages[currentStageID-1].CompletedAt = &now

	c.logger.CtxInfof(ctx, "[GrayRelease] ========== 阶段%d完成 ==========", currentStageID)
	c.logMetrics(ctx)

	// 检查是否可以晋升到下一阶段
	nextStageID := currentStageID + 1
	if nextStageID > len(c.stages) {
		c.logger.CtxInfof(ctx, "[GrayRelease] ========== 灰度发布全部完成！==========")
		c.currentStage = len(c.stages)
		return nil
	}

	// 启动下一阶段
	c.currentStage = nextStageID
	return c.startStage(ctx, nextStageID)
}

// Rollback 回滚到上一阶段
func (c *GrayReleaseController) Rollback(ctx context.Context, reason string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.CtxErrorf(ctx, "[GrayRelease] ========== 开始回滚 ==========")
	c.logger.CtxErrorf(ctx, "[GrayRelease] 回滚原因: %s", reason)

	currentStage := c.stages[c.currentStage-1]
	currentStage.Status = "rolled_back"

	// 回滚到上一阶段
	if c.currentStage > 1 {
		c.currentStage--
		c.stages[c.currentStage-1].Status = "active"
		c.logger.CtxInfof(ctx, "[GrayRelease] 已回滚到阶段 %d", c.currentStage)
	} else {
		c.logger.CtxWarnf(ctx, "[GrayRelease] 已是第一阶段，无法回滚")
		return fmt.Errorf("已是第一阶段，无法回滚")
	}

	c.logger.CtxErrorf(ctx, "[GrayRelease] ========== 回滚完成 ==========")

	return nil
}

// ShouldRouteToNew 判断是否应该路由到新架构
func (c *GrayReleaseController) ShouldRouteToNew(ctx context.Context, tenantID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 1. 检查白名单
	if c.enabledTenants != nil && c.enabledTenants.IsEnabled(tenantID) {
		return true
	}

	// 2. 检查当前阶段
	if c.currentStage == 0 {
		return false
	}

	stage := c.stages[c.currentStage-1]
	if stage.Status != "active" {
		return false
	}

	// 3. 根据流量百分比决定（简化实现）
	// 实际应该基于租户ID的哈希值或用户ID
	// 这里简化为随机模拟
	// return rand.Intn(100) < stage.Percentage

	return false // 默认返回false，由外部决定
}

// GetCurrentStage 获取当前阶段
func (c *GrayReleaseController) GetCurrentStage() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentStage
}

// GetStages 获取所有阶段
func (c *GrayReleaseController) GetStages() []ReleaseStage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stages
}

// RecordMetrics 记录指标
func (c *GrayReleaseController) RecordMetrics(isNewRoute bool, latency float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics.TotalRequests++

	if isNewRoute {
		c.metrics.NewRouteRequests++
	} else {
		c.metrics.OldRouteRequests++
	}

	c.metrics.UpdatedAt = time.Now()
}

// logMetrics 输出指标
func (c *GrayReleaseController) logMetrics(ctx context.Context) {
	metrics := c.metrics

	newRouteRate := float64(metrics.NewRouteRequests) / float64(metrics.TotalRequests) * 100

	c.logger.CtxInfof(ctx, "[GrayRelease] ========== 阶段指标 ==========")
	c.logger.CtxInfof(ctx, "[GrayRelease] 总请求数: %d", metrics.TotalRequests)
	c.logger.CtxInfof(ctx, "[GrayRelease] 新路由请求: %d (%.2f%%)", metrics.NewRouteRequests, newRouteRate)
	c.logger.CtxInfof(ctx, "[GrayRelease] 旧路由请求: %d (%.2f%%)", metrics.OldRouteRequests, 100-newRouteRate)
	c.logger.CtxInfof(ctx, "[GrayRelease] 错误率: %.2f%%", metrics.ErrorRate*100)
	c.logger.CtxInfof(ctx, "[GrayRelease] P95延迟: %.2fms", metrics.LatencyP95)
	c.logger.CtxInfof(ctx, "[GrayRelease] =================================")
}

// GetMetrics 获取指标
func (c *GrayReleaseController) GetMetrics() *GrayReleaseMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metrics
}
