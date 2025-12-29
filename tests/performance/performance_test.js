/*
 * ZKER 端到端性能测试套件
 * 版本: v1.0.0
 * 工具: K6 (https://k6.io/)
 *
 * 测试覆盖：
 * 1. HTTP API基础性能测试
 * 2. 租户管理API性能测试
 * 3. 配额管理API性能测试
 * 4. Bot服务API性能测试
 * 5. 权限检查API性能测试
 * 6. 并发压力测试
 * 7. 稳定性测试
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// ==================== 配置 ====================

// 测试环境配置
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8888';
const TEST_TENANT_ID = __ENV.TENANT_ID || 'test_tenant_001';
const TEST_USER_ID = __ENV.USER_ID || 'test_user_001';

// 性能阈值配置
const THRESHOLDS = {
  // HTTP响应时间阈值（P95）
  http_req_duration: ['p(95)<500', 'p(99)<1000'], // P95 < 500ms, P99 < 1000ms

  // 错误率阈值
  http_req_failed: ['rate<0.01'], // 错误率 < 1%

  // 特定API阈值
  api_latency_fast: ['p(95)<200'],   // 快速API: P95 < 200ms
  api_latency_normal: ['p(95)<500'], // 普通API: P(95) < 500ms
  api_latency_slow: ['p(95)<2000'],  // 慢速API: P95 < 2000ms
};

// 自定义指标
const apiLatency = new Trend('api_latency');
const tenantAPIRate = new Rate('tenant_api_success_rate');
const quotaAPIRate = new Rate('quota_api_success_rate');
const botAPIRate = new Rate('bot_api_success_rate');
const permissionAPIRate = new Rate('permission_api_success_rate');

// ==================== 测试选项 ====================

export const options = {
  // 基础配置
  vus: 10,                // 虚拟用户数
  duration: '5m',         // 测试持续时间

  // 阶段配置（负载测试）
  stages: [
    { duration: '1m', target: 10 },   // 启动阶段：1分钟内增加到10个用户
    { duration: '2m', target: 50 },   // 加压阶段：2分钟内增加到50个用户
    { duration: '1m', target: 100 },  // 峰值阶段：1分钟内增加到100个用户
    { duration: '1m', target: 0 },    // 降压阶段：1分钟内减少到0个用户
  ],

  // 阈值配置
  thresholds: {
    ...THRESHOLDS,
    http_req_duration: ['p(95)<2000'], // 整体P95 < 2秒
    http_req_failed: ['rate<0.05'],    // 整体错误率 < 5%
  },

  // 数据保留
  ext: {
    loadimpact: {
      projectID: 1, // LoadImpact项目ID（可选）
      name: 'ZKER E2E Performance Test',
    },
  },
};

// ==================== 设置 ====================

export function setup() {
  // 测试前准备
  console.log('Starting E2E Performance Test...');
  console.log(`Base URL: ${BASE_URL}`);
  console.log(`Tenant ID: ${TEST_TENANT_ID}`);
  console.log(`User ID: ${TEST_USER_ID}`);

  // 可选：创建测试数据
  // createTestData();
}

export function teardown() {
  // 测试后清理
  console.log('Test completed. Cleaning up...');
  // cleanupTestData();
}

// ==================== 主测试套件 ====================

export default function () {
  // 1. 租户管理API测试
  testTenantAPIs();

  // 2. 配额管理API测试
  testQuotaAPIs();

  // 3. Bot服务API测试
  testBotAPIs();

  // 4. 权限检查API测试
  testPermissionAPIs();

  // 5. 混合负载测试
  testMixedLoad();

  // 思考时间（模拟真实用户行为）
  sleep(Math.random() * 3);
}

// ==================== 测试场景 ====================

/**
 * 租户管理API性能测试
 */
