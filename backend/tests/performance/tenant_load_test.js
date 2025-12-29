#!/usr/bin/env k6
/*
 * ZKER Tenant API Load Test
 *
 * 测试场景: 租户管理API负载测试
 * 目标: 验证API在高并发下的性能表现
 *
 * 使用方法:
 *   k6 run --env BASE_URL=http://localhost:8080 tenant_load_test.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const latency = new Trend('latency');

// 测试配置
export const options = {
    // 阶段性负载测试
    stages: [
        { duration: '1m', target: 10 },   // 1分钟爬坡到10用户
        { duration: '3m', target: 10 },   // 维持10用户3分钟
        { duration: '1m', target: 50 },   // 1分钟爬坡到50用户
        { duration: '3m', target: 50 },   // 维持50用户3分钟
        { duration: '1m', target: 100 },  // 1分钟爬坡到100用户
        { duration: '5m', target: 100 },  // 维持100用户5分钟
        { duration: '1m', target: 200 },  // 1分钟爬坡到200用户
        { duration: '5m', target: 200 },  // 维持200用户5分钟
        { duration: '2m', target: 0 },    // 2分钟降到0
    ],

    // 性能阈值
    thresholds: {
        http_req_duration: ['p(95)<500', 'p(99)<1000'],  // 95%请求<500ms, 99%请求<1000ms
        http_req_failed: ['rate<0.01'],                   // 错误率<1%
        errors: ['rate<0.01'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080/api/v1';

// 测试数据
let testData = {
    tenantID: null,
    authToken: 'test-token', // 实际测试时应该通过登录获取
};

// setup: 在测试开始前执行
export function setup() {
    console.log('🚀 开始性能测试 - 设置测试数据');

    // 创建测试租户
    const createResp = http.post(
        `${BASE_URL}/tenants`,
        JSON.stringify({
            tenant_name: 'Performance Test Tenant',
            tenant_type: 'enterprise',
            admin_email: 'perftest@example.com',
            admin_password: 'SecurePassword123!',
        }),
        {
            headers: { 'Content-Type': 'application/json' },
            tags: { name: 'CreateTenant' },
        }
    );

    if (createResp.status !== 201) {
        console.error('❌ 创建测试租户失败:', createResp.status);
        return null;
    }

    const tenant = JSON.parse(createResp.body).data;
    console.log(`✅ 创建测试租户成功: ${tenant.tenant_id}`);

    return {
        tenant_id: tenant.tenant_id,
    };
}

// default: 主测试逻辑
export default function (data) {
    if (!data || !data.tenant_id) {
        console.error('❌ 测试数据未初始化');
        return;
    }

    // 场景1: 获取租户列表
    scenario1_listTenants();

    // 场景2: 获取租户详情
    scenario2_getTenant(data.tenant_id);

    // 场景3: 获取租户配额
    scenario3_getQuotas(data.tenant_id);

    // 场景4: 检查Bot配额
    scenario4_checkBotQuota(data.tenant_id);

    // 场景5: 获取配额列表
    scenario5_listQuotas(data.tenant_id);

    // 每个VU每秒执行一次
    sleep(1);
}

// teardown: 在测试结束后执行
export function teardown(data) {
    if (!data || !data.tenant_id) {
        return;
    }

    console.log('🧹 清理测试数据');

    // 删除测试租户
    const deleteResp = http.del(
        `${BASE_URL}/tenants/${data.tenant_id}`,
        null,
        {
            headers: { 'Authorization': `Bearer ${testData.authToken}` },
            tags: { name: 'DeleteTenant' },
        }
    );

    if (deleteResp.status === 204) {
        console.log('✅ 测试租户已删除');
    } else {
        console.warn('⚠️ 删除测试租户失败:', deleteResp.status);
    }
}

// ========== 测试场景 ==========

// 场景1: 获取租户列表
function scenario1_listTenants() {
    const resp = http.get(
        `${BASE_URL}/tenants?page_size=20`,
        {
            headers: { 'Authorization': `Bearer ${testData.authToken}` },
            tags: { name: 'ListTenants' },
        }
    );

    const success = check(resp, {
        'list status is 200': (r) => r.status === 200,
        'list response time < 500ms': (r) => r.timings.duration < 500,
        'list has tenants array': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.tenants);
            } catch (e) {
                return false;
            }
        },
    });

    errorRate.add(!success);
    latency.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ ListTenants failed: status=${resp.status}`);
    }
}

// 场景2: 获取租户详情
function scenario2_getTenant(tenantID) {
    const resp = http.get(
        `${BASE_URL}/tenants/${tenantID}`,
        {
            headers: { 'Authorization': `Bearer ${testData.authToken}` },
            tags: { name: 'GetTenant' },
        }
    );

    const success = check(resp, {
        'detail status is 200': (r) => r.status === 200,
        'detail response time < 300ms': (r) => r.timings.duration < 300,
        'detail has tenant_id': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && body.data.tenant_id === tenantID;
            } catch (e) {
                return false;
            }
        },
    });

    errorRate.add(!success);
    latency.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ GetTenant failed: status=${resp.status}`);
    }
}

// 场景3: 获取租户配额
function scenario3_getQuotas(tenantID) {
    const resp = http.get(
        `${BASE_URL}/tenants/${tenantID}/quotas`,
        {
            headers: { 'Authorization': `Bearer ${testData.authToken}` },
            tags: { name: 'GetQuotas' },
        }
    );

    const success = check(resp, {
        'quotas status is 200': (r) => r.status === 200,
        'quotas response time < 300ms': (r) => r.timings.duration < 300,
        'quotas has array': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.quotas);
            } catch (e) {
                return false;
            }
        },
    });

    errorRate.add(!success);
    latency.add(resp.timings.duration);
}

// 场景4: 检查Bot配额
function scenario4_checkBotQuota(tenantID) {
    const resp = http.post(
        `${BASE_URL}/tenants/${tenantID}/quotas/bots/check`,
        JSON.stringify({
            amount: 1,
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${testData.authToken}`,
            },
            tags: { name: 'CheckBotQuota' },
        }
    );

    const success = check(resp, {
        'check status is 200': (r) => r.status === 200,
        'check response time < 100ms': (r) => r.timings.duration < 100,
        'check has allowed field': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && typeof body.data.allowed === 'boolean';
            } catch (e) {
                return false;
            }
        },
    });

    errorRate.add(!success);
    latency.add(resp.timings.duration);
}

// 场景5: 获取特定资源配额
function scenario5_listQuotas(tenantID) {
    const resourceTypes = ['bots', 'knowledge', 'workflows', 'api_calls'];

    resourceTypes.forEach(resourceType => {
        const resp = http.get(
            `${BASE_URL}/tenants/${tenantID}/quotas/${resourceType}`,
            {
                headers: { 'Authorization': `Bearer ${testData.authToken}` },
                tags: { name: `GetQuota_${resourceType}` },
            }
        );

        const success = check(resp, {
            [`quota (${resourceType}) status is 200`]: (r) => r.status === 200,
            [`quota (${resourceType}) response time < 200ms`]: (r) => r.timings.duration < 200,
        });

        errorRate.add(!success);
        latency.add(resp.timings.duration);
    });
}
