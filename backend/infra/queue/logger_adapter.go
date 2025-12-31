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

package queue

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// nsqLoggerAdapter NSQ日志适配器
// 将zap.Logger适配到nsq.logger接口
type nsqLoggerAdapter struct {
	logger *zap.Logger
	level  nsq.LogLevel
}

// Output 实现nsq.logger接口
func (a *nsqLoggerAdapter) Output(calldepth int, s string) error {
	switch a.level {
	case nsq.LogLevelDebug:
		a.logger.Debug(s)
	case nsq.LogLevelInfo:
		a.logger.Info(s)
	case nsq.LogLevelWarning:
		a.logger.Warn(s)
	case nsq.LogLevelError:
		a.logger.Error(s)
	default:
		a.logger.Info(s)
	}
	return nil
}

// newNSQLogger 创建NSQ日志适配器
func newNSQLogger(logger *zap.Logger, level zapLogLevel) *nsqLoggerAdapter {
	return &nsqLoggerAdapter{
		logger: logger,
		level:  toNSQLogLevel(level),
	}
}

// zapLogLevel zap日志级别类型（用于类型转换）
type zapLogLevel int

// Zap到NSQ日志级别映射
const (
	LogLevelDebug zapLogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// toNSQLogLevel 转换zap日志级别到NSQ日志级别
func toNSQLogLevel(level zapLogLevel) nsq.LogLevel {
	switch level {
	case LogLevelDebug:
		return nsq.LogLevelDebug
	case LogLevelInfo:
		return nsq.LogLevelInfo
	case LogLevelWarning:
		return nsq.LogLevelWarning
	case LogLevelError:
		return nsq.LogLevelError
	default:
		return nsq.LogLevelInfo
	}
}