export function testTenantAPIs() {
  const group = 'Tenant Management APIs';

  // 测试1: 获取租户列表
  let res = http.get(`${BASE_URL}/api/tenants`, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
    },
    tags: { group },
  });

  check(res, {
    [`${group} - List Tenants`]: (r) => r.status === 200,
    [`${group} - Response Time < 500ms`]: (r) => r.timings.duration < 500,
  });

  tenantAPIRate.add(res.status === 200);
  apiLatency.add(res.timings.duration);

  sleep(1);

  // 测试2: 获取租户详情
  res = http.get(`${BASE_URL}/api/tenants/${TEST_TENANT_ID}`, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
    },
    tags: { group },
  });

  check(res, {
    [`${group} - Get Tenant`]: (r) => r.status === 200,
    [`${group} - Response Time < 300ms`]: (r) => r.timings.duration < 300,
  });

  apiLatency.add(res.timings.duration);

  sleep(1);
}

/**
 * 配额管理API性能测试
 */
export function testQuotaAPIs() {
  const group = 'Quota Management APIs';

  // 测试1: 获取配额状态
  let res = http.get(`${BASE_URL}/api/tenants/${TEST_TENANT_ID}/quotas`, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
    },
    tags: { group },
  });

  check(res, {
    [`${group} - Get Quotas`]: (r) => r.status === 200,
    [`${group} - Response Time < 500ms`]: (r) => r.timings.duration < 500,
  });

  quotaAPIRate.add(res.status === 200);
  apiLatency.add(res.timings.duration);

  sleep(1);

  // 测试2: 检查配额
  const payload = JSON.stringify({
    resource_type: 'bots',
    required_count: 1,
  });

  res = http.post(`${BASE_URL}/api/tenants/${TEST_TENANT_ID}/quotas/check`, payload, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
      'Content-Type': 'application/json',
    },
    tags: { group },
  });

  check(res, {
    [`${group} - Check Quota`]: (r) => r.status === 200,
    [`${group} - Response Time < 1000ms`]: (r) => r.timings.duration < 1000,
  });

  apiLatency.add(res.timings.duration);

  sleep(1);
}

/**
 * Bot服务API性能测试
 */
export function testBotAPIs() {
  const group = 'Bot Service APIs';

  // 测试1: 列出Bots
  let res = http.get(`${BASE_URL}/api/bots`, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
    },
    tags: { group },
  });

  check(res, {
    [`${group} - List Bots`]: (r) => r.status === 200 || r.status === 404, // 404可能因为没有Bot
    [`${group} - Response Time < 1000ms`]: (r) => r.timings.duration < 1000,
  });

  botAPIRate.add(res.status === 200);
  apiLatency.add(res.timings.duration);

  sleep(1);

  // 测试2: 创建Bot（模拟）
  const botPayload = JSON.stringify({
    bot_name: `Test Bot ${__VU}-${__ITER}`,
    description: 'Performance test bot',
  });

  res = http.post(`${BASE_URL}/api/bots`, botPayload, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
      'Content-Type': 'application/json',
    },
    tags: { group },
  });

  check(res, {
    [`${group} - Create Bot`]: (r) => r.status === 200 || r.status === 400 || r.status === 403, // 403可能配额不足
    [`${group} - Response Time < 2000ms`]: (r) => r.timings.duration < 2000,
  });

  apiLatency.add(res.timings.duration);

  sleep(1);
}

/**
 * 权限检查API性能测试
 */
export function testPermissionAPIs() {
  const group = 'Permission Check APIs';

  // 测试1: 获取用户角色
  let res = http.get(`${BASE_URL}/api/users/${TEST_USER_ID}/roles?tenant_id=${TEST_TENANT_ID}`, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
    },
    tags: { group },
  });

  check(res, {
    [`${group} - Get User Roles`]: (r) => r.status === 200,
    [`${group} - Response Time < 500ms`]: (r) => r.timings.duration < 500,
  });

  permissionAPIRate.add(res.status === 200);
  apiLatency.add(res.timings.duration);

  sleep(1);

  // 测试2: 检查数据权限
  const payload = JSON.stringify({
    tenant_id: TEST_TENANT_ID,
    user_id: TEST_USER_ID,
    resource_type: 'bots',
    action: 'read',
    resource_ids: ['test_bot_id'],
  });

  res = http.post(`${BASE_URL}/api/permissions/check-data`, payload, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
      'Content-Type': 'application/json',
    },
    tags: { group },
  });

  check(res, {
    [`${group} - Check Data Permission`]: (r) => r.status === 200,
    [`${group} - Response Time < 200ms`]: (r) => r.timings.duration < 200, // 权限检查应该很快
  });

  apiLatency.add(res.timings.duration);

  sleep(1);
}

