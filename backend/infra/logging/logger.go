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

package logging

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/coze-studio/coze-studio/backend/pkg/config"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Init 初始化日志系统
func Init(env string) error {
	var zapConfig zap.Config

	if env == "production" {
		// 生产环境配置: JSON格式,文件输出,日志轮转
		zapConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding: "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths: []string{
				"/var/log/zker/app.log",
				"stdout",
			},
			ErrorOutputPaths: []string{
				"/var/log/zker/error.log",
				"stderr",
			},
			InitialFields: map[string]interface{}{
				"service": "zker-api",
				"env":     env,
			},
		}
	} else {
		// 开发环境配置: Console格式,彩色输出,Debug级别
		zapConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			Development: true,
			Encoding:    "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	// 构建logger
	var err error
	logger, err = zapConfig.Build(
		zap.AddCaller(), // 添加调用者信息
		zap.AddStacktrace(zap.ErrorLevel), // Error级别及以上添加堆栈
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	sugar = logger.Sugar()

	return nil
}

// InitWithFileRotation 使用日志轮转初始化
func InitWithFileRotation(env string, configFile string) error {
	var zapConfig zap.Config

	// 日志轮转配置
	rotateConfig := &lumberjack.Logger{
		Filename:   "/var/log/zker/app.log",
		MaxSize:    100, // MB
		MaxBackups: 10,  // 保留10个备份
		MaxAge:     30,  // 保留30天
		Compress:   true, // 压缩旧文件
	}

	if env == "production" {
		zapConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		// 使用 lumberjack 作为 WriteSyncer
		writer := zapcore.AddSync(rotateConfig)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
			writer,
			zapConfig.Level,
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		// 开发环境不需要日志轮转
		return Init(env)
	}

	sugar = logger.Sugar()
	return nil
}

// ========== 基础日志方法 ==========

// Debug 记录Debug日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 记录Info日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 记录Warn日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 记录Error日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal 记录Fatal日志并退出
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Panic 记录Panic日志并panic
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// ========== Sugared日志方法(更方便但性能较低) ==========

// Debugf 记录Debug日志(格式化)
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Infof 记录Info日志(格式化)
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Warnf 记录Warn日志(格式化)
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Errorf 记录Error日志(格式化)
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Fatalf 记录Fatal日志并退出(格式化)
func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

// ========== 带Context的日志方法 ==========

// WithContext 从context中提取信息并创建logger
func WithContext(ctx context.Context) *zap.Logger {
	fields := extractContextFields(ctx)
	return logger.With(fields...)
}

// WithContextSugar 从context中提取信息并创建sugar logger
func WithContextSugar(ctx context.Context) *zap.SugaredLogger {
	fields := extractContextFields(ctx)
	return sugar.With(fields...)
}

// DebugContext 记录带context的Debug日志
func DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Debug(msg, fields...)
}

// InfoContext 记录带context的Info日志
func InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Info(msg, fields...)
}

// WarnContext 记录带context的Warn日志
func WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Warn(msg, fields...)
}

// ErrorContext 记录带context的Error日志
func ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Error(msg, fields...)
}

// ========== 辅助函数 ==========

// extractContextFields 从context中提取字段
func extractContextFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	// 从context中提取request_id
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields = append(fields, zap.String("request_id", id))
		}
	}

	// 从context中提取user_id
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			fields = append(fields, zap.String("user_id", id))
		}
	}

	// 从context中提取tenant_id
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		if id, ok := tenantID.(string); ok {
			fields = append(fields, zap.String("tenant_id", id))
		}
	}

	// 从context中提取trace_id
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if id, ok := traceID.(string); ok {
			fields = append(fields, zap.String("trace_id", id))
		}
	}

	return fields
}

// GetLogger 获取底层logger
func GetLogger() *zap.Logger {
	return logger
}

// GetSugaredLogger 获取sugar logger
func GetSugaredLogger() *zap.SugaredLogger {
	return sugar
}

