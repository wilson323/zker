#!/usr/bin/env k6
/*
 * ZKER 组织中心性能测试套件
 *
 * 测试场景: 组织中心API完整性能测试
 * 目标: 验证组织管理在高并发下的性能表现
 *
 * API覆盖:
 *   - GET /api/organizations (组织列表)
 *   - GET /api/organizations/:id (组织详情)
 *   - POST /api/organizations (创建组织)
 *   - PUT /api/organizations/:id (更新组织)
 *   - GET /api/organizations/tree (组织树)
 *   - GET /api/org/departments (部门列表)
 *   - GET /api/org/employees (员工列表)
 *   - GET /api/org/employees/search (员工搜索)
 *   - GET /api/directory/organization (通讯录)
 *
 * 使用方法:
 *   # 基准测试 (100用户, 5分钟)
 *   k6 run --env BASE_URL=http://localhost:8080 org_load_test.k6.js
 *
 *   # 峰值测试 (1000用户, 短时间冲击)
 *   k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=spike org_load_test.k6.js
 *
 *   # 压力测试 (500用户, 30分钟)
 *   k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=stress org_load_test.k6.js
 *
 *   # 耐久测试 (200用户, 2小时)
 *   k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=endurance org_load_test.k6.js
 *
 *   # 可扩展性测试 (阶梯式负载)
 *   k6 run --env BASE_URL=http://localhost:8080 --env TEST_TYPE=scalability org_load_test.k6.js
 *
 * @author 研发B (后端工程师)
 * @version 1.0
 * @date 2025-01-01
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// ================================
// 自定义指标
// ================================
const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');
const throughput = new Trend('throughput');

// API特定指标
const orgListRate = new Rate('org_list_errors');
const orgDetailRate = new Rate('org_detail_errors');
const orgCreateRate = new Rate('org_create_errors');
const orgUpdateRate = new Rate('org_update_errors');
const orgTreeRate = new Rate('org_tree_errors');
const deptListRate = new Rate('dept_list_errors');
const empListRate = new Rate('emp_list_errors');
const empSearchRate = new Rate('emp_search_errors');
const directoryRate = new Rate('directory_errors');

// ================================
// 测试配置
// ================================
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TEST_TYPE = __ENV.TEST_TYPE || 'base'; // base, spike, stress, endurance, scalability
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'test-token-perf-';

// 测试场景配置
const testConfigs = {
    base: {
        // 基准测试: 100用户, 持续5分钟
        stages: [
            { duration: '1m', target: 10 },   // 预热
            { duration: '3m', target: 100 },  // 爬坡到100用户
            { duration: '5m', target: 100 },  // 维持100用户5分钟
            { duration: '1m', target: 0 },    // 降温
        ],
        thresholds: {
            http_req_duration: ['p(95)<200', 'p(99)<500'],
            http_req_failed: ['rate<0.001'],
            errors: ['rate<0.001'],
        },
    },
    spike: {
        // 峰值测试: 1000用户, 短时间冲击
        stages: [
            { duration: '2m', target: 100 },   // 预热到100用户
            { duration: '1m', target: 1000 },  // 突增到1000用户
            { duration: '3m', target: 1000 },  // 维持1000用户3分钟
            { duration: '2m', target: 0 },     // 快速降温
        ],
        thresholds: {
            http_req_duration: ['p(95)<500', 'p(99)<1000'],
            http_req_failed: ['rate<0.01'],
            errors: ['rate<0.01'],
        },
    },
    stress: {
        // 压力测试: 500用户, 持续30分钟
        stages: [
            { duration: '5m', target: 50 },    // 预热
            { duration: '5m', target: 100 },   // 爬坡到100
            { duration: '5m', target: 200 },   // 爬坡到200
            { duration: '5m', target: 300 },   // 爬坡到300
            { duration: '5m', target: 400 },   // 爬坡到400
            { duration: '30m', target: 500 },  // 维持500用户30分钟
            { duration: '5m', target: 0 },     // 降温
        ],
        thresholds: {
            http_req_duration: ['p(95)<300', 'p(99)<800'],
            http_req_failed: ['rate<0.005'],
            errors: ['rate<0.005'],
        },
    },
    endurance: {
        // 耐久测试: 200用户, 持续2小时
        stages: [
            { duration: '5m', target: 50 },    // 预热
            { duration: '5m', target: 100 },   // 爬坡
            { duration: '5m', target: 200 },   // 爬坡到目标负载
            { duration: '2h', target: 200 },   // 维持200用户2小时
            { duration: '5m', target: 0 },     // 降温
        ],
        thresholds: {
            http_req_duration: ['p(95)<250', 'p(99)<600'],
            http_req_failed: ['rate<0.002'],
            errors: ['rate<0.002'],
        },
    },
    scalability: {
        // 可扩展性测试: 阶梯式增加负载
        stages: [
            { duration: '3m', target: 10 },    // 10用户
            { duration: '3m', target: 50 },    // 50用户
            { duration: '3m', target: 100 },   // 100用户
            { duration: '3m', target: 200 },   // 200用户
            { duration: '3m', target: 300 },   // 300用户
            { duration: '3m', target: 500 },   // 500用户
            { duration: '3m', target: 700 },   // 700用户
            { duration: '3m', target: 1000 },  // 1000用户
            { duration: '5m', target: 0 },     // 降温
        ],
        thresholds: {
            http_req_duration: ['p(95)<400', 'p(99)<1000'],
            http_req_failed: ['rate<0.01'],
            errors: ['rate<0.01'],
        },
    },
};

// 导出配置
export const options = function() {
    const config = testConfigs[TEST_TYPE] || testConfigs.base;
    return {
        stages: config.stages,
        thresholds: config.thresholds,
        // 全局设置
        noConnectionReuse: false,
        userAgent: 'K6-Performance-Test/1.0',
        // 超时设置
        timeout: '30s',
        // 批量设置
        batch: 20,
        batchPerHost: 10,
    };
}();

// ================================
// 测试数据
// ================================
let testData = {
    tenantID: null,
    orgID: null,
    deptID: null,
    empID: null,
    authToken: AUTH_TOKEN,
};

// ================================
// Setup: 测试初始化
// ================================
export function setup() {
    console.log('🚀 组织中心性能测试 - 初始化测试环境');
    console.log(`📊 测试类型: ${TEST_TYPE}`);
    console.log(`🌐 BASE_URL: ${BASE_URL}`);

    // 1. 创建测试租户
    const tenantResp = http.post(
        `${BASE_URL}/api/tenants`,
        JSON.stringify({
            tenant_name: `OrgPerfTest_${Date.now()}`,
            tenant_type: 'enterprise',
            admin_email: `orgperftest${Date.now()}@example.com`,
            admin_password: 'SecurePassword123!',
        }),
        {
            headers: { 'Content-Type': 'application/json' },
            tags: { name: 'Setup_CreateTenant' },
        }
    );

    if (tenantResp.status !== 201 && tenantResp.status !== 200) {
        console.error(`❌ 创建测试租户失败: ${tenantResp.status}`);
        console.error(`响应: ${tenantResp.body}`);
        return null;
    }

    const tenant = JSON.parse(tenantResp.body).data;
    console.log(`✅ 创建测试租户成功: ${tenant.tenant_id}`);

    // 2. 创建测试组织
    const orgResp = http.post(
        `${BASE_URL}/api/organizations`,
        JSON.stringify({
            tenant_id: tenant.tenant_id,
            org_name: '性能测试组织',
            org_type: 'company',
            org_code: 'PERF_ORG_001',
            description: '性能测试专用组织',
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'X-Tenant-ID': tenant.tenant_id,
            },
            tags: { name: 'Setup_CreateOrganization' },
        }
    );

    if (orgResp.status !== 201 && orgResp.status !== 200) {
        console.warn(`⚠️ 创建测试组织失败: ${orgResp.status}, 将使用查询获取`);
        // 尝试获取现有组织
        const listResp = http.get(
            `${BASE_URL}/api/organizations?tenant_id=${tenant.tenant_id}&page_size=1`,
            {
                headers: { 'X-Tenant-ID': tenant.tenant_id },
                tags: { name: 'Setup_ListOrganizations' },
            }
        );

        if (listResp.status === 200) {
            const orgs = JSON.parse(listResp.body).data.organizations;
            if (orgs && orgs.length > 0) {
                testData.orgID = orgs[0].org_id;
                console.log(`✅ 使用现有组织: ${testData.orgID}`);
            }
        }
    } else {
        const org = JSON.parse(orgResp.body).data;
        testData.orgID = org.org_id;
        console.log(`✅ 创建测试组织成功: ${org.org_id}`);
    }

    // 3. 创建测试部门
    const deptResp = http.post(
        `${BASE_URL}/api/org/departments`,
        JSON.stringify({
            tenant_id: tenant.tenant_id,
            org_id: testData.orgID,
            dept_name: '性能测试部门',
            dept_code: 'PERF_DEPT_001',
            description: '性能测试专用部门',
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'X-Tenant-ID': tenant.tenant_id,
            },
            tags: { name: 'Setup_CreateDepartment' },
        }
    );

    if (deptResp.status === 201 || deptResp.status === 200) {
        const dept = JSON.parse(deptResp.body).data;
        testData.deptID = dept.dept_id;
        console.log(`✅ 创建测试部门成功: ${dept.dept_id}`);
    }

    // 4. 获取测试员工
    const empListResp = http.get(
        `${BASE_URL}/api/org/employees?tenant_id=${tenant.tenant_id}&page_size=1`,
        {
            headers: { 'X-Tenant-ID': tenant.tenant_id },
            tags: { name: 'Setup_ListEmployees' },
        }
    );

    if (empListResp.status === 200) {
        const emps = JSON.parse(empListResp.body).data.employees;
        if (emps && emps.length > 0) {
            testData.empID = emps[0].emp_id;
            console.log(`✅ 使用现有员工: ${testData.empID}`);
        }
    }

    console.log('📊 测试数据初始化完成');
    console.log(`  - 租户ID: ${tenant.tenant_id}`);
    console.log(`  - 组织ID: ${testData.orgID}`);
    console.log(`  - 部门ID: ${testData.deptID}`);
    console.log(`  - 员工ID: ${testData.empID}`);

    return {
        tenant_id: tenant.tenant_id,
        org_id: testData.orgID,
        dept_id: testData.deptID,
        emp_id: testData.empID,
    };
}

// ================================
// Default: 主测试逻辑
// ================================
export default function(data) {
    if (!data || !data.tenant_id) {
        console.error('❌ 测试数据未初始化');
        return;
    }

    // 执行所有测试场景（随机权重）
    const scenarios = [
        { func: () => scenario1_listOrganizations(data), weight: 20 },
        { func: () => scenario2_getOrganization(data.org_id), weight: 15 },
        { func: () => scenario3_getOrganizationTree(data), weight: 10 },
        { func: () => scenario4_listDepartments(data), weight: 15 },
        { func: () => scenario5_listEmployees(data), weight: 20 },
        { func: () => scenario6_searchEmployees(data), weight: 10 },
        { func: () => scenario7_getDirectory(data), weight: 10 },
    ];

    // 随机选择场景执行
    const totalWeight = scenarios.reduce((sum, s) => sum + s.weight, 0);
    let random = Math.random() * totalWeight;

    for (const scenario of scenarios) {
        random -= scenario.weight;
        if (random <= 0) {
            scenario.func();
            break;
        }
    }

    // 每个VU每秒执行一次
    sleep(1);
}

// ================================
// Teardown: 清理测试数据
// ================================
export function teardown(data) {
    if (!data || !data.tenant_id) {
        return;
    }

    console.log('🧹 清理测试数据');

    // 删除测试租户（会级联删除所有相关数据）
    const deleteResp = http.del(
        `${BASE_URL}/api/tenants/${data.tenant_id}`,
        null,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'Teardown_DeleteTenant' },
        }
    );

    if (deleteResp.status === 204 || deleteResp.status === 200) {
        console.log('✅ 测试租户已删除');
    } else {
        console.warn(`⚠️ 删除测试租户失败: ${deleteResp.status}`);
    }

    console.log('✅ 测试清理完成');
}

// ================================
// 测试场景
// ================================

/**
 * 场景1: 获取组织列表
 * API: GET /api/organizations
 * 性能目标: p95 < 200ms
 */