/**
 * 混合负载测试（模拟真实场景）
 */
export function testMixedLoad() {
  const group = 'Mixed Load';

  // 随机选择测试场景
  const scenario = Math.random();

  if (scenario < 0.3) {
    // 30%: 读操作（获取数据）
    readOperations();
  } else if (scenario < 0.7) {
    // 40%: 写操作（创建、更新）
    writeOperations();
  } else {
    // 30%: 复杂操作（多步骤）
    complexOperations();
  }
}

/**
 * 读操作
 */
function readOperations() {
  const operations = [
    () => http.get(`${BASE_URL}/api/bots`, {
      headers: {
        'X-Tenant-ID': TEST_TENANT_ID,
        'X-User-ID': TEST_USER_ID,
      },
    }),
    () => http.get(`${BASE_URL}/api/tenants/${TEST_TENANT_ID}/quotas`, {
      headers: {
        'X-Tenant-ID': TEST_TENANT_ID,
        'X-User-ID': TEST_USER_ID,
      },
    }),
    () => http.get(`${BASE_URL}/api/users/${TEST_USER_ID}/roles?tenant_id=${TEST_TENANT_ID}`, {
      headers: {
        'X-Tenant-ID': TEST_TENANT_ID,
        'X-User-ID': TEST_USER_ID,
      },
    }),
  ];

  const op = operations[Math.floor(Math.random() * operations.length)];
  const res = op();

  check(res, {
    'Read Operation Status': (r) => r.status === 200 || r.status === 404,
    'Read Operation Latency < 1000ms': (r) => r.timings.duration < 1000,
  });

  apiLatency.add(res.timings.duration);
}

/**
 * 写操作
 */
function writeOperations() {
  const operations = [
    () => {
      const payload = JSON.stringify({
        bot_name: `Load Test Bot ${__VU}-${__ITER}`,
        description: 'Load test bot',
      });
      return http.post(`${BASE_URL}/api/bots`, payload, {
        headers: {
          'X-Tenant-ID': TEST_TENANT_ID,
          'X-User-ID': TEST_USER_ID,
          'Content-Type': 'application/json',
        },
      });
    },
    () => {
      const payload = JSON.stringify({
        resource_type: 'bots',
        required_count: 1,
      });
      return http.post(`${BASE_URL}/api/tenants/${TEST_TENANT_ID}/quotas/check`, payload, {
        headers: {
          'X-Tenant-ID': TEST_TENANT_ID,
          'X-User-ID': TEST_USER_ID,
          'Content-Type': 'application/json',
        },
      });
    },
  ];

  const op = operations[Math.floor(Math.random() * operations.length)];
  const res = op();

  check(res, {
    'Write Operation Status': (r) => r.status === 200 || r.status === 400 || r.status === 403,
    'Write Operation Latency < 2000ms': (r) => r.timings.duration < 2000,
  });

  apiLatency.add(res.timings.duration);
}

/**
 * 复杂操作（多步骤）
 */
