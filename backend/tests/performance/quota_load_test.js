#!/usr/bin/env k6
/*
 * ZKER Quota Check Performance Test
 *
 * 测试场景: 配额检查性能测试
 * 目标: 验证配额检查在高并发下的性能
 * 关键指标: 配额检查延迟应<100ms
 *
 * 使用方法:
 *   k6 run --env BASE_URL=http://localhost:8080 --env TENANT_ID=test-tenant quota_load_test.js
 */

import http from 'k6/http';
import { check } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('quota_errors');
const checkLatency = new Trend('quota_check_latency');

// 测试配置
export const options = {
    scenarios: {
        // 场景1: 恒定负载测试
        constant_load: {
            executor: 'constant-vus',
            vus: 100,              // 100个虚拟用户
            duration: '5m',        // 持续5分钟
            gracefulStop: '30s',
        },

        // 场景2: 爬坡测试
        ramping_load: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '1m', target: 50 },   // 1分钟爬坡到50
                { duration: '2m', target: 50 },   // 维持50用户2分钟
                { duration: '1m', target: 100 },  // 爬坡到100
                { duration: '2m', target: 100 },  // 维持100用户2分钟
                { duration: '1m', target: 200 },  // 爬坡到200
                { duration: '2m', target: 200 },  // 维持200用户2分钟
            ],
            gracefulStop: '30s',
        },

        // 场景3: 压力测试
        stress_test: {
            executor: 'ramping-arrival-rate',
            startRate: 100,
            timeUnit: '1s',
            preAllocatedVUs: 500,
            stages: [
                { duration: '2m', target: 200 },  // 每秒200请求
                { duration: '2m', target: 500 },  // 每秒500请求
                { duration: '2m', target: 1000 }, // 每秒1000请求
                { duration: '1m', target: 0 },    // 降到0
            ],
            gracefulStop: '30s',
        },
    },

    // 性能阈值
    thresholds: {
        http_req_duration: ['p(95)<100', 'p(99)<200'], // 配额检查应该很快
        http_req_failed: ['rate<0.01'],
        quota_errors: ['rate<0.01'],
        quota_check_latency: ['p(95)<100', 'p(99)<200'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080/api/v1';
const TENANT_ID = __ENV.TENANT_ID || 'test-tenant-id';
const AUTH_TOKEN = 'test-token';

export default function () {
    // 测试不同资源类型的配额检查
    const resourceTypes = [
        { type: 'bots', amount: 1 },
        { type: 'knowledge', amount: 1 },
        { type: 'workflows', amount: 1 },
        { type: 'api_calls', amount: 10 },
    ];

    resourceTypes.forEach(resource => {
        checkQuota(resource.type, resource.amount);
    });
}

// 检查配额
function checkQuota(resourceType, amount) {
    const startTime = Date.now();

    const resp = http.post(
        `${BASE_URL}/tenants/${TENANT_ID}/quotas/${resourceType}/check`,
        JSON.stringify({ amount }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${AUTH_TOKEN}`,
            },
            tags: { name: `CheckQuota_${resourceType}` },
        }
    );

    const endTime = Date.now();
    const duration = endTime - startTime;

    const success = check(resp, {
        [`quota check (${resourceType}) status is 200`]: (r) => r.status === 200,
        [`quota check (${resourceType}) latency < 50ms`]: (r) => duration < 50,
        [`quota check (${resourceType}) latency < 100ms`]: (r) => duration < 100,
        'response has allowed field': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && typeof body.data.allowed === 'boolean';
            } catch (e) {
                return false;
            }
        },
    });

    errorRate.add(!success, { tag: resourceType });
    checkLatency.add(duration, { tag: resourceType });

    if (!success) {
        console.error(`❌ Quota check failed for ${resourceType}: status=${resp.status}`);
    }
}

// 辅助函数: 获取配额使用情况
export function handleSummary(options) {
    console.log('\n========== 配额检查性能测试报告 ==========');
    console.log(`测试租户: ${TENANT_ID}`);
    console.log(`请求数: ${options.metrics.http_reqs.values.count}`);
    console.log(`P95延迟: ${options.metrics.http_req_duration.values['p(95)']}ms`);
    console.log(`P99延迟: ${options.metrics.http_req_duration.values['p(99)']}ms`);
    console.log(`错误率: ${(options.metrics.http_req_failed.values.rate * 100).toFixed(2)}%`);
    console.log('==========================================\n');
}