// Sync 同步日志缓冲区
func Sync() error {
	return logger.Sync()
}

// ========== 常用字段辅助函数 ==========

// String 创建string字段
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int 创建int字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 创建int64字段
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Float64 创建float64字段
func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// Bool 创建bool字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Any 创建任意类型字段
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Err 创建error字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration 创建duration字段
func Duration(key string, value interface{}) zap.Field {
	return zap.Duration(key, value)
}

// ========== HTTP日志辅助方法 ==========

// LogHTTPRequest 记录HTTP请求日志
func LogHTTPRequest(ctx context.Context, method, path, query, clientIP, userAgent string) {
	WithContext(ctx).Info("HTTP Request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent),
	)
}

// LogHTTPResponse 记录HTTP响应日志
func LogHTTPResponse(ctx context.Context, statusCode int, duration float64, responseSize int) {
	WithContext(ctx).Info("HTTP Response",
		zap.Int("status_code", statusCode),
		zap.Float64("duration_ms", duration),
		zap.Int("response_size", responseSize),
	)
}

// LogHTTPError 记录HTTP错误日志
func LogHTTPError(ctx context.Context, err error, statusCode int) {
	WithContext(ctx).Error("HTTP Error",
		zap.Error(err),
		zap.Int("status_code", statusCode),
	)
}

// ========== 数据库日志辅助方法 ==========

// LogDBQuery 记录数据库查询日志
func LogDBQuery(ctx context.Context, db, table, operation string, duration float64) {
	WithContext(ctx).Debug("DB Query",
		zap.String("database", db),
		zap.String("table", table),
		zap.String("operation", operation),
		zap.Float64("duration_ms", duration),
	)
}

// LogDBError 记录数据库错误日志
func LogDBError(ctx context.Context, db, operation string, err error) {
	WithContext(ctx).Error("DB Error",
		zap.String("database", db),
		zap.String("operation", operation),
		zap.Error(err),
	)
}

// ========== 业务日志辅助方法 ==========

// LogQuotaCheck 记录配额检查日志
func LogQuotaCheck(ctx context.Context, resourceType string, amount int, allowed bool, current, limit int) {
	WithContext(ctx).Info("Quota Check",
		zap.String("resource_type", resourceType),
		zap.Int("amount", amount),
		zap.Bool("allowed", allowed),
		zap.Int("current_usage", current),
		zap.Int("quota_limit", limit),
	)
}

// LogTenantOperation 记录租户操作日志
func LogTenantOperation(ctx context.Context, operation, tenantID string, success bool) {
	WithContext(ctx).Info("Tenant Operation",
		zap.String("operation", operation),
		zap.String("tenant_id", tenantID),
		zap.Bool("success", success),
	)
}

// LogSubscriptionChange 记录订阅变更日志
func LogSubscriptionChange(ctx context.Context, tenantID, operation, fromTier, toTier string) {
	WithContext(ctx).Info("Subscription Change",
		zap.String("tenant_id", tenantID),
		zap.String("operation", operation),
		zap.String("from_tier", fromTier),
		zap.String("to_tier", toTier),
	)
}

// ========== 性能日志辅助方法 ==========

// LogSlowQuery 记录慢查询日志
func LogSlowQuery(ctx context.Context, query string, duration float64, threshold float64) {
	if duration > threshold {
		WithContext(ctx).Warn("Slow Query",
			zap.String("query", query),
			zap.Float64("duration_ms", duration),
			zap.Float64("threshold_ms", threshold),
		)
	}
}

// LogSlowAPI 记录慢API日志
func LogSlowAPI(ctx context.Context, endpoint string, duration float64, threshold float64) {
	if duration > threshold {
		WithContext(ctx).Warn("Slow API",
			zap.String("endpoint", endpoint),
			zap.Float64("duration_ms", duration),
			zap.Float64("threshold_ms", threshold),
		)
	}
}