function scenario1_listOrganizations(data) {
    const params = {
        page_size: 20,
        page_number: Math.floor(Math.random() * 10) + 1,
    };

    const resp = http.get(
        `${BASE_URL}/api/organizations?${new URLSearchParams(params)}`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'ListOrganizations' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
        'has organizations array': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.organizations);
            } catch (e) {
                return false;
            }
        },
    });

    orgListRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ ListOrganizations failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景2: 获取组织详情
 * API: GET /api/organizations/:id
 * 性能目标: p95 < 100ms
 */
function scenario2_getOrganization(orgID) {
    if (!orgID) {
        return;
    }

    const resp = http.get(
        `${BASE_URL}/api/organizations/${orgID}`,
        {
            headers: { 'X-Tenant-ID': testData.tenant_id },
            tags: { name: 'GetOrganization' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 100ms': (r) => r.timings.duration < 100,
        'has org_id': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && body.data.org_id === orgID;
            } catch (e) {
                return false;
            }
        },
    });

    orgDetailRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ GetOrganization failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景3: 获取组织树
 * API: GET /api/organizations/tree
 * 性能目标: p95 < 300ms (树形查询较复杂)
 */
function scenario3_getOrganizationTree(data) {
    const resp = http.get(
        `${BASE_URL}/api/organizations/tree`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'GetOrganizationTree' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 300ms': (r) => r.timings.duration < 300,
        'has tree structure': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.tree);
            } catch (e) {
                return false;
            }
        },
    });

    orgTreeRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ GetOrganizationTree failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景4: 获取部门列表
 * API: GET /api/org/departments
 * 性能目标: p95 < 200ms
 */
