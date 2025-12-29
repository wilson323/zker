#!/usr/bin/env k6
/*
 * ZKER API Stress Test
 *
 * 测试场景: API压力测试
 * 目标: 发现系统的性能瓶颈和极限
 *
 * 使用方法:
 *   k6 run --env BASE_URL=http://localhost:8080 stress_test.js
 */

import http from 'k6/http';
import { check, group } from 'k6';
import { Rate } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('stress_errors');

// 测试配置
export const options = {
    scenarios: {
        // 压力测试场景
        stress_test: {
            executor: 'ramping-arrival-rate',
            startRate: 10,
            timeUnit: '1s',
            preAllocatedVUs: 1000,
            maxVUs: 2000,
            stages: [
                { duration: '2m', target: 100 },   // 每秒100请求
                { duration: '2m', target: 500 },   // 每秒500请求
                { duration: '2m', target: 1000 },  // 每秒1000请求
                { duration: '2m', target: 2000 },  // 每秒2000请求
                { duration: '1m', target: 0 },     // 降到0
            ],
            gracefulStop: '30s',
        },
    },

    thresholds: {
        http_req_duration: ['p(95)<1000', 'p(99)<2000'],
        http_req_failed: ['rate<0.05'], // 压力测试允许5%错误率
        stress_errors: ['rate<0.05'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080/api/v1';
const AUTH_TOKEN = 'test-token';

// 测试租户ID列表(用于压力测试)
const TENANT_IDS = [
    'tenant-1',
    'tenant-2',
    'tenant-3',
    'tenant-4',
    'tenant-5',
];

export default function () {
    group('Read Operations', () => {
        // 随机选择一个租户ID
        const tenantID = TENANT_IDS[Math.floor(Math.random() * TENANT_IDS.length)];

        // 操作1: 获取租户列表
        const listResp = http.get(
            `${BASE_URL}/tenants?page_size=20`,
            {
                headers: { 'Authorization': `Bearer ${AUTH_TOKEN}` },
                tags: { name: 'ListTenants' },
            }
        );

        errorRate.add(listResp.status !== 200);

        // 操作2: 获取租户详情
        const detailResp = http.get(
            `${BASE_URL}/tenants/${tenantID}`,
            {
                headers: { 'Authorization': `Bearer ${AUTH_TOKEN}` },
                tags: { name: 'GetTenant' },
            }
        );

        errorRate.add(detailResp.status !== 200);

        // 操作3: 获取配额
        const quotaResp = http.get(
            `${BASE_URL}/tenants/${tenantID}/quotas`,
            {
                headers: { 'Authorization': `Bearer ${AUTH_TOKEN}` },
                tags: { name: 'GetQuotas' },
            }
        );

        errorRate.add(quotaResp.status !== 200);
    });

    group('Write Operations', () => {
        // 操作4: 创建租户(压力测试)
        const createResp = http.post(
            `${BASE_URL}/tenants`,
            JSON.stringify({
                tenant_name: `Stress Test Tenant ${Date.now()}`,
                tenant_type: 'individual',
                admin_email: `stresstest${Date.now()}@example.com`,
                admin_password: 'SecurePassword123!',
            }),
            {
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${AUTH_TOKEN}`,
                },
                tags: { name: 'CreateTenant' },
            }
        );

        errorRate.add(createResp.status !== 201 && createResp.status !== 409);
    });

    group('Mixed Operations', () => {
        const tenantID = TENANT_IDS[Math.floor(Math.random() * TENANT_IDS.length)];

        // 操作5: 配额检查
        const checkResp = http.post(
            `${BASE_URL}/tenants/${tenantID}/quotas/bots/check`,
            JSON.stringify({ amount: 1 }),
            {
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${AUTH_TOKEN}`,
                },
                tags: { name: 'CheckQuota' },
            }
        );

        errorRate.add(checkResp.status !== 200);
    });
}
