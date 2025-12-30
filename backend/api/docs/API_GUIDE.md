# 组织中心管理API使用指南

**版本**: v1.0.0
**最后更新**: 2025-01-01
**维护团队**: Coze Studio

---

## 📖 目录

- [API概览](#api概览)
- [认证授权](#认证授权)
- [通用请求格式](#通用请求格式)
- [通用响应格式](#通用响应格式)
- [错误码说明](#错误码说明)
- [API端点详细说明](#api端点详细说明)
- [请求示例](#请求示例)
- [最佳实践](#最佳实践)
- [常见问题FAQ](#常见问题faq)

---

## API概览

### 基础信息

| 项目 | 说明 |
|------|------|
| **API名称** | 组织中心管理API |
| **版本** | v1.0.0 |
| **协议** | HTTPS/HTTP |
| **数据格式** | JSON |
| **字符编码** | UTF-8 |
| **架构风格** | RESTful |

### 服务端点

| 环境 | URL | 说明 |
|------|-----|------|
| 开发环境 | http://localhost:8080 | 本地开发 |
| 测试环境 | https://api-test.coze.com | 测试验证 |
| 生产环境 | https://api.coze.com | 生产使用 |

### 核心功能

组织中心管理API提供以下核心功能：

1. **组织管理** - 管理公司、分公司、部门、项目组的树形结构
2. **部门管理** - 部门层级关系维护和查询
3. **员工管理** - 员工全生命周期管理（入职、转正、调岗、离职）
4. **岗位管理** - 岗位定义、职级体系、职责要求
5. **通讯录服务** - 组织架构目录查询、员工搜索

---

## 认证授权

### 认证方式

API支持两种认证方式：

#### 1. Bearer Token认证（推荐）

所有API请求都需要在HTTP Header中携带JWT Token：

```http
Authorization: Bearer <your_jwt_token>
```

**Token获取方式**：

```bash
# 登录获取Token
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your_password"
}
```

**响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200
  }
}
```

#### 2. API Key认证

```http
X-API-Key: <your_api_key>
```

### 租户隔离

所有API请求都需要携带租户标识：

**方式一：通过Header传递**

```http
X-Tenant-ID: tenant_001
```

**方式二：通过Query参数传递**

```
GET /api/organizations?tenant_id=tenant_001
```

**方式三：从Token解析**（推荐）

JWT Token中包含租户ID，后端自动解析并进行租户隔离。

### 权限控制

API采用RBAC（基于角色的访问控制）模型：

| 角色 | 权限范围 |
|------|---------|
| **系统管理员** | 所有操作权限 |
| **租户管理员** | 租户内所有操作权限 |
| **部门经理** | 本部门及子部门的员工查看和管理权限 |
| **普通员工** | 仅查看权限（通讯录、员工信息） |

---

## 通用请求格式

### HTTP方法

| 方法 | 说明 | 示例 |
|------|------|------|
| GET | 查询资源 | GET /api/organizations/{id} |
| POST | 创建资源 | POST /api/organizations |
| PUT | 更新资源 | PUT /api/organizations/{id} |
| DELETE | 删除资源 | DELETE /api/organizations/{id} |

### 请求Header

```http
Content-Type: application/json
Authorization: Bearer <token>
X-Tenant-ID: tenant_001
X-Request-ID: uuid-string
```

### 请求体示例

```json
{
  "tenant_id": "tenant_001",
  "org_name": "技术部",
  "org_type": "department",
  "org_code": "TECH",
  "description": "负责产品研发",
  "sort_order": 1
}
```

### 分页参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量（最大100） |
| sort_by | string | 否 | created_at | 排序字段 |
| sort_order | string | 否 | desc | 排序方向（asc/desc） |

---

## 通用响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "message_en": "success",
  "data": {
    // 业务数据
  }
}
```

### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": [
    // 数据列表
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### 错误响应

```json
{
  "code": 40001,
  "message": "参数错误",
  "message_en": "Invalid parameters",
  "data": {
    "description": "tenant_id is required",
    "field": "tenant_id"
  }
}
```

---

## 错误码说明

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（未登录或token无效） |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如唯一索引冲突） |
| 500 | 服务器内部错误 |

### 业务错误码

| 错误码 | HTTP状态码 | 说明（中文） | 说明（英文） |
|--------|-----------|-------------|-------------|
| 0 | 200 | 成功 | Success |
| 40001 | 400 | 参数错误 | Invalid parameters |
| 40002 | 400 | 参数缺失 | Missing required parameter |
| 40003 | 400 | 参数格式错误 | Invalid parameter format |
| 40101 | 401 | 未授权访问 | Unauthorized access |
| 40102 | 401 | Token无效 | Invalid token |
| 40103 | 401 | Token已过期 | Token expired |
| 40301 | 403 | 权限不足 | Insufficient permissions |
| 40302 | 403 | 跨租户访问拒绝 | Cross-tenant access denied |
| 40401 | 404 | 资源不存在 | Resource not found |
| 40402 | 404 | 组织不存在 | Organization not found |
| 40403 | 404 | 部门不存在 | Department not found |
| 40404 | 404 | 员工不存在 | Employee not found |
| 40405 | 404 | 岗位不存在 | Position not found |
| 40901 | 409 | 资源已存在 | Resource already exists |
| 40902 | 409 | 组织编码重复 | Organization code duplicate |
| 40903 | 409 | 员工工号重复 | Employee code duplicate |
| 40904 | 409 | 邮箱已被使用 | Email already in use |
| 50001 | 500 | 服务器内部错误 | Internal server error |
| 50002 | 500 | 数据库错误 | Database error |
| 50003 | 500 | 外部服务错误 | External service error |

---

## API端点详细说明

### 组织管理 (Organizations)

#### 1. 创建组织

**接口**: `POST /api/organizations`

**请求参数**:

```json
{
  "tenant_id": "string (必填)",
  "org_name": "string (必填, 1-200字符)",
  "org_type": "string (必填) - 枚举: company|division|department|project",
  "parent_id": "string (可选)",
  "org_code": "string (必填, 1-50字符)",
  "leader_id": "string (可选)",
  "description": "string (可选)",
  "sort_order": "integer (可选)"
}
```

**业务规则**:
- 组织编码在租户内唯一
- 根组织（公司）不能有父组织
- 自动计算组织层级和路径

**示例**:

```bash
curl -X POST https://api.coze.com/api/organizations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant_001",
    "org_name": "研发中心",
    "org_type": "division",
    "parent_id": "org_company_001",
    "org_code": "R&D",
    "leader_id": "emp_001",
    "description": "负责产品研发和技术创新"
  }'
```

#### 2. 获取组织详情

**接口**: `GET /api/organizations/{id}`

**路径参数**:
- `id`: 组织ID

**示例**:

```bash
curl -X GET https://api.coze.com/api/organizations/org_001 \
  -H "Authorization: Bearer <token>"
```

#### 3. 更新组织信息

**接口**: `PUT /api/organizations/{id}`

**请求参数**:

```json
{
  "org_name": "string (可选)",
  "leader_id": "string (可选)",
  "description": "string (可选)",
  "sort_order": "integer (可选)",
  "status": "string (可选) - 枚举: active|inactive|frozen"
}
```

#### 4. 删除组织

**接口**: `DELETE /api/organizations/{id}`

**业务规则**:
- 有子组织的组织不能删除
- 使用软删除机制

#### 5. 获取组织树

**接口**: `GET /api/organizations/tree`

**查询参数**:
- `tenant_id` (必填): 租户ID

#### 6. 分页查询组织列表

**接口**: `GET /api/organizations`

**查询参数**:
- `tenant_id` (必填): 租户ID
- `org_type` (可选): 组织类型
- `status` (可选): 状态
- `parent_id` (可选): 父组织ID
- `keyword` (可选): 搜索关键词
- `page` (可选): 页码，默认1
- `page_size` (可选): 每页数量，默认20
- `sort_by` (可选): 排序字段
- `sort_order` (可选): 排序方向（asc/desc）

#### 7. 移动组织

**接口**: `POST /api/organizations/{id}/move`

**请求参数**:

```json
{
  "org_id": "string (必填)",
  "new_parent_id": "string (可选)"
}
```

**业务规则**:
- 不能移动到自己的后代组织下
- 移动后自动更新层级和路径

#### 8. 获取子组织/祖先/后代

**接口**:
- `GET /api/organizations/{id}/children` - 获取直接子组织
- `GET /api/organizations/{id}/ancestors` - 获取所有祖先
- `GET /api/organizations/{id}/descendants` - 获取所有后代

---

### 部门管理 (Departments)

#### 1. 创建部门

**接口**: `POST /api/org/departments`

**请求参数**:

```json
{
  "tenant_id": "string (必填)",
  "org_id": "string (必填)",
  "parent_id": "string (可选)",
  "dept_name": "string (必填, 1-200字符)",
  "dept_code": "string (必填, 1-50字符)",
  "leader_id": "string (可选)",
  "parent_leader": "string (可选)",
  "description": "string (可选)",
  "sort_order": "integer (可选)"
}
```

#### 2-8. 其他部门接口

与组织管理类似，包括：
- 获取部门详情: `GET /api/org/departments/{id}`
- 更新部门: `PUT /api/org/departments/{id}`
- 删除部门: `DELETE /api/org/departments/{id}`
- 获取部门树: `GET /api/org/departments/tree`
- 分页查询: `GET /api/org/departments`
- 移动部门: `POST /api/org/departments/{id}/move`
- 获取子部门/祖先/后代

---

### 员工管理 (Employees)

#### 1. 创建员工

**接口**: `POST /api/v1/employees`

**请求参数**:

```json
{
  "tenant_id": "string (必填)",
  "org_id": "string (必填)",
  "dept_id": "string (可选)",
  "position_id": "string (可选)",
  "emp_name": "string (必填)",
  "emp_code": "string (必填)",
  "employee_type": "string (必填) - 枚举: full_time|part_time|intern|outsourcing|contractor",
  "job_level": "integer (必填, 1-10)",
  "job_title": "string (可选)",
  "hire_date": "integer (必填, 毫秒时间戳)",
  "probation_days": "integer (可选)",
  "email": "string (可选)",
  "phone": "string (可选)"
}
```

**业务规则**:
- 工号在租户内唯一
- 邮箱在租户内唯一（如果提供）
- 默认状态为试用期(trial)

**示例**:

```bash
curl -X POST https://api.coze.com/api/v1/employees \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant_001",
    "org_id": "org_001",
    "dept_id": "dept_001",
    "position_id": "pos_001",
    "emp_name": "张三",
    "emp_code": "E001",
    "employee_type": "full_time",
    "job_level": 5,
    "job_title": "高级工程师",
    "hire_date": 1704067200000,
    "probation_days": 90,
    "email": "zhangsan@example.com",
    "phone": "13800138000"
  }'
```

#### 2. 获取员工详情

**接口**: `GET /api/v1/employees/{id}`

#### 3. 更新员工信息

**接口**: `PUT /api/v1/employees/{id}`

**业务规则**:
- 不能修改tenant_id、emp_id、org_id等关键字段
- 调换部门时需验证新部门属于同一组织

#### 4. 删除员工（软删除）

**接口**: `DELETE /api/v1/employees/{id}`

**业务规则**:
- 在职员工不能删除，必须先办理离职
- 使用软删除机制

#### 5. 更新员工状态

**接口**: `PUT /api/v1/employees/{id}/status`

**请求参数**:

```json
{
  "status": "string (必填) - 枚举: active|trial|probation|resigned|suspended|retired"
}
```

**状态转换规则**:
- trial -> active: 转正
- active -> resigned: 离职
- active -> suspended: 停职
- suspended -> active: 复职

#### 6. 分页查询员工列表

**接口**: `GET /api/v1/employees`

**查询参数**:
- `org_id` (可选): 组织ID
- `dept_id` (可选): 部门ID
- `position_id` (可选): 岗位ID
- `status` (可选): 员工状态
- `employee_type` (可选): 员工类型
- `keyword` (可选): 搜索关键词
- `page_size` (可选): 每页数量，默认20
- `page_token` (可选): 分页令牌

#### 7-10. 其他员工查询接口

- 根据工号查询: `GET /api/v1/employees/by-code/{code}`
- 根据用户ID查询: `GET /api/v1/employees/by-user/{user_id}`
- 获取部门员工: `GET /api/v1/employees/by-dept/{dept_id}`
- 获取组织员工: `GET /api/v1/employees/by-org/{org_id}`
- 搜索员工: `GET /api/v1/employees/search`
- 按拼音查询: `GET /api/v1/employees/pinyin/{pinyin}`

---

### 岗位管理 (Positions)

#### 1. 创建岗位

**接口**: `POST /api/positions`

**请求参数**:

```json
{
  "tenant_id": "string (必填)",
  "dept_id": "string (可选)",
  "position_name": "string (必填, 1-100字符)",
  "position_code": "string (必填, 1-50字符)",
  "level": "integer (必填, 1-10)",
  "category": "string (必填, 1-50字符)",
  "responsibilities": "string (可选)",
  "requirements": "string (可选)",
  "sort_order": "integer (可选)"
}
```

#### 2-7. 其他岗位接口

- 获取岗位详情: `GET /api/positions/{id}`
- 更新岗位: `PUT /api/positions/{id}`
- 删除岗位: `DELETE /api/positions/{id}`
- 分页查询: `GET /api/positions`
- 根据编码查询: `GET /api/positions/by-code/{code}`
- 获取部门岗位: `GET /api/positions/by-dept/{dept_id}`
- 按职级查询: `GET /api/positions/by-level/{level}`
- 按类别查询: `GET /api/positions/by-category/{category}`

---

### 通讯录服务 (Directory)

#### 1. 获取组织架构目录

**接口**: `GET /api/directory/organization`

#### 2. 获取部门目录

**接口**: `GET /api/directory/department`

#### 3. 搜索员工

**接口**: `GET /api/directory/search`

**查询参数**:
- `keyword` (必填): 搜索关键词
- `limit` (可选): 返回数量限制，默认20，最大100
- `tenant_id` (可选): 租户ID

#### 4. 获取部门/组织员工列表

- `GET /api/directory/departments/{dept_id}/employees`
- `GET /api/directory/organizations/{org_id}/employees`

#### 5. 根据工号查询员工

**接口**: `GET /api/directory/employee/by-code/{code}`

---

## 请求示例

### cURL示例

#### 创建组织

```bash
curl -X POST https://api.coze.com/api/organizations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant_001" \
  -d '{
    "tenant_id": "tenant_001",
    "org_name": "技术部",
    "org_type": "department",
    "org_code": "TECH",
    "description": "负责产品研发"
  }'
```

#### 查询员工列表

```bash
curl -X GET "https://api.coze.com/api/v1/employees?page=1&page_size=20&dept_id=dept_001" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### JavaScript示例

#### 使用Fetch API

```javascript
// 创建组织
async function createOrganization() {
  const response = await fetch('https://api.coze.com/api/organizations', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + token,
      'Content-Type': 'application/json',
      'X-Tenant-ID': 'tenant_001'
    },
    body: JSON.stringify({
      tenant_id: 'tenant_001',
      org_name: '技术部',
      org_type: 'department',
      org_code: 'TECH',
      description: '负责产品研发'
    })
  });

  const data = await response.json();
  console.log(data);
}
```

#### 使用Axios

```javascript
// 查询员工列表
async function getEmployees() {
  try {
    const response = await axios.get('https://api.coze.com/api/v1/employees', {
      headers: {
        'Authorization': 'Bearer ' + token
      },
      params: {
        page: 1,
        page_size: 20,
        dept_id: 'dept_001'
      }
    });

    console.log(response.data);
  } catch (error) {
    console.error('Error:', error.response.data);
  }
}
```

### Python示例

#### 使用Requests

```python
import requests

# 创建组织
def create_organization():
    url = 'https://api.coze.com/api/organizations'
    headers = {
        'Authorization': 'Bearer ' + token,
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'tenant_001'
    }
    data = {
        'tenant_id': 'tenant_001',
        'org_name': '技术部',
        'org_type': 'department',
        'org_code': 'TECH',
        'description': '负责产品研发'
    }

    response = requests.post(url, headers=headers, json=data)
    return response.json()

# 查询员工列表
def get_employees():
    url = 'https://api.coze.com/api/v1/employees'
    headers = {
        'Authorization': 'Bearer ' + token
    }
    params = {
        'page': 1,
        'page_size': 20,
        'dept_id': 'dept_001'
    }

    response = requests.get(url, headers=headers, params=params)
    return response.json()
```

---

## 最佳实践

### 1. 错误处理

始终检查HTTP状态码和业务错误码：

```javascript
async function safeApiCall() {
  try {
    const response = await fetch(url, options);
    const data = await response.json();

    if (data.code !== 0) {
      console.error(`API Error: ${data.message} (${data.code})`);
      // 处理业务错误
      return;
    }

    // 处理成功响应
    return data.data;
  } catch (error) {
    console.error('Network Error:', error);
    // 处理网络错误
  }
}
```

### 2. Token管理

- Token有效期2小时，建议在过期前刷新
- 使用localStorage或sessionStorage存储Token
- 每次请求前检查Token有效性

```javascript
function getToken() {
  const token = localStorage.getItem('api_token');
  const expiresAt = localStorage.getItem('token_expires_at');

  if (!token || Date.now() > expiresAt) {
    // Token过期，重新登录
    return refreshToken();
  }

  return token;
}
```

### 3. 租户隔离

始终确保请求携带租户标识：

```javascript
const headers = {
  'Authorization': 'Bearer ' + token,
  'X-Tenant-ID': getCurrentTenantId() // 从用户会话获取
};
```

### 4. 分页查询

使用游标分页获取大量数据：

```javascript
async function fetchAllEmployees() {
  let pageToken = null;
  let allEmployees = [];

  do {
    const response = await fetch(
      `/api/v1/employees?page_size=100&page_token=${pageToken || ''}`
    );
    const data = await response.json();

    allEmployees = allEmployees.concat(data.data.employees);
    pageToken = data.data.next_page_token;
  } while (pageToken);

  return allEmployees;
}
```

### 5. 请求去重

对于幂等的GET请求，使用请求ID去重：

```javascript
const requestId = crypto.randomUUID();

const headers = {
  'X-Request-ID': requestId
};
```

### 6. 重试机制

对于网络错误，实现指数退避重试：

```javascript
async function fetchWithRetry(url, options, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(url, options);
      if (response.ok) {
        return await response.json();
      }
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await new Promise(resolve => setTimeout(resolve, Math.pow(2, i) * 1000));
    }
  }
}
```

---

## 常见问题FAQ

### Q1: 如何获取API Token？

**答**: 通过登录接口获取：

```bash
POST /api/auth/login
{
  "username": "your_username",
  "password": "your_password"
}
```

返回的`token`字段即为JWT Token。

### Q2: Token过期了怎么办？

**答**: Token有效期为2小时，可以使用Refresh Token刷新，或重新登录获取新Token。

### Q3: 如何处理跨租户访问？

**答**: API自动进行租户隔离，请求中的`tenant_id`必须与Token中的租户ID一致，否则返回403错误。

### Q4: 分页查询的最大page_size是多少？

**答**: 最大为100。建议使用20-50的页面大小以平衡性能和用户体验。

### Q5: 如何查询所有员工（包括离职的）？

**答**: 不指定`status`参数即可查询所有状态的员工。

### Q6: 软删除的数据如何恢复？

**答**: 软删除的数据仅设置`deleted_at`字段，可通过管理员手动恢复或调用专门的数据恢复接口。

### Q7: API有请求频率限制吗？

**答**: 是的，每个租户每分钟最多1000次请求。超出限制将返回429错误。

### Q8: 如何批量导入员工数据？

**答**: 使用批量导入接口 `POST /api/v1/employees/batch`，支持CSV和Excel格式。

### Q9: 组织和部门有什么区别？

**答**:
- **组织**: 用于管理公司、分公司、项目组等大型业务单元
- **部门**: 用于细分组织内部的职能单元

两者都是树形结构，可以独立管理或关联使用。

### Q10: 如何获取完整的组织架构树？

**答**: 调用 `GET /api/organizations/tree` 或 `GET /api/directory/organization` 接口。

---

## 附录

### A. 数据模型

详细的数据模型请参考 [Swagger文档](./swagger/swagger.yaml)。

### B. 错误码完整列表

完整错误码列表请参考 [统一错误码规范](../../企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)。

### C. SDK和客户端库

- [Go SDK](./CLIENT_SDK.md#go-sdk)
- [JavaScript SDK](./CLIENT_SDK.md#javascript-sdk)
- [Python SDK](./CLIENT_SDK.md#python-sdk)
- [Java SDK](./CLIENT_SDK.md#java-sdk)

### D. 相关文档

- [数据库设计文档](../../企业级功能完善与统一性设计方案/数据库设计完整交付清单.md)
- [API设计规范](../../企业级功能完善与统一性设计方案/API设计规范文档.md)

---

## 更新日志

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0.0 | 2025-01-01 | 初始版本发布 |

---

## 技术支持

如有问题，请联系：

- 邮箱: api-support@coze.com
- 文档: https://docs.coze.com
- GitHub Issues: https://github.com/coze-dev/coze-studio/issues