function scenario4_listDepartments(data) {
    const params = {
        page_size: 20,
        page_number: Math.floor(Math.random() * 10) + 1,
    };

    const resp = http.get(
        `${BASE_URL}/api/org/departments?${new URLSearchParams(params)}`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'ListDepartments' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
        'has departments array': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.departments);
            } catch (e) {
                return false;
            }
        },
    });

    deptListRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ ListDepartments failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景5: 获取员工列表
 * API: GET /api/org/employees
 * 性能目标: p95 < 200ms
 * 说明: 员工表数据量大，是最常见的查询场景
 */
function scenario5_listEmployees(data) {
    const params = {
        page_size: 20,
        page_number: Math.floor(Math.random() * 50) + 1,
    };

    const resp = http.get(
        `${BASE_URL}/api/org/employees?${new URLSearchParams(params)}`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'ListEmployees' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
        'has employees array': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.employees);
            } catch (e) {
                return false;
            }
        },
    });

    empListRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ ListEmployees failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景6: 搜索员工
 * API: GET /api/org/employees/search
 * 性能目标: p95 < 300ms
 * 说明: 全文搜索是性能瓶颈点
 */
function scenario6_searchEmployees(data) {
    const searchTerms = ['张', '李', '王', '工程师', '经理', '开发', '测试'];

    const params = {
        keyword: searchTerms[Math.floor(Math.random() * searchTerms.length)],
        page_size: 10,
    };

    const resp = http.get(
        `${BASE_URL}/api/org/employees/search?${new URLSearchParams(params)}`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'SearchEmployees' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 300ms': (r) => r.timings.duration < 300,
        'has search results': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && Array.isArray(body.data.employees);
            } catch (e) {
                return false;
            }
        },
    });

    empSearchRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ SearchEmployees failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景7: 获取通讯录
 * API: GET /api/directory/organization
 * 性能目标: p95 < 500ms
 * 说明: 通讯录聚合查询，涉及多表JOIN
 */