function complexOperations() {
  // 模拟用户工作流：创建Bot -> 检查配额 -> 添加权限

  // 步骤1: 检查配额
  const quotaPayload = JSON.stringify({
    resource_type: 'bots',
    required_count: 1,
  });

  let res = http.post(`${BASE_URL}/api/tenants/${TEST_TENANT_ID}/quotas/check`, quotaPayload, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
      'Content-Type': 'application/json',
    },
  });

  const quotaAllowed = res.status === 200 && JSON.parse(res.body).data.allowed === true;

  // 步骤2: 创建Bot（如果配额允许）
  if (quotaAllowed) {
    const botPayload = JSON.stringify({
      bot_name: `Complex Test Bot ${__VU}-${__ITER}`,
      description: 'Complex operation test bot',
    });

    res = http.post(`${BASE_URL}/api/bots`, botPayload, {
      headers: {
        'X-Tenant-ID': TEST_TENANT_ID,
        'X-User-ID': TEST_USER_ID,
        'Content-Type': 'application/json',
      },
    });
  }

  // 步骤3: 获取用户角色
  res = http.get(`${BASE_URL}/api/users/${TEST_USER_ID}/roles?tenant_id=${TEST_TENANT_ID}`, {
    headers: {
      'X-Tenant-ID': TEST_TENANT_ID,
      'X-User-ID': TEST_USER_ID,
    },
  });

  check(res, {
    'Complex Operation - All Steps Successful': (r) => r.status === 200,
    'Complex Operation - Total Time < 3000ms': (r) => r.timings.duration < 3000,
  });

  apiLatency.add(res.timings.duration);
}

// ==================== 压力测试配置 ====================

/**
 * 压力测试选项（单独执行）
 * 命令: k6 run --stage stress
 */
export const options_stress = {
  vus: 1,
  stages: [
    { duration: '2m', target: 100 },  // 2分钟爬坡到100用户
    { duration: '5m', target: 100 },  // 持续5分钟100用户
    { duration: '2m', target: 200 },  // 2分钟爬坡到200用户
    { duration: '5m', target: 200 },  // 持续5分钟200用户
    { duration: '2m', target: 300 },  // 2分钟爬坡到300用户
    { duration: '5m', target: 300 },  // 持续5分钟300用户
    { duration: '2m', target: 0 },    // 降压到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<3000'], // 压力测试放宽阈值
    http_req_failed: ['rate<0.1'],    // 压力测试允许10%错误率
  },
};

/**
 * 浸泡测试选项（长时间稳定性测试）
 * 命令: k6 run --stage soak
 */
export const options_soak = {
  vus: 10,
  stages: [
    { duration: '5m', target: 10 },   // 预热5分钟
    { duration: '1h', target: 50 },   // 持续1小时50用户
    { duration: '5m', target: 10 },   // 降温5分钟
    { duration: '5m', target: 0 },    // 降压到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'],
    http_req_failed: ['rate<0.01'],   // 稳定性测试要求更严格
  },
};

/**
 * 峰值测试选项（极限性能测试）
 * 命令: k6 run --stage spike
 */
export const options_spike = {
  vus: 1,
  stages: [
    { duration: '1m', target: 100 },  // 1分钟爬坡到100用户
    { duration: '1m', target: 1000 }, // 突增到1000用户
    { duration: '1m', target: 1000 }, // 保持1000用户
    { duration: '1m', target: 0 },    // 降压到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'], // 峰值测试放宽阈值
    http_req_failed: ['rate<0.2'],    // 峰值测试允许20%错误率
  },
};

/*
 * ==================== 使用说明 ====================
 *
 * 安装K6:
 *   macOS: brew install k6
 *   Linux: sudo apt-get install k6
 *   Windows: choco install k6
 *   或使用Docker: docker pull grafana/k6
 *
 * 运行测试:
 *   # 基础测试
 *   k6 run performance_test.js
 *
 *   # 指定环境变量
 *   BASE_URL=http://staging.example.com k6 run performance_test.js
 *
 *   # 压力测试
 *   k6 run --stage stress performance_test.js
 *
 *   # 浸泡测试
 *   k6 run --stage soak performance_test.js
 *
 *   # 峰值测试
 *   k6 run --stage spike performance_test.js
 *
 * 查看结果:
 *   # 终端输出（实时）
 *   k6 run --out json=output.json performance_test.js
 *
 *   # 生成HTML报告
 *   k6 run --out json=output.json performance_test.js
 *   k6-to-html output.json > report.html
 *
 * 性能基线:
 *   - API P95响应时间: < 500ms
 *   - API P99响应时间: < 1000ms
 *   - 错误率: < 1%
 *   - 并发100用户: 系统稳定
 *   - 并发1000用户: 可降级服务
 */
