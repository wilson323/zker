# API接口文档：组织中心模块（增强版）

**模块名称**: 组织中心增强版 (Organization Center Plus)
**设计文档**: 06-企业管理中台_组织中心_增强版.md
**版本**: v2.0.0
**最后更新**: 2025-01-03
**优先级**: P0

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 企业通讯录API](#2-企业通讯录api)
- [3. 组织架构管理API](#3-组织架构管理api)
- [4. 员工全生命周期API](#4-员工全生命周期api)
- [5. 员工档案管理API](#5-员工档案管理api)
- [6. 虚拟组织管理API](#6-虚拟组织管理api)
- [7. 组织代理管理API](#7-组织代理管理api)
- [8. 数据模型](#8-数据模型)
- [9. 后端代码示例](#9-后端代码示例)
- [10. 前端代码示例](#10-前端代码示例)
- [11. 错误码定义](#11-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

组织中心增强版提供企业级的组织与人员管理能力，包括：

- ✅ **企业通讯录** - 按部门/按拼音查看同事信息、快速搜索、组织关系可视化
- ✅ **组织架构管理** - 部门树管理、组织关系、汇报关系
- ✅ **员工全生命周期** - 入职、试用期、调岗、晋升、离职、黑名单管理
- ✅ **员工档案** - 基本信息扩展、工作经历、教育背景、技能标签、绩效记录
- ✅ **虚拟组织** - 项目组、委员会、虚拟团队
- ✅ **组织代理** - 直属代理、授权委托、委托记录

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0
- 搜索: Elasticsearch（用于通讯录搜索）

**前端技术栈**:
- 框架: React 18 + TypeScript
- UI库: Semi Design
- 状态管理: Zustand

---

## 2. 企业通讯录API

### 2.1 获取通讯录（按部门视图）

**接口地址**: `GET /api/v1/tenants/{tenant_id}/contacts/by-department`

**功能说明**: 获取按部门组织的通讯录，支持树形视图

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| include_members | boolean | 否 | 是否包含成员列表，默认true |
| tree_view | boolean | 否 | 是否返回树形结构，默认true |
| include_inactive | boolean | 否 | 是否包含离职员工，默认false |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "organizations": [
      {
        "id": "1",
        "name": "研发中心",
        "type": "department",
        "level": 1,
        "parent_id": null,
        "path": "/1",
        "manager": {
          "id": "10",
          "name": "张三",
          "avatar": "https://cdn.example.com/avatar/10.jpg",
          "position": "技术总监"
        },
        "member_count": 150,
        "children": [
          {
            "id": "11",
            "name": "前端开发部",
            "type": "department",
            "level": 2,
            "parent_id": "1",
            "path": "/1/11",
            "manager": {
              "id": "11",
              "name": "李四",
              "avatar": "https://cdn.example.com/avatar/11.jpg",
              "position": "前端经理"
            },
            "member_count": 50,
            "members": [
              {
                "id": "101",
                "name": "王五",
                "avatar": "https://cdn.example.com/avatar/101.jpg",
                "position": "前端工程师",
                "job_level": "P5",
                "phone": "138****1234",
                "email": "wangwu@example.com",
                "entry_status": "regular",
                "office_location": "A栋3楼"
              }
            ],
            "children": []
          },
          {
            "id": "12",
            "name": "后端开发部",
            "type": "department",
            "level": 2,
            "parent_id": "1",
            "path": "/1/12",
            "manager": {
              "id": "12",
              "name": "赵六",
              "avatar": "https://cdn.example.com/avatar/12.jpg",
              "position": "后端经理"
            },
            "member_count": 60,
            "members": [],
            "children": []
          }
        ]
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.2 获取通讯录（按拼音视图）

**接口地址**: `GET /api/v1/tenants/{tenant_id}/contacts/by-pinyin`

**功能说明**: 获取按拼音首字母分组的通讯录

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| initial | string | 否 | 拼音首字母过滤（A-Z） |
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认50 |
| include_inactive | boolean | 否 | 是否包含离职员工，默认false |

**请求示例**:
```http
GET /api/v1/tenants/tenant-123/contacts/by-pinyin?initial=A&page=1&page_size=50
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 500,
    "A": [
      {
        "id": "101",
        "name": "艾伦",
        "pinyin": "ai lun",
        "pinyin_initial": "A",
        "avatar": "https://cdn.example.com/avatar/101.jpg",
        "organization": {
          "id": "11",
          "name": "前端开发部"
        },
        "position": "产品经理",
        "job_level": "P6",
        "phone": "138****1234",
        "email": "ailun@example.com",
        "entry_status": "regular",
        "office_location": "A栋3楼"
      }
    ],
    "B": [
      {
        "id": "102",
        "name": "白云",
        "pinyin": "bai yun",
        "pinyin_initial": "B",
        "avatar": "https://cdn.example.com/avatar/102.jpg",
        "organization": {
          "id": "12",
          "name": "后端开发部"
        },
        "position": "后端工程师",
        "job_level": "P5",
        "phone": "139****5678",
        "email": "baiyun@example.com",
        "entry_status": "probation",
        "office_location": "B栋2楼"
      }
    ],
    "Z": [...]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.3 获取同事详情

**接口地址**: `GET /api/v1/tenants/{tenant_id}/contacts/{member_id}`

**功能说明**: 获取同事的详细信息（包含组织关系、档案信息等）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "101",
    "name": "张三",
    "avatar": "https://cdn.example.com/avatar/101.jpg",
    "email": "zhangsan@example.com",
    "phone": "138****1234",
    "work_email": "zhangsan@company.com",
    "work_phone": "010-12345678",
    "office_location": "A栋3楼",

    "organization": {
      "id": "11",
      "name": "前端开发部",
      "path": "/1/11",
      "type": "department"
    },

    "position": {
      "id": "20",
      "name": "前端工程师",
      "job_level": "P5",
      "salary_grade": "P5-2"
    },

    "direct_manager": {
      "id": "11",
      "name": "李四",
      "position": "前端经理",
      "email": "lisi@example.com"
    },

    "direct_subordinates": [
      {
        "id": "201",
        "name": "小明",
        "position": "前端实习生"
      }
    ],

    "virtual_organizations": [
      {
        "id": "v1",
        "name": "AI项目组",
        "type": "project",
        "role": "前端负责人"
      }
    ],

    "entry_status": "regular",
    "joined_at": "2024-01-15",

    "work_experiences": [
      {
        "company": "某科技公司",
        "position": "前端工程师",
        "start_date": "2020-01-01",
        "end_date": "2023-12-31",
        "description": "负责公司官网前端开发"
      }
    ],

    "educations": [
      {
        "school": "XX大学",
        "degree": "master",
        "major": "计算机科学",
        "start_date": "2017-09-01",
        "end_date": "2020-06-30",
        "gpa": 3.8
      }
    ],

    "skills": [
      {
        "name": "React",
        "type": "tech",
        "proficiency": "expert",
        "certified": false
      },
      {
        "name": "英语",
        "type": "language",
        "proficiency": "advanced",
        "certified": true
      }
    ],

    "performance_reviews": [
      {
        "review_type": "annual",
        "review_period": "2024",
        "rating": 4,
        "result": "good",
        "review_date": "2024-12-31"
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.4 搜索同事

**接口地址**: `GET /api/v1/tenants/{tenant_id}/contacts/search`

**功能说明**: 全文搜索同事（支持姓名、手机、邮箱、岗位、组织）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| include_inactive | boolean | 否 | 是否包含离职员工，默认false |

**请求示例**:
```http
GET /api/v1/tenants/tenant-123/contacts/search?keyword=张三&page=1&page_size=20
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": "101",
        "name": "张三",
        "avatar": "https://cdn.example.com/avatar/101.jpg",
        "organization": "前端开发部",
        "position": "前端工程师",
        "phone": "138****1234",
        "email": "zhangsan@example.com",
        "highlight": {
          "name": "<em>张三</em>",
          "organization": "<em>前端</em>开发部",
          "position": "<em>前端</em>工程师"
        }
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.5 获取我的下属

**接口地址**: `GET /api/v1/tenants/{tenant_id}/contacts/my-subordinates`

**功能说明**: 获取当前用户的所有下属（直接和间接）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| level | string | 否 | 层级: direct（直接）/ indirect（间接）/ all（全部），默认all |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "direct_count": 5,
    "indirect_count": 10,
    "items": [
      {
        "id": "201",
        "name": "小明",
        "avatar": "https://cdn.example.com/avatar/201.jpg",
        "organization": "前端开发部",
        "position": "前端实习生",
        "reporting_line": "direct",
        "entry_status": "probation"
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.6 获取我的上级

**接口地址**: `GET /api/v1/tenants/{tenant_id}/contacts/my-managers`

**功能说明**: 获取当前用户的所有上级（直接和间接）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| level | string | 否 | 层级: direct（直接）/ indirect（间接）/ all（全部），默认all |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "direct_manager": {
      "id": "11",
      "name": "李四",
      "avatar": "https://cdn.example.com/avatar/11.jpg",
      "position": "前端经理",
      "phone": "139****5678",
      "email": "lisi@example.com"
    },
    "indirect_managers": [
      {
        "id": "10",
        "name": "张三",
        "avatar": "https://cdn.example.com/avatar/10.jpg",
        "position": "技术总监",
        "level": 2
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

## 3. 组织架构管理API

### 3.1 获取组织树

**接口地址**: `GET /api/v1/tenants/{tenant_id}/organizations/tree`

**功能说明**: 获取完整的组织架构树

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| include_members | boolean | 否 | 是否包含成员，默认false |
| include_inactive | boolean | 否 | 是否包含已禁用组织，默认false |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "1",
    "name": "公司",
    "type": "company",
    "level": 1,
    "parent_id": null,
    "path": "/1",
    "manager": null,
    "member_count": 500,
    "children": [
      {
        "id": "2",
        "name": "研发中心",
        "type": "department",
        "level": 2,
        "parent_id": "1",
        "path": "/1/2",
        "manager": {
          "id": "10",
          "name": "技术总监"
        },
        "member_count": 150,
        "children": []
      },
      {
        "id": "3",
        "name": "产品中心",
        "type": "department",
        "level": 2,
        "parent_id": "1",
        "path": "/1/3",
        "manager": {
          "id": "20",
          "name": "产品总监"
        },
        "member_count": 80,
        "children": []
      }
    ]
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.2 创建组织

**接口地址**: `POST /api/v1/tenants/{tenant_id}/organizations`

**功能说明**: 创建新组织（部门/团队）

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "name": "测试部",
  "type": "department",
  "parent_id": "2",
  "manager_id": 30,
  "description": "负责产品测试"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 组织名称 |
| type | string | 是 | 组织类型: company/department/team/group |
| parent_id | int64 | 是 | 父组织ID |
| manager_id | int64 | 否 | 负责人ID |
| description | string | 否 | 描述 |

**响应示例**:
```json
{
  "code": 0,
  "message": "组织创建成功",
  "data": {
    "id": 100,
    "name": "测试部",
    "type": "department",
    "parent_id": "2",
    "path": "/1/2/100",
    "created_at": "2025-01-03T12:00:00Z"
  },
  "timestamp": "2025-01-03T12:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.3 更新组织

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/organizations/{org_id}`

**功能说明**: 更新组织信息

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| org_id | int64 | 组织ID |

**请求参数**:
```json
{
  "name": "测试中心",
  "manager_id": 31,
  "description": "负责产品测试与质量保证"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "组织更新成功",
  "data": {
    "id": 100,
    "updated_at": "2025-01-03T13:00:00Z"
  },
  "timestamp": "2025-01-03T13:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.4 删除组织

**接口地址**: `DELETE /api/v1/tenants/{tenant_id}/organizations/{org_id}`

**功能说明**: 删除组织（需先转移成员）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| org_id | int64 | 组织ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "组织删除成功",
  "data": {
    "id": 100,
    "deleted_at": "2025-01-03T14:00:00Z"
  },
  "timestamp": "2025-01-03T14:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.5 移动组织

**接口地址**: `POST /api/v1/tenants/{tenant_id}/organizations/{org_id}/move`

**功能说明**: 移动组织到新的父组织下

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| org_id | int64 | 组织ID |

**请求参数**:
```json
{
  "new_parent_id": 3,
  "reason": "组织架构调整"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "组织移动成功",
  "data": {
    "id": 100,
    "old_path": "/1/2/100",
    "new_path": "/1/3/100",
    "updated_at": "2025-01-03T15:00:00Z"
  },
  "timestamp": "2025-01-03T15:00:00Z",
  "trace_id": "abc123"
}
```

---

## 4. 员工全生命周期API

### 4.1 试用期管理 - 转正申请

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/probation-review`

**功能说明**: 提交员工转正申请

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "status": "passed",
  "feedback": "表现优秀，同意转正",
  "new_position_id": 20,
  "new_salary_grade": "P5-2",
  "effective_date": "2025-02-01"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 是 | 转正结果: passed（通过）/ failed（未通过）/ extended（延长试用期） |
| feedback | string | 是 | 反馈意见 |
| new_position_id | int64 | 否 | 转正后新岗位ID |
| new_salary_grade | string | 否 | 转正后薪级 |
| effective_date | string | 是 | 生效日期 |

**响应示例**:
```json
{
  "code": 0,
  "message": "转正申请已提交",
  "data": {
    "review_id": "123",
    "status": "pending_approval",
    "created_at": "2025-01-03T16:00:00Z"
  },
  "timestamp": "2025-01-03T16:00:00Z",
  "trace_id": "abc123"
}
```

---

### 4.2 调岗申请

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/transfer`

**功能说明**: 提交调岗申请

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "change_type": "transfer",
  "new_organization_id": 3,
  "new_position_id": 30,
  "new_job_level": "P6",
  "effective_date": "2025-02-01",
  "reason": "业务调整"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| change_type | string | 是 | 变更类型: transfer（调岗）/ promotion（晋升）/ demotion（降职）/ adjustment（调整） |
| new_organization_id | int64 | 是 | 新组织ID |
| new_position_id | int64 | 是 | 新岗位ID |
| new_job_level | string | 否 | 新职级 |
| effective_date | string | 是 | 生效日期 |
| reason | string | 是 | 变更原因 |

**响应示例**:
```json
{
  "code": 0,
  "message": "调岗申请已提交",
  "data": {
    "change_id": "456",
    "status": "pending_approval",
    "created_at": "2025-01-03T17:00:00Z"
  },
  "timestamp": "2025-01-03T17:00:00Z",
  "trace_id": "abc123"
}
```

---

### 4.3 离职申请

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/resign`

**功能说明**: 提交离职申请

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "last_working_day": "2025-02-15",
  "reason": "个人发展",
  "handover_to": 102
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| last_working_day | string | 是 | 最后工作日 |
| reason | string | 是 | 离职原因 |
| handover_to | int64 | 是 | 交接人ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "离职申请已提交",
  "data": {
    "resign_id": "789",
    "status": "pending_approval",
    "created_at": "2025-01-03T18:00:00Z"
  },
  "timestamp": "2025-01-03T18:00:00Z",
  "trace_id": "abc123"
}
```

---

### 4.4 获取调岗记录

**接口地址**: `GET /api/v1/tenants/{tenant_id}/members/{member_id}/position-changes`

**功能说明**: 获取员工的调岗历史记录

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "items": [
      {
        "id": 456,
        "change_type": "promotion",
        "old_organization": {
          "id": 11,
          "name": "前端开发部"
        },
        "old_position": {
          "id": 20,
          "name": "前端工程师"
        },
        "old_job_level": "P5",
        "new_organization": {
          "id": 11,
          "name": "前端开发部"
        },
        "new_position": {
          "id": 21,
          "name": "高级前端工程师"
        },
        "new_job_level": "P6",
        "effective_date": "2024-06-01",
        "reason": "晋升",
        "status": "approved",
        "created_at": "2024-05-15T10:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T19:00:00Z",
  "trace_id": "abc123"
}
```

---

## 5. 员工档案管理API

### 5.1 更新员工基本信息

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/members/{member_id}/profile`

**功能说明**: 更新员工基本信息（档案扩展字段）

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "gender": "male",
  "birth_date": "1995-01-01",
  "marital_status": "single",
  "nationality": "中国",
  "native_place": "浙江杭州",
  "address": "浙江省杭州市西湖区",
  "emergency_contact_name": "张父",
  "emergency_contact_phone": "139****8765",
  "emergency_contact_relation": "父亲"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "档案更新成功",
  "data": {
    "member_id": 101,
    "updated_at": "2025-01-03T20:00:00Z"
  },
  "timestamp": "2025-01-03T20:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.2 添加工作经历

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/work-experiences`

**功能说明**: 为员工添加工作经历

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "company_name": "某科技公司",
  "industry": "互联网",
  "position": "前端工程师",
  "start_date": "2020-01-01",
  "end_date": "2023-12-31",
  "description": "负责公司官网前端开发，使用React和TypeScript"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "工作经历添加成功",
  "data": {
    "experience_id": 501,
    "created_at": "2025-01-03T21:00:00Z"
  },
  "timestamp": "2025-01-03T21:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.3 添加教育背景

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/educations`

**功能说明**: 为员工添加教育背景

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "school_name": "XX大学",
  "degree": "master",
  "major": "计算机科学",
  "start_date": "2017-09-01",
  "end_date": "2020-06-30",
  "gpa": 3.8,
  "description": "硕士研究生，主修人工智能方向"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "教育背景添加成功",
  "data": {
    "education_id": 601,
    "created_at": "2025-01-03T22:00:00Z"
  },
  "timestamp": "2025-01-03T22:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.4 添加技能标签

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/skills`

**功能说明**: 为员工添加技能标签

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "skill_type": "tech",
  "skill_name": "React",
  "proficiency": "expert",
  "certified": false
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| skill_type | string | 是 | 技能类型: tech（技术）/ language（语言）/ certificate（证书）/ soft_skill（软技能） |
| skill_name | string | 是 | 技能名称 |
| proficiency | string | 是 | 熟练度: beginner（初级）/ intermediate（中级）/ advanced（高级）/ expert（专家） |
| certified | boolean | 否 | 是否有证书 |

**响应示例**:
```json
{
  "code": 0,
  "message": "技能标签添加成功",
  "data": {
    "skill_id": 701,
    "created_at": "2025-01-03T23:00:00Z"
  },
  "timestamp": "2025-01-03T23:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.5 添加绩效记录

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/performance-reviews`

**功能说明**: 为员工添加绩效记录

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "review_type": "annual",
  "review_period": "2024",
  "reviewer_id": 11,
  "rating": 4,
  "goals_achievement": 95.5,
  "feedback": "表现优秀，完成了年度目标",
  "result": "good",
  "review_date": "2024-12-31"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "绩效记录添加成功",
  "data": {
    "review_id": 801,
    "created_at": "2025-01-03T23:30:00Z"
  },
  "timestamp": "2025-01-03T23:30:00Z",
  "trace_id": "abc123"
}
```

---

### 5.6 添加培训记录

**接口地址**: `POST /api/v1/tenants/{tenant_id}/members/{member_id}/training-records`

**功能说明**: 为员工添加培训记录

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| member_id | int64 | 成员ID |

**请求参数**:
```json
{
  "training_name": "React高级培训",
  "training_type": "external",
  "provider": "某培训机构",
  "start_date": "2024-06-01",
  "end_date": "2024-06-03",
  "duration_hours": 24,
  "certificate_obtained": true,
  "notes": "学习了React性能优化、Hooks高级用法等"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "培训记录添加成功",
  "data": {
    "training_id": 901,
    "created_at": "2025-01-03T23:45:00Z"
  },
  "timestamp": "2025-01-03T23:45:00Z",
  "trace_id": "abc123"
}
```

---

## 6. 虚拟组织管理API

### 6.1 创建虚拟组织

**接口地址**: `POST /api/v1/tenants/{tenant_id}/virtual-organizations`

**功能说明**: 创建虚拟组织（项目组、委员会等）

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求参数**:
```json
{
  "name": "AI项目组",
  "type": "project",
  "description": "负责AI产品研发",
  "leader_id": 101,
  "start_date": "2025-01-01",
  "end_date": "2025-12-31",
  "member_ids": [101, 102, 103, 104]
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 虚拟组织名称 |
| type | string | 是 | 类型: project（项目组）/ committee（委员会）/ team（团队）/ other（其他） |
| description | string | 否 | 描述 |
| leader_id | int64 | 是 | 负责人ID |
| start_date | string | 是 | 开始日期 |
| end_date | string | 否 | 结束日期（NULL表示长期） |
| member_ids | array | 是 | 成员ID列表 |

**响应示例**:
```json
{
  "code": 0,
  "message": "虚拟组织创建成功",
  "data": {
    "id": "v1",
    "name": "AI项目组",
    "type": "project",
    "status": "active",
    "member_count": 4,
    "created_at": "2025-01-03T10:00:00Z"
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.2 查询虚拟组织列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/virtual-organizations`

**功能说明**: 查询虚拟组织列表

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 类型过滤 |
| status | string | 否 | 状态过滤: active/inactive/archived |
| member_id | int64 | 否 | 过滤包含指定成员的组织 |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "items": [
      {
        "id": "v1",
        "name": "AI项目组",
        "type": "project",
        "description": "负责AI产品研发",
        "leader": {
          "id": 101,
          "name": "张三"
        },
        "member_count": 4,
        "start_date": "2025-01-01",
        "end_date": "2025-12-31",
        "status": "active",
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.3 获取虚拟组织详情

**接口地址**: `GET /api/v1/tenants/{tenant_id}/virtual-organizations/{virtual_org_id}`

**功能说明**: 获取虚拟组织详情

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| virtual_org_id | string | 虚拟组织ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "v1",
    "name": "AI项目组",
    "type": "project",
    "description": "负责AI产品研发",
    "leader": {
      "id": 101,
      "name": "张三",
      "avatar": "https://cdn.example.com/avatar/101.jpg"
    },
    "members": [
      {
        "id": 101,
        "name": "张三",
        "role": "项目负责人",
        "is_leader": true,
        "joined_at": "2025-01-01T00:00:00Z"
      },
      {
        "id": 102,
        "name": "李四",
        "role": "后端开发",
        "is_leader": false,
        "joined_at": "2025-01-01T00:00:00Z"
      }
    ],
    "start_date": "2025-01-01",
    "end_date": "2025-12-31",
    "status": "active",
    "created_at": "2025-01-03T10:00:00Z"
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.4 更新虚拟组织

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/virtual-organizations/{virtual_org_id}`

**功能说明**: 更新虚拟组织信息

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| virtual_org_id | string | 虚拟组织ID |

**请求参数**:
```json
{
  "name": "AI研发项目组",
  "description": "负责AI产品研发与优化",
  "end_date": "2026-12-31"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "虚拟组织更新成功",
  "data": {
    "id": "v1",
    "updated_at": "2025-01-03T12:00:00Z"
  },
  "timestamp": "2025-01-03T12:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.5 添加虚拟组织成员

**接口地址**: `POST /api/v1/tenants/{tenant_id}/virtual-organizations/{virtual_org_id}/members`

**功能说明**: 向虚拟组织添加成员

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| virtual_org_id | string | 虚拟组织ID |

**请求参数**:
```json
{
  "member_ids": [105, 106],
  "role": "开发人员"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "成员添加成功",
  "data": {
    "added_count": 2,
    "total_members": 6,
    "updated_at": "2025-01-03T13:00:00Z"
  },
  "timestamp": "2025-01-03T13:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.6 移除虚拟组织成员

**接口地址**: `DELETE /api/v1/tenants/{tenant_id}/virtual-organizations/{virtual_org_id}/members/{member_id}`

**功能说明**: 从虚拟组织移除成员

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| virtual_org_id | string | 虚拟组织ID |
| member_id | int64 | 成员ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "成员移除成功",
  "data": {
    "virtual_org_id": "v1",
    "member_id": 106,
    "removed_at": "2025-01-03T14:00:00Z"
  },
  "timestamp": "2025-01-03T14:00:00Z",
  "trace_id": "abc123"
}
```

---

## 7. 组织代理管理API

### 7.1 创建组织代理

**接口地址**: `POST /api/v1/tenants/{tenant_id}/organization-delegates`

**功能说明**: 创建组织代理（授权委托）

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求参数**:
```json
{
  "organization_id": 11,
  "delegator_id": 11,
  "delegate_id": 102,
  "delegate_type": "manager",
  "start_time": "2025-01-10T00:00:00Z",
  "end_time": "2025-01-20T23:59:59Z",
  "reason": "经理年假期间代理审批"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| organization_id | int64 | 是 | 组织ID |
| delegator_id | int64 | 是 | 授权人ID |
| delegate_id | int64 | 是 | 被代理人ID |
| delegate_type | string | 是 | 代理类型: manager（管理）/ approval（审批）/ view（查看） |
| start_time | string | 是 | 开始时间 |
| end_time | string | 否 | 结束时间（NULL表示永久） |
| reason | string | 否 | 代理原因 |

**响应示例**:
```json
{
  "code": 0,
  "message": "组织代理创建成功",
  "data": {
    "id": 1001,
    "status": "active",
    "created_at": "2025-01-03T15:00:00Z"
  },
  "timestamp": "2025-01-03T15:00:00Z",
  "trace_id": "abc123"
}
```

---

### 7.2 查询组织代理列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/organization-delegates`

**功能说明**: 查询组织代理列表

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| organization_id | int64 | 否 | 组织ID过滤 |
| delegator_id | int64 | 否 | 授权人ID过滤 |
| delegate_id | int64 | 否 | 被代理人ID过滤 |
| delegate_type | string | 否 | 代理类型过滤 |
| status | string | 否 | 状态过滤: active/expired/cancelled |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": 1001,
        "organization": {
          "id": 11,
          "name": "前端开发部"
        },
        "delegator": {
          "id": 11,
          "name": "李四"
        },
        "delegate": {
          "id": 102,
          "name": "王五"
        },
        "delegate_type": "manager",
        "start_time": "2025-01-10T00:00:00Z",
        "end_time": "2025-01-20T23:59:59Z",
        "reason": "经理年假期间代理审批",
        "status": "active",
        "created_at": "2025-01-03T15:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T16:00:00Z",
  "trace_id": "abc123"
}
```

---

### 7.3 取消组织代理

**接口地址**: `POST /api/v1/tenants/{tenant_id}/organization-delegates/{delegate_id}/cancel`

**功能说明**: 取消组织代理

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| delegate_id | int64 | 代理ID |

**请求参数**:
```json
{
  "reason": "提前返回工作岗位"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "组织代理已取消",
  "data": {
    "id": 1001,
    "status": "cancelled",
    "cancelled_at": "2025-01-03T17:00:00Z"
  },
  "timestamp": "2025-01-03T17:00:00Z",
  "trace_id": "abc123"
}
```

---

## 8. 数据模型

### 8.1 Organization（组织）

```typescript
interface Organization {
  id: number;
  name: string;
  type: 'company' | 'department' | 'team' | 'group';
  level: number;
  parent_id: number | null;
  path: string;
  manager?: {
    id: number;
    name: string;
    avatar?: string;
    position?: string;
  };
  member_count: number;
  children?: Organization[];
  members?: Member[];
}
```

### 8.2 Member（成员）

```typescript
interface Member {
  id: number;
  name: string;
  avatar?: string;
  email: string;
  phone?: string;
  work_email?: string;
  work_phone?: string;
  office_location?: string;

  organization?: {
    id: number;
    name: string;
    path: string;
    type: string;
  };

  position?: {
    id: number;
    name: string;
    job_level?: string;
    salary_grade?: string;
  };

  direct_manager?: {
    id: number;
    name: string;
    position?: string;
    email?: string;
  };

  direct_subordinates?: Member[];

  virtual_organizations?: VirtualOrganization[];

  entry_status: 'probation' | 'regular' | 'resigned' | 'blacklist';
  joined_at: string;

  // 档案扩展
  gender?: 'male' | 'female' | 'other';
  birth_date?: string;
  marital_status?: 'single' | 'married' | 'divorced' | 'widowed';
  address?: string;
  emergency_contact_name?: string;
  emergency_contact_phone?: string;
  emergency_contact_relation?: string;

  // 工作经历
  work_experiences?: WorkExperience[];
  // 教育背景
  educations?: Education[];
  // 技能标签
  skills?: Skill[];
  // 绩效记录
  performance_reviews?: PerformanceReview[];
  // 培训记录
  training_records?: TrainingRecord[];
}
```

### 8.3 WorkExperience（工作经历）

```typescript
interface WorkExperience {
  id: number;
  company_name: string;
  industry?: string;
  position?: string;
  start_date: string;
  end_date?: string;
  description?: string;
}
```

### 8.4 Education（教育背景）

```typescript
interface Education {
  id: number;
  school_name: string;
  degree: 'high_school' | 'associate' | 'bachelor' | 'master' | 'doctor' | 'other';
  major?: string;
  start_date: string;
  end_date?: string;
  gpa?: number;
  description?: string;
}
```

### 8.5 Skill（技能标签）

```typescript
interface Skill {
  id: number;
  skill_type: 'tech' | 'language' | 'certificate' | 'soft_skill';
  skill_name: string;
  proficiency: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  certified: boolean;
}
```

### 8.6 PerformanceReview（绩效记录）

```typescript
interface PerformanceReview {
  id: number;
  review_type: 'annual' | 'project' | 'probation' | 'promotion';
  review_period?: string;
  reviewer_id?: number;
  rating?: number;
  goals_achievement?: number;
  feedback?: string;
  result?: 'excellent' | 'good' | 'satisfactory' | 'needs_improvement' | 'unsatisfactory';
  review_date: string;
}
```

### 8.7 TrainingRecord（培训记录）

```typescript
interface TrainingRecord {
  id: number;
  training_name: string;
  training_type: 'internal' | 'external' | 'online' | 'offline';
  provider?: string;
  start_date: string;
  end_date?: string;
  duration_hours?: number;
  status: 'planned' | 'ongoing' | 'completed' | 'cancelled';
  certificate_obtained: boolean;
  notes?: string;
}
```

### 8.8 VirtualOrganization（虚拟组织）

```typescript
interface VirtualOrganization {
  id: string;
  name: string;
  type: 'project' | 'committee' | 'team' | 'other';
  description?: string;
  leader?: {
    id: number;
    name: string;
    avatar?: string;
  };
  members?: VirtualOrgMember[];
  start_date: string;
  end_date?: string;
  status: 'active' | 'inactive' | 'archived';
  created_at: string;
}
```

### 8.9 VirtualOrgMember（虚拟组织成员）

```typescript
interface VirtualOrgMember {
  id: number;
  name: string;
  role?: string;
  is_leader: boolean;
  joined_at: string;
}
```

### 8.10 OrganizationDelegate（组织代理）

```typescript
interface OrganizationDelegate {
  id: number;
  organization: {
    id: number;
    name: string;
  };
  delegator: {
    id: number;
    name: string;
  };
  delegate: {
    id: number;
    name: string;
  };
  delegate_type: 'manager' | 'approval' | 'view';
  start_time: string;
  end_time?: string;
  reason?: string;
  status: 'active' | 'expired' | 'cancelled';
  created_at: string;
}
```

---

## 9. 后端代码示例

### 9.1 组织服务

```go
// application/service/organization_service.go
package service

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app/server"
)

type OrganizationService struct {
    orgRepo           repository.OrganizationRepository
    memberRepo        repository.MemberRepository
    virtualOrgRepo    repository.VirtualOrganizationRepository
    delegateRepo      repository.OrganizationDelegateRepository
    dataPermService   *DataPermissionService
    cache             *redis.Client
}

// GetContactsByDepartment 获取按部门组织的通讯录
func (s *OrganizationService) GetContactsByDepartment(
    ctx context.Context,
    tenantID string,
    includeMembers bool,
    treeView bool,
) ([]*dto.OrganizationTreeNode, error) {
    // 1. 查询所有组织
    orgs, err := s.orgRepo.GetByTenantID(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 2. 构建树形结构
    if treeView {
        return s.buildOrgTree(ctx, orgs, includeMembers)
    }

    // 3. 返回扁平列表
    return s.buildOrgList(ctx, orgs, includeMembers)
}

// buildOrgTree 构建组织树
func (s *OrganizationService) buildOrgTree(
    ctx context.Context,
    orgs []*entity.Organization,
    includeMembers bool,
) ([]*dto.OrganizationTreeNode, error) {
    // 1. 构建ID到组织的映射
    orgMap := make(map[int64]*dto.OrganizationTreeNode)
    for _, org := range orgs {
        orgMap[org.ID] = &dto.OrganizationTreeNode{
            ID:   org.ID,
            Name: org.Name,
            Type: org.Type,
            Level: org.Level,
            ParentID: org.ParentID,
            Path: org.Path,
        }
    }

    // 2. 填充负责人信息
    for _, org := range orgs {
        if org.ManagerID != nil {
            manager, _ := s.memberRepo.GetByID(ctx, *org.ManagerID)
            if manager != nil {
                node := orgMap[org.ID]
                node.Manager = &dto.MemberSummary{
                    ID:       manager.ID,
                    Name:     manager.Name,
                    Avatar:   manager.Avatar,
                    Position: manager.Position,
                }
            }
        }

        // 填充成员数量
        count, _ := s.memberRepo.CountByOrgID(ctx, org.ID)
        node := orgMap[org.ID]
        node.MemberCount = count

        // 填充成员列表
        if includeMembers {
            members, _ := s.memberRepo.GetByOrgID(ctx, org.ID)
            for _, member := range members {
                node.Members = append(node.Members, &dto.MemberSummary{
                    ID:       member.ID,
                    Name:     member.Name,
                    Avatar:   member.Avatar,
                    Position: member.Position,
                })
            }
        }
    }

    // 3. 构建树形结构
    var roots []*dto.OrganizationTreeNode
    for _, org := range orgs {
        node := orgMap[org.ID]
        if org.ParentID == nil {
            roots = append(roots, node)
        } else {
            parent := orgMap[*org.ParentID]
            parent.Children = append(parent.Children, node)
        }
    }

    return roots, nil
}

// GetContactsByPinyin 获取按拼音分组的通讯录
func (s *OrganizationService) GetContactsByPinyin(
    ctx context.Context,
    tenantID string,
    initial string,
    page, pageSize int,
) (*dto.PinyinContactsResponse, error) {
    // 1. 查询所有成员
    members, total, err := s.memberRepo.GetByTenantIDWithPinyin(
        ctx,
        tenantID,
        initial,
        page,
        pageSize,
    )
    if err != nil {
        return nil, err
    }

    // 2. 按拼音首字母分组
    grouped := make(map[string][]*dto.ContactItem)
    for _, member := range members {
        initial := s.getPinyinInitial(member.Name)
        if _, ok := grouped[initial]; !ok {
            grouped[initial] = []*dto.ContactItem{}
        }
        grouped[initial] = append(grouped[initial], &dto.ContactItem{
            ID:             member.ID,
            Name:           member.Name,
            Pinyin:         member.Pinyin,
            PinyinInitial:  initial,
            Avatar:         member.Avatar,
            Organization:   member.Organization.Name,
            Position:       member.Position,
            Phone:          s.maskPhone(member.Phone),
            Email:          member.Email,
            EntryStatus:    member.EntryStatus,
            OfficeLocation: member.OfficeLocation,
        })
    }

    return &dto.PinyinContactsResponse{
        Total:   total,
        Contacts: grouped,
    }, nil
}

// SearchContacts 搜索同事
func (s *OrganizationService) SearchContacts(
    ctx context.Context,
    tenantID string,
    keyword string,
    page, pageSize int,
) (*dto.SearchContactsResponse, error) {
    // 使用Elasticsearch全文搜索
    results, total, err := s.memberRepo.Search(
        ctx,
        tenantID,
        keyword,
        page,
        pageSize,
    )
    if err != nil {
        return nil, err
    }

    items := make([]*dto.ContactSearchItem, len(results))
    for i, member := range results {
        items[i] = &dto.ContactSearchItem{
            ID:           member.ID,
            Name:         member.Name,
            Avatar:       member.Avatar,
            Organization: member.Organization.Name,
            Position:     member.Position,
            Phone:        s.maskPhone(member.Phone),
            Email:        member.Email,
            Highlight:    s.buildHighlight(member, keyword),
        }
    }

    return &dto.SearchContactsResponse{
        Total: total,
        Items: items,
    }, nil
}
```

### 9.2 员工生命周期服务

```go
// application/service/lifecycle_service.go
package service

type LifecycleService struct {
    memberRepo          repository.MemberRepository
    positionChangeRepo  repository.PositionChangeRepository
    performanceRepo     repository.PerformanceReviewRepository
    workflowService     *WorkflowService
    notificationService *NotificationService
}

// HandleProbationReview 处理转正申请
func (s *LifecycleService) HandleProbationReview(
    ctx context.Context,
    memberID int64,
    req *dto.ProbationReviewRequest,
) (*dto.ProbationReviewResponse, error) {
    // 1. 获取成员信息
    member, err := s.memberRepo.GetByID(ctx, memberID)
    if err != nil {
        return nil, err
    }

    // 2. 验证试用期状态
    if member.EntryStatus != "probation" {
        return nil, errors.New("员工不在试用期")
    }

    // 3. 创建绩效记录
    review := &entity.PerformanceReview{
        TenantID:     member.TenantID,
        MemberID:     memberID,
        ReviewType:   "probation",
        ReviewPeriod: fmt.Sprintf("%d", time.Now().Year()),
        Result:       req.Status,
        Feedback:     req.Feedback,
        ReviewDate:   time.Now(),
    }
    s.performanceRepo.Create(ctx, review)

    // 4. 提交审批流程
    workflowID, err := s.workflowService.StartProbationApprovalWorkflow(
        ctx,
        memberID,
        req,
    )
    if err != nil {
        return nil, err
    }

    // 5. 发送通知
    s.notificationService.NotifyProbationReview(ctx, memberID, req)

    return &dto.ProbationReviewResponse{
        ReviewID:   fmt.Sprintf("%d", review.ID),
        Status:     "pending_approval",
        WorkflowID: workflowID,
    }, nil
}

// HandleTransfer 处理调岗申请
func (s *LifecycleService) HandleTransfer(
    ctx context.Context,
    memberID int64,
    req *dto.TransferRequest,
) (*dto.TransferResponse, error) {
    // 1. 获取成员信息
    member, err := s.memberRepo.GetByID(ctx, memberID)
    if err != nil {
        return nil, err
    }

    // 2. 创建调岗记录
    change := &entity.PositionChange{
        TenantID:           member.TenantID,
        MemberID:           memberID,
        ChangeType:         req.ChangeType,
        OldOrganizationID:  member.OrganizationID,
        OldPositionID:      member.PositionID,
        OldJobLevel:        member.JobLevel,
        NewOrganizationID:  req.NewOrganizationID,
        NewPositionID:      req.NewPositionID,
        NewJobLevel:        req.NewJobLevel,
        EffectiveDate:      req.EffectiveDate,
        Reason:             req.Reason,
        Status:             "pending",
    }
    err = s.positionChangeRepo.Create(ctx, change)
    if err != nil {
        return nil, err
    }

    // 3. 提交审批流程
    workflowID, err := s.workflowService.StartTransferApprovalWorkflow(
        ctx,
        memberID,
        change,
    )
    if err != nil {
        return nil, err
    }

    // 4. 发送通知
    s.notificationService.NotifyTransfer(ctx, memberID, req)

    return &dto.TransferResponse{
        ChangeID:   fmt.Sprintf("%d", change.ID),
        Status:     "pending_approval",
        WorkflowID: workflowID,
    }, nil
}

// HandleResignation 处理离职申请
func (s *LifecycleService) HandleResignation(
    ctx context.Context,
    memberID int64,
    req *dto.ResignationRequest,
) (*dto.ResignationResponse, error) {
    // 1. 获取成员信息
    member, err := s.memberRepo.GetByID(ctx, memberID)
    if err != nil {
        return nil, err
    }

    // 2. 创建离职记录
    resignation := &entity.Resignation{
        TenantID:       member.TenantID,
        MemberID:       memberID,
        LastWorkingDay: req.LastWorkingDay,
        Reason:         req.Reason,
        HandoverTo:     req.HandoverTo,
        Status:         "pending",
    }
    err = s.resignationRepo.Create(ctx, resignation)
    if err != nil {
        return nil, err
    }

    // 3. 提交审批流程
    workflowID, err := s.workflowService.StartResignationApprovalWorkflow(
        ctx,
        memberID,
        resignation,
    )
    if err != nil {
        return nil, err
    }

    // 4. 发送通知（给交接人）
    s.notificationService.NotifyResignation(ctx, memberID, req)

    return &dto.ResignationResponse{
        ResignID:   fmt.Sprintf("%d", resignation.ID),
        Status:     "pending_approval",
        WorkflowID: workflowID,
    }, nil
}
```

---

## 10. 前端代码示例

### 10.1 组织API客户端

```typescript
// api/organization-api.ts
import { request } from '@shared/utils/request';

export interface OrganizationApiClient {
  // 企业通讯录
  getContactsByDepartment(
    tenantId: string,
    includeMembers?: boolean,
    treeView?: boolean
  ): Promise<OrganizationTree[]>;
  getContactsByPinyin(
    tenantId: string,
    initial?: string,
    page?: number,
    pageSize?: number
  ): Promise<PinyinContactsResponse>;
  getContactDetail(tenantId: string, memberId: number): Promise<ContactDetail>;
  searchContacts(
    tenantId: string,
    keyword: string,
    page?: number,
    pageSize?: number
  ): Promise<SearchResult<ContactItem>>;

  // 组织架构管理
  getOrganizationTree(tenantId: string): Promise<OrganizationTree>;
  createOrganization(tenantId: string, data: CreateOrganizationRequest): Promise<Organization>;
  updateOrganization(tenantId: string, orgId: number, data: UpdateOrganizationRequest): Promise<void>;
  deleteOrganization(tenantId: string, orgId: number): Promise<void>;
  moveOrganization(tenantId: string, orgId: number, newParentId: number): Promise<void>;

  // 员工全生命周期
  submitProbationReview(tenantId: string, memberId: number, data: ProbationReviewRequest): Promise<ProbationReviewResponse>;
  submitTransfer(tenantId: string, memberId: number, data: TransferRequest): Promise<TransferResponse>;
  submitResignation(tenantId: string, memberId: number, data: ResignationRequest): Promise<ResignationResponse>;
  getPositionChanges(tenantId: string, memberId: number): Promise<PositionChange[]>;

  // 员工档案
  updateMemberProfile(tenantId: string, memberId: number, data: UpdateProfileRequest): Promise<void>;
  addWorkExperience(tenantId: string, memberId: number, data: WorkExperienceRequest): Promise<WorkExperience>;
  addEducation(tenantId: string, memberId: number, data: EducationRequest): Promise<Education>;
  addSkill(tenantId: string, memberId: number, data: SkillRequest): Promise<Skill>;

  // 虚拟组织
  createVirtualOrganization(tenantId: string, data: CreateVirtualOrgRequest): Promise<VirtualOrganization>;
  listVirtualOrganizations(tenantId: string, params?: ListVirtualOrgParams): Promise<PageResult<VirtualOrganization>>;
  getVirtualOrganization(tenantId: string, virtualOrgId: string): Promise<VirtualOrganization>;
  addVirtualOrgMembers(tenantId: string, virtualOrgId: string, memberIds: number[]): Promise<void>;
  removeVirtualOrgMember(tenantId: string, virtualOrgId: string, memberId: number): Promise<void>;

  // 组织代理
  createOrganizationDelegate(tenantId: string, data: CreateDelegateRequest): Promise<OrganizationDelegate>;
  listOrganizationDelegates(tenantId: string, params?: ListDelegateParams): Promise<PageResult<OrganizationDelegate>>;
  cancelOrganizationDelegate(tenantId: string, delegateId: number, reason?: string): Promise<void>;
}

class OrganizationApiClientImpl implements OrganizationApiClient {
  private baseURL = '/api/v1';

  // 企业通讯录
  async getContactsByDepartment(
    tenantId: string,
    includeMembers = true,
    treeView = true
  ): Promise<OrganizationTree[]> {
    const result = await request.get<{ organizations: OrganizationTree[] }>(
      `${this.baseURL}/tenants/${tenantId}/contacts/by-department`,
      { params: { include_members: includeMembers, tree_view: treeView } }
    );
    return result.organizations;
  }

  async getContactsByPinyin(
    tenantId: string,
    initial?: string,
    page = 1,
    pageSize = 50
  ): Promise<PinyinContactsResponse> {
    return request.get(
      `${this.baseURL}/tenants/${tenantId}/contacts/by-pinyin`,
      { params: { initial, page, page_size: pageSize } }
    );
  }

  async getContactDetail(tenantId: string, memberId: number): Promise<ContactDetail> {
    return request.get(`${this.baseURL}/tenants/${tenantId}/contacts/${memberId}`);
  }

  async searchContacts(
    tenantId: string,
    keyword: string,
    page = 1,
    pageSize = 20
  ): Promise<SearchResult<ContactItem>> {
    return request.get(
      `${this.baseURL}/tenants/${tenantId}/contacts/search`,
      { params: { keyword, page, page_size: pageSize } }
    );
  }

  // 组织架构管理
  async getOrganizationTree(tenantId: string): Promise<OrganizationTree> {
    const result = await request.get<OrganizationTree>(
      `${this.baseURL}/tenants/${tenantId}/organizations/tree`
    );
    return result;
  }

  async createOrganization(tenantId: string, data: CreateOrganizationRequest): Promise<Organization> {
    return request.post(`${this.baseURL}/tenants/${tenantId}/organizations`, data);
  }

  async updateOrganization(
    tenantId: string,
    orgId: number,
    data: UpdateOrganizationRequest
  ): Promise<void> {
    return request.put(`${this.baseURL}/tenants/${tenantId}/organizations/${orgId}`, data);
  }

  async deleteOrganization(tenantId: string, orgId: number): Promise<void> {
    return request.delete(`${this.baseURL}/tenants/${tenantId}/organizations/${orgId}`);
  }

  // 员工全生命周期
  async submitProbationReview(
    tenantId: string,
    memberId: number,
    data: ProbationReviewRequest
  ): Promise<ProbationReviewResponse> {
    return request.post(
      `${this.baseURL}/tenants/${tenantId}/members/${memberId}/probation-review`,
      data
    );
  }

  async submitTransfer(
    tenantId: string,
    memberId: number,
    data: TransferRequest
  ): Promise<TransferResponse> {
    return request.post(
      `${this.baseURL}/tenants/${tenantId}/members/${memberId}/transfer`,
      data
    );
  }

  async submitResignation(
    tenantId: string,
    memberId: number,
    data: ResignationRequest
  ): Promise<ResignationResponse> {
    return request.post(
      `${this.baseURL}/tenants/${tenantId}/members/${memberId}/resign`,
      data
    );
  }

  // 虚拟组织
  async createVirtualOrganization(
    tenantId: string,
    data: CreateVirtualOrgRequest
  ): Promise<VirtualOrganization> {
    return request.post(`${this.baseURL}/tenants/${tenantId}/virtual-organizations`, data);
  }

  async listVirtualOrganizations(
    tenantId: string,
    params?: ListVirtualOrgParams
  ): Promise<PageResult<VirtualOrganization>> {
    return request.get(`${this.baseURL}/tenants/${tenantId}/virtual-organizations`, { params });
  }

  // 组织代理
  async createOrganizationDelegate(
    tenantId: string,
    data: CreateDelegateRequest
  ): Promise<OrganizationDelegate> {
    return request.post(`${this.baseURL}/tenants/${tenantId}/organization-delegates`, data);
  }

  async listOrganizationDelegates(
    tenantId: string,
    params?: ListDelegateParams
  ): Promise<PageResult<OrganizationDelegate>> {
    return request.get(`${this.baseURL}/tenants/${tenantId}/organization-delegates`, { params });
  }

  async cancelOrganizationDelegate(
    tenantId: string,
    delegateId: number,
    reason?: string
  ): Promise<void> {
    return request.post(
      `${this.baseURL}/tenants/${tenantId}/organization-delegates/${delegateId}/cancel`,
      { reason }
    );
  }
}

export const organizationApi = new OrganizationApiClientImpl();
```

### 10.2 企业通讯录组件

```typescript
// pages/contacts/ContactsPage.tsx
import React, { useState } from 'react';
import { Tabs, Tree, Input, Avatar, Card, Modal } from '@douyinfe/semi-ui';
import { organizationApi } from '@api/organization-api';

export const ContactsPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('department');
  const [searchKeyword, setSearchKeyword] = useState('');
  const [selectedMember, setSelectedMember] = useState(null);
  const [orgTree, setOrgTree] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadContactsByDepartment();
  }, []);

  const loadContactsByDepartment = async () => {
    setLoading(true);
    try {
      const tree = await organizationApi.getContactsByDepartment('tenant-123');
      setOrgTree(tree);
    } catch (error) {
      console.error('加载通讯录失败', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (keyword: string) => {
    setSearchKeyword(keyword);
    if (keyword.length >= 2) {
      const result = await organizationApi.searchContacts('tenant-123', keyword);
      console.log('搜索结果:', result);
    }
  };

  return (
    <div className="contacts-page">
      <div className="search-bar">
        <Input
          placeholder="搜索姓名/手机/邮箱/岗位"
          value={searchKeyword}
          onChange={setSearchSearchKeyword}
          showClear
        />
      </div>

      <Tabs activeKey={activeTab} onChange={setActiveTab}>
        <Tabs.TabPane tab="按部门" itemKey="department">
          <Tree
            treeData={orgTree}
            loading={loading}
            renderTitle={(nodeData) => (
              <div className="org-node">
                <Avatar size="small">{nodeData.name[0]}</Avatar>
                <span className="org-name">{nodeData.name}</span>
                <span className="member-count">({nodeData.member_count})</span>
              </div>
            )}
          />
        </Tabs.TabPane>

        <Tabs.TabPane tab="按拼音" itemKey="pinyin">
          <ContactListByPinyin />
        </Tabs.TabPane>
      </Tabs>

      <Modal visible={!!selectedMember} footer={null}>
        <ContactDetail member={selectedMember} />
      </Modal>
    </div>
  );
};
```

---

## 11. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 20001 | 400 | 组织名称已存在 |
| 20002 | 404 | 组织不存在 |
| 20003 | 400 | 组织仍有成员，无法删除 |
| 20004 | 400 | 无法移动到子组织下 |
| 20005 | 400 | 组织循环引用 |
| 20101 | 404 | 成员不存在 |
| 20102 | 403 | 无权限查看该成员信息 |
| 20103 | 400 | 成员状态异常，无法操作 |
| 20104 | 400 | 试用期未结束，无法转正 |
| 20105 | 400 | 已有审批中的调岗申请 |
| 20106 | 400 | 离职申请已存在 |
| 20201 | 400 | 虚拟组织名称已存在 |
| 20202 | 404 | 虚拟组织不存在 |
| 20203 | 400 | 成员已在该虚拟组织中 |
| 20204 | 400 | 虚拟组织已过期 |
| 20301 | 400 | 代理关系已存在 |
| 20302 | 400 | 代理时间范围无效 |
| 20303 | 404 | 代理记录不存在 |
| 20304 | 400 | 代理已过期或已取消 |

---

**文档结束**
