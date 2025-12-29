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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// DBWriteQPS 数据库写QPS
	DBWriteQPS = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_write_qps_total",
			Help: "Database write queries per second",
		},
		[]string{"database"}, // database: zker
	)

	// DBReadQPS 数据库读QPS
	DBReadQPS = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_read_qps_total",
			Help: "Database read queries per second",
		},
		[]string{"database", "replica"}, // replica: slave1, slave2
	)

	// DBReplicationLag 数据库复制延迟
	DBReplicationLag = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_replication_lag_seconds",
			Help: "Database replication lag in seconds",
		},
		[]string{"master", "slave"}, // master: mysql-master, slave: mysql-slave1, mysql-slave2
	)

	// DBWriteLatency 数据库写操作延迟
	DBWriteLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_write_latency_seconds",
			Help:    "Database write operation latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"database", "operation"}, // operation: insert, update, delete
	)

	// DBReadLatency 数据库读操作延迟
	DBReadLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_read_latency_seconds",
			Help:    "Database read operation latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"database", "replica", "operation"}, // operation: select
	)

	// DBConnectionPoolUsage 数据库连接池使用率
	DBConnectionPoolUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connection_pool_usage",
			Help: "Database connection pool usage (0-1)",
		},
		[]string{"database", "type"}, // type: write, read
	)

	// DBReplicationStatus 数据库复制状态
	DBReplicationStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_replication_status",
			Help: "Database replication status (1=running, 0=stopped)",
		},
		[]string{"master", "slave"},
	)
)

// ========== 辅助函数 ==========

// RecordDBWrite 记录数据库写操作
func RecordDBWrite(database, operation string, duration float64) {
	DBWriteQPS.WithLabelValues(database).Inc()
	DBWriteLatency.WithLabelValues(database, operation).Observe(duration)
}

// RecordDBRead 记录数据库读操作
func RecordDBRead(database, replica, operation string, duration float64) {
	DBReadQPS.WithLabelValues(database, replica).Inc()
	DBReadLatency.WithLabelValues(database, replica, operation).Observe(duration)
}

// UpdateReplicationLag 更新复制延迟
func UpdateReplicationLag(master, slave string, lagSeconds float64) {
	DBReplicationLag.WithLabelValues(master, slave).Set(lagSeconds)
}

// UpdateReplicationStatus 更新复制状态
func UpdateReplicationStatus(master, slave string, status float64) {
	DBReplicationStatus.WithLabelValues(master, slave).Set(status)
}

// UpdateConnectionPoolUsage 更新连接池使用率
func UpdateConnectionPoolUsage(database, dbType string, usage float64) {
	DBConnectionPoolUsage.WithLabelValues(database, dbType).Set(usage)
}