function scenario7_getDirectory(data) {
    const resp = http.get(
        `${BASE_URL}/api/directory/organization`,
        {
            headers: { 'X-Tenant-ID': data.tenant_id },
            tags: { name: 'GetDirectory' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
        'has directory data': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && (body.data.organizations || body.data.departments);
            } catch (e) {
                return false;
            }
        },
    });

    directoryRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ GetDirectory failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景8: 创建组织 (写操作测试)
 * API: POST /api/organizations
 * 性能目标: p95 < 300ms
 * 说明: 测试写操作性能（树形结构计算较重）
 */
function scenario8_createOrganization(data) {
    const timestamp = Date.now();
    const orgData = {
        tenant_id: data.tenant_id,
        org_name: `PerfTest组织_${timestamp}`,
        org_type: 'department',
        org_code: `PERF_ORG_${timestamp}`,
        description: '性能测试创建的组织',
    };

    const resp = http.post(
        `${BASE_URL}/api/organizations`,
        JSON.stringify(orgData),
        {
            headers: {
                'Content-Type': 'application/json',
                'X-Tenant-ID': data.tenant_id,
            },
            tags: { name: 'CreateOrganization' },
        }
    );

    const success = check(resp, {
        'status is 201 or 200': (r) => r.status === 201 || r.status === 200,
        'response time < 300ms': (r) => r.timings.duration < 300,
        'has org_id': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && body.data.org_id;
            } catch (e) {
                return false;
            }
        },
    });

    orgCreateRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ CreateOrganization failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}

/**
 * 场景9: 更新组织
 * API: PUT /api/organizations/:id
 * 性能目标: p95 < 250ms
 */
function scenario9_updateOrganization(data) {
    if (!data.org_id) {
        return;
    }

    const updateData = {
        org_name: `更新组织_${Date.now()}`,
        description: '性能测试更新的组织',
    };

    const resp = http.put(
        `${BASE_URL}/api/organizations/${data.org_id}`,
        JSON.stringify(updateData),
        {
            headers: {
                'Content-Type': 'application/json',
                'X-Tenant-ID': data.tenant_id,
            },
            tags: { name: 'UpdateOrganization' },
        }
    );

    const success = check(resp, {
        'status is 200': (r) => r.status === 200,
        'response time < 250ms': (r) => r.timings.duration < 250,
    });

    orgUpdateRate.add(!success);
    errorRate.add(!success);
    responseTime.add(resp.timings.duration);

    if (!success) {
        console.error(`❌ UpdateOrganization failed: status=${resp.status}, time=${resp.timings.duration}ms`);
    }
}
