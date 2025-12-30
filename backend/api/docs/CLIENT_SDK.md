# 组织中心API客户端SDK文档

**版本**: v1.0.0
**最后更新**: 2025-01-01

---

## 📖 目录

- [SDK概览](#sdk概览)
- [Go SDK](#go-sdk)
- [JavaScript/TypeScript SDK](#javascripttypescript-sdk)
- [Python SDK](#python-sdk)
- [Java SDK](#java-sdk)
- [最佳实践](#最佳实践)

---

## SDK概览

组织中心API提供了多语言客户端SDK，简化API调用流程，提供类型安全和更好的开发体验。

### 支持的语言

| 语言 | SDK版本 | 状态 | 文档链接 |
|------|---------|------|---------|
| Go | v1.0.0 | ✅ 稳定 | [查看文档](#go-sdk) |
| JavaScript/TypeScript | v1.0.0 | ✅ 稳定 | [查看文档](#javascripttypescript-sdk) |
| Python | v1.0.0 | ✅ 稳定 | [查看文档](#python-sdk) |
| Java | v1.0.0 | ✅ 稳定 | [查看文档](#java-sdk) |

### 核心功能

所有SDK提供以下核心功能：

- ✅ 自动Token管理和刷新
- ✅ 请求/响应拦截
- ✅ 错误处理和重试
- ✅ 类型安全（TypeScript/Java）
- ✅ 日志记录
- ✅ 超时配置
- ✅ 租户隔离

---

## Go SDK

### 安装

```bash
go get github.com/coze-dev/coze-studio-sdk-go
```

### 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    orgsdk "github.com/coze-dev/coze-studio-sdk-go/org"
)

func main() {
    // 1. 创建客户端
    client := orgsdk.NewClient(&orgsdk.Config{
        BaseURL:    "https://api.coze.com",
        APIKey:     "your-api-key",
        TenantID:   "tenant_001",
        Timeout:    30 * time.Second,
    })

    // 2. 创建组织
    org, err := client.Organizations.Create(context.Background(), &orgsdk.CreateOrganizationRequest{
        TenantID:    "tenant_001",
        OrgName:     "技术部",
        OrgType:     "department",
        OrgCode:     "TECH",
        Description: "负责产品研发",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("创建组织成功: %s\n", org.OrgName)

    // 3. 查询组织列表
    list, err := client.Organizations.List(context.Background(), &orgsdk.ListOrganizationsRequest{
        TenantID: "tenant_001",
        Page:     1,
        PageSize: 20,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("共有 %d 个组织\n", list.Total)
}
```

### 配置选项

```go
type Config struct {
    // API基础URL
    BaseURL string

    // API密钥或JWT Token
    APIKey  string
    Token   string

    // 租户ID（可选，也可以在每次请求时指定）
    TenantID string

    // 请求超时时间
    Timeout time.Duration

    // 重试配置
    MaxRetries int
    RetryDelay time.Duration

    // 日志配置
    LogLevel LogLevel
    Logger   Logger

    // HTTP客户端（可选）
    HTTPClient *http.Client

    // 自定义拦截器
    Interceptors []Interceptor
}
```

### 组织管理

```go
// 创建组织
org, err := client.Organizations.Create(ctx, &CreateOrganizationRequest{
    TenantID:    "tenant_001",
    OrgName:     "研发中心",
    OrgType:     "division",
    ParentID:    stringPtr("org_company_001"),
    OrgCode:     "R&D",
    LeaderID:    stringPtr("emp_001"),
    Description: "负责产品研发",
    SortOrder:   1,
})

// 获取组织详情
org, err := client.Organizations.Get(ctx, "org_001")

// 更新组织
org, err := client.Organizations.Update(ctx, &UpdateOrganizationRequest{
    OrgID:       "org_001",
    OrgName:     "新名称",
    Description: "新描述",
})

// 删除组织
err := client.Organizations.Delete(ctx, "org_001")

// 获取组织树
tree, err := client.Organizations.GetTree(ctx, "tenant_001")

// 分页查询组织列表
list, err := client.Organizations.List(ctx, &ListOrganizationsRequest{
    TenantID:  "tenant_001",
    OrgType:   "department",
    Status:    "active",
    Page:      1,
    PageSize:  20,
    SortBy:    "created_at",
    SortOrder: "desc",
})

// 移动组织
err := client.Organizations.Move(ctx, &MoveOrganizationRequest{
    OrgID:       "org_001",
    NewParentID: stringPtr("org_new_parent_001"),
})

// 获取子组织
children, err := client.Organizations.GetChildren(ctx, "org_001")

// 获取祖先组织
ancestors, err := client.Organizations.GetAncestors(ctx, "org_001")

// 获取后代组织
descendants, err := client.Organizations.GetDescendants(ctx, "org_001")
```

### 员工管理

```go
// 创建员工
emp, err := client.Employees.Create(ctx, &CreateEmployeeRequest{
    TenantID:      "tenant_001",
    OrgID:         "org_001",
    DeptID:        stringPtr("dept_001"),
    PositionID:    stringPtr("pos_001"),
    EmpName:       "张三",
    EmpCode:       "E001",
    EmployeeType:  "full_time",
    JobLevel:      5,
    JobTitle:      "高级工程师",
    HireDate:      time.Now().UnixMilli(),
    ProbationDays: 90,
    Email:         stringPtr("zhangsan@example.com"),
    Phone:         stringPtr("13800138000"),
})

// 获取员工详情
emp, err := client.Employees.Get(ctx, "emp_001")

// 更新员工信息
emp, err := client.Employees.Update(ctx, &UpdateEmployeeRequest{
    EmpID:    "emp_001",
    DeptID:   stringPtr("dept_002"),
    JobLevel: intPtr(6),
    Email:    stringPtr("new@example.com"),
})

// 更新员工状态
err := client.Employees.UpdateStatus(ctx, "emp_001", "active")

// 删除员工（软删除）
err := client.Employees.Delete(ctx, "emp_001")

// 分页查询员工列表
list, err := client.Employees.List(ctx, &ListEmployeesRequest{
    TenantID:    "tenant_001",
    DeptID:      stringPtr("dept_001"),
    Status:      "active",
    PageSize:    20,
})

// 根据工号查询
emp, err := client.Employees.GetByCode(ctx, "tenant_001", "E001")

// 搜索员工
results, err := client.Employees.Search(ctx, "tenant_001", "张三", 20)
```

### 错误处理

```go
org, err := client.Organizations.Create(ctx, req)
if err != nil {
    // 检查错误类型
    if apiErr, ok := err.(*orgsdk.APIError); ok {
        fmt.Printf("API Error: %s (code: %d)\n", apiErr.Message, apiErr.Code)

        // 处理特定错误码
        switch apiErr.Code {
        case 40901:
            fmt.Println("组织编码重复")
        case 40301:
            fmt.Println("权限不足")
        default:
            fmt.Println("未知错误")
        }
    } else {
        // 网络错误或其他错误
        log.Fatal(err)
    }
}
```

### 自定义拦截器

```go
// 日志拦截器
func loggingInterceptor(ctx context.Context, req *http.Request) error {
    fmt.Printf("[Request] %s %s\n", req.Method, req.URL)
    return nil
}

client.AddInterceptor(loggingInterceptor)

// 重试拦截器
client.AddInterceptor(func(ctx context.Context, req *http.Request) error {
    var lastErr error
    for i := 0; i < 3; i++ {
        err := client.Do(ctx, req)
        if err == nil {
            return nil
        }
        lastErr = err
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return lastErr
})
```

---

## JavaScript/TypeScript SDK

### 安装

```bash
npm install @coze-dev/org-sdk
# 或
yarn add @coze-dev/org-sdk
```

### 快速开始

```typescript
import { OrgClient } from '@coze-dev/org-sdk';

// 1. 创建客户端
const client = new OrgClient({
  baseURL: 'https://api.coze.com',
  token: 'your-jwt-token',
  tenantId: 'tenant_001',
  timeout: 30000,
});

// 2. 创建组织
const org = await client.organizations.create({
  tenantId: 'tenant_001',
  orgName: '技术部',
  orgType: OrgType.Department,
  orgCode: 'TECH',
  description: '负责产品研发',
});

console.log('创建组织成功:', org.orgName);

// 3. 查询组织列表
const list = await client.organizations.list({
  tenantId: 'tenant_001',
  page: 1,
  pageSize: 20,
});

console.log(`共有 ${list.total} 个组织`);
```

### 配置选项

```typescript
interface ClientConfig {
  // API基础URL
  baseURL: string;

  // JWT Token或API Key
  token?: string;
  apiKey?: string;

  // 租户ID
  tenantId?: string;

  // 请求超时时间（毫秒）
  timeout?: number;

  // 重试配置
  maxRetries?: number;
  retryDelay?: number;

  // 日志配置
  logLevel?: 'debug' | 'info' | 'warn' | 'error';
  logger?: Logger;

  // HTTP Agent（Node.js）
  httpAgent?: any;
  httpsAgent?: any;

  // 自定义拦截器
  interceptors?: {
    request?: (config: RequestConfig) => RequestConfig | Promise<RequestConfig>;
    response?: (response: Response) => Response | Promise<Response>;
    error?: (error: Error) => any | Promise<any>;
  };
}
```

### 组织管理

```typescript
// 创建组织
const org = await client.organizations.create({
  tenantId: 'tenant_001',
  orgName: '研发中心',
  orgType: OrgType.Division,
  parentId: 'org_company_001',
  orgCode: 'R&D',
  leaderId: 'emp_001',
  description: '负责产品研发',
  sortOrder: 1,
});

// 获取组织详情
const org = await client.organizations.get('org_001');

// 更新组织
const updated = await client.organizations.update('org_001', {
  orgName: '新名称',
  description: '新描述',
});

// 删除组织
await client.organizations.delete('org_001');

// 获取组织树
const tree = await client.organizations.getTree('tenant_001');

// 分页查询
const list = await client.organizations.list({
  tenantId: 'tenant_001',
  orgType: OrgType.Department,
  status: OrgStatus.Active,
  page: 1,
  pageSize: 20,
  sortBy: 'created_at',
  sortOrder: 'desc',
});

// 移动组织
await client.organizations.move('org_001', {
  newParentId: 'org_new_parent_001',
});

// 获取子组织
const children = await client.organizations.getChildren('org_001');

// 获取祖先组织
const ancestors = await client.organizations.getAncestors('org_001');

// 获取后代组织
const descendants = await client.organizations.getDescendants('org_001');
```

### 员工管理

```typescript
// 创建员工
const emp = await client.employees.create({
  tenantId: 'tenant_001',
  orgId: 'org_001',
  deptId: 'dept_001',
  positionId: 'pos_001',
  empName: '张三',
  empCode: 'E001',
  employeeType: EmployeeType.FullTime,
  jobLevel: 5,
  jobTitle: '高级工程师',
  hireDate: Date.now(),
  probationDays: 90,
  email: 'zhangsan@example.com',
  phone: '13800138000',
});

// 获取员工详情
const emp = await client.employees.get('emp_001');

// 更新员工
const updated = await client.employees.update('emp_001', {
  deptId: 'dept_002',
  jobLevel: 6,
  email: 'new@example.com',
});

// 更新状态
await client.employees.updateStatus('emp_001', EmployeeStatus.Active);

// 删除员工
await client.employees.delete('emp_001');

// 分页查询
const list = await client.employees.list({
  tenantId: 'tenant_001',
  deptId: 'dept_001',
  status: EmployeeStatus.Active,
  pageSize: 20,
});

// 根据工号查询
const emp = await client.employees.getByCode('tenant_001', 'E001');

// 搜索员工
const results = await client.employees.search('tenant_001', '张三', 20);
```

### 错误处理

```typescript
try {
  const org = await client.organizations.create(req);
} catch (error) {
  if (error instanceof APIError) {
    console.error(`API Error: ${error.message} (code: ${error.code})`);

    // 处理特定错误码
    switch (error.code) {
      case ErrorCode.OrganizationCodeDuplicate:
        console.log('组织编码重复');
        break;
      caseErrorCode.InsufficientPermissions:
        console.log('权限不足');
        break;
      default:
        console.log('未知错误');
    }
  } else {
    // 网络错误或其他错误
    console.error('Error:', error);
  }
}
```

### React Hook

```typescript
import { useOrganization, useEmployeeList } from '@coze-dev/org-sdk/react';

function OrganizationDetail({ orgId }: { orgId: string }) {
  const { data: org, loading, error } = useOrganization(orgId);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h1>{org.orgName}</h1>
      <p>{org.description}</p>
    </div>
  );
}

function EmployeeList({ tenantId }: { tenantId: string }) {
  const { data, loading, error, refetch } = useEmployeeList({
    tenantId,
    pageSize: 20,
  });

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data.employees.map(emp => (
        <div key={emp.empId}>{emp.empName}</div>
      ))}
    </div>
  );
}
```

### 自定义拦截器

```typescript
// 请求拦截器
const client = new OrgClient({
  interceptors: {
    request: (config) => {
      console.log(`[Request] ${config.method?.toUpperCase()} ${config.url}`);
      // 添加自定义Header
      config.headers = {
        ...config.headers,
        'X-Custom-Header': 'value',
      };
      return config;
    },
    response: (response) => {
      console.log(`[Response] ${response.status}`);
      return response;
    },
    error: (error) => {
      console.error('[Error]', error);
      // 错误重试
      if (error.code === ErrorCode.RateLimitExceeded) {
        return new Promise((resolve) => {
          setTimeout(() => {
            // 重试请求
          }, 1000);
        });
      }
      throw error;
    },
  },
});
```

---

## Python SDK

### 安装

```bash
pip install coze-org-sdk
```

### 快速开始

```python
from coze_org_sdk import OrgClient
from coze_org_sdk.models import OrgType

# 1. 创建客户端
client = OrgClient(
    base_url="https://api.coze.com",
    token="your-jwt-token",
    tenant_id="tenant_001",
    timeout=30,
)

# 2. 创建组织
org = client.organizations.create(
    tenant_id="tenant_001",
    org_name="技术部",
    org_type=OrgType.DEPARTMENT,
    org_code="TECH",
    description="负责产品研发",
)

print(f"创建组织成功: {org.org_name}")

# 3. 查询组织列表
result = client.organizations.list(
    tenant_id="tenant_001",
    page=1,
    page_size=20,
)

print(f"共有 {result.total} 个组织")
```

### 配置选项

```python
class ClientConfig:
    # API基础URL
    base_url: str

    # JWT Token或API Key
    token: Optional[str] = None
    api_key: Optional[str] = None

    # 租户ID
    tenant_id: Optional[str] = None

    # 请求超时时间（秒）
    timeout: int = 30

    # 重试配置
    max_retries: int = 3
    retry_delay: int = 1

    # 日志配置
    log_level: str = "INFO"
    logger: Optional[Logger] = None

    # 自定义Session
    session: Optional[requests.Session] = None

    # SSL验证
    verify_ssl: bool = True

    # 代理配置
    proxies: Optional[Dict[str, str]] = None
```

### 组织管理

```python
# 创建组织
org = client.organizations.create(
    tenant_id="tenant_001",
    org_name="研发中心",
    org_type=OrgType.DIVISION,
    parent_id="org_company_001",
    org_code="R&D",
    leader_id="emp_001",
    description="负责产品研发",
    sort_order=1,
)

# 获取组织详情
org = client.organizations.get("org_001")

# 更新组织
updated = client.organizations.update(
    org_id="org_001",
    org_name="新名称",
    description="新描述",
)

# 删除组织
client.organizations.delete("org_001")

# 获取组织树
tree = client.organizations.get_tree("tenant_001")

# 分页查询
result = client.organizations.list(
    tenant_id="tenant_001",
    org_type=OrgType.DEPARTMENT,
    status="active",
    page=1,
    page_size=20,
    sort_by="created_at",
    sort_order="desc",
)

# 移动组织
client.organizations.move(
    org_id="org_001",
    new_parent_id="org_new_parent_001",
)

# 获取子组织
children = client.organizations.get_children("org_001")

# 获取祖先组织
ancestors = client.organizations.get_ancestors("org_001")

# 获取后代组织
descendants = client.organizations.get_descendants("org_001")
```

### 员工管理

```python
from coze_org_sdk.models import EmployeeType

# 创建员工
emp = client.employees.create(
    tenant_id="tenant_001",
    org_id="org_001",
    dept_id="dept_001",
    position_id="pos_001",
    emp_name="张三",
    emp_code="E001",
    employee_type=EmployeeType.FULL_TIME,
    job_level=5,
    job_title="高级工程师",
    hire_date=int(time.time() * 1000),
    probation_days=90,
    email="zhangsan@example.com",
    phone="13800138000",
)

# 获取员工详情
emp = client.employees.get("emp_001")

# 更新员工
updated = client.employees.update(
    emp_id="emp_001",
    dept_id="dept_002",
    job_level=6,
    email="new@example.com",
)

# 更新状态
client.employees.update_status("emp_001", "active")

# 删除员工
client.employees.delete("emp_001")

# 分页查询
result = client.employees.list(
    tenant_id="tenant_001",
    dept_id="dept_001",
    status="active",
    page_size=20,
)

# 根据工号查询
emp = client.employees.get_by_code("tenant_001", "E001")

# 搜索员工
results = client.employees.search("tenant_001", "张三", limit=20)
```

### 错误处理

```python
from coze_org_sdk.exceptions import APIError, OrganizationCodeDuplicate

try:
    org = client.organizations.create(**kwargs)
except APIError as e:
    print(f"API Error: {e.message} (code: {e.code})")

    # 处理特定错误
    if e.code == 40901:
        print("组织编码重复")
    elif e.code == 40301:
        print("权限不足")
except Exception as e:
    # 网络错误或其他错误
    print(f"Error: {e}")
```

### 异步支持

```python
import asyncio
from coze_org_sdk import AsyncOrgClient

async def main():
    client = AsyncOrgClient(
        base_url="https://api.coze.com",
        token="your-token",
    )

    # 异步创建组织
    org = await client.organizations.create(
        tenant_id="tenant_001",
        org_name="技术部",
        org_type=OrgType.DEPARTMENT,
        org_code="TECH",
    )

    print(f"创建成功: {org.org_name}")

asyncio.run(main())
```

### 上下文管理器

```python
# 使用with语句自动关闭连接
with OrgClient(token="your-token") as client:
    org = client.organizations.get("org_001")
    print(org.org_name)
```

---

## Java SDK

### 安装（Maven）

```xml
<dependency>
    <groupId>com.coze</groupId>
    <artifactId>org-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Gradle

```gradle
implementation 'com.coze:org-sdk:1.0.0'
```

### 快速开始

```java
import com.coze.org.sdk.*;
import com.coze.org.sdk.models.*;

public class Main {
    public static void main(String[] args) {
        // 1. 创建客户端
        OrgClient client = OrgClient.builder()
            .baseUrl("https://api.coze.com")
            .token("your-jwt-token")
            .tenantId("tenant_001")
            .timeout(30)
            .build();

        // 2. 创建组织
        Organization org = client.organizations().create(CreateOrganizationRequest.builder()
            .tenantId("tenant_001")
            .orgName("技术部")
            .orgType(OrgType.DEPARTMENT)
            .orgCode("TECH")
            .description("负责产品研发")
            .build());

        System.out.println("创建组织成功: " + org.getOrgName());

        // 3. 查询组织列表
        PaginatedResponse<Organization> list = client.organizations().list(ListOrganizationsRequest.builder()
            .tenantId("tenant_001")
            .page(1)
            .pageSize(20)
            .build());

        System.out.println("共有 " + list.getTotal() + " 个组织");
    }
}
```

### 配置选项

```java
OrgClient client = OrgClient.builder()
    .baseUrl("https://api.coze.com")
    .token("your-jwt-token")
    .apiKey("your-api-key")
    .tenantId("tenant_001")
    .timeout(30)  // 秒
    .maxRetries(3)
    .retryDelay(1)  // 秒
    .logLevel(LogLevel.INFO)
    .build();
```

### 组织管理

```java
// 创建组织
Organization org = client.organizations().create(CreateOrganizationRequest.builder()
    .tenantId("tenant_001")
    .orgName("研发中心")
    .orgType(OrgType.DIVISION)
    .parentId("org_company_001")
    .orgCode("R&D")
    .leaderId("emp_001")
    .description("负责产品研发")
    .sortOrder(1)
    .build());

// 获取组织详情
Organization org = client.organizations().get("org_001");

// 更新组织
Organization updated = client.organizations().update(UpdateOrganizationRequest.builder()
    .orgId("org_001")
    .orgName("新名称")
    .description("新描述")
    .build());

// 删除组织
client.organizations().delete("org_001");

// 获取组织树
List<Organization> tree = client.organizations().getTree("tenant_001");

// 分页查询
PaginatedResponse<Organization> list = client.organizations().list(ListOrganizationsRequest.builder()
    .tenantId("tenant_001")
    .orgType(OrgType.DEPARTMENT)
    .status("active")
    .page(1)
    .pageSize(20)
    .sortBy("created_at")
    .sortOrder("desc")
    .build());

// 移动组织
client.organizations().move(MoveOrganizationRequest.builder()
    .orgId("org_001")
    .newParentId("org_new_parent_001")
    .build());
```

### 员工管理

```java
import com.coze.org.sdk.models.EmployeeType;

// 创建员工
Employee emp = client.employees().create(CreateEmployeeRequest.builder()
    .tenantId("tenant_001")
    .orgId("org_001")
    .deptId("dept_001")
    .positionId("pos_001")
    .empName("张三")
    .empCode("E001")
    .employeeType(EmployeeType.FULL_TIME)
    .jobLevel(5)
    .jobTitle("高级工程师")
    .hireDate(System.currentTimeMillis())
    .probationDays(90)
    .email("zhangsan@example.com")
    .phone("13800138000")
    .build());

// 获取员工详情
Employee emp = client.employees().get("emp_001");

// 更新员工
Employee updated = client.employees().update(UpdateEmployeeRequest.builder()
    .empId("emp_001")
    .deptId("dept_002")
    .jobLevel(6)
    .email("new@example.com")
    .build());

// 更新状态
client.employees().updateStatus("emp_001", "active");

// 删除员工
client.employees().delete("emp_001");

// 分页查询
PaginatedResponse<Employee> list = client.employees().list(ListEmployeesRequest.builder()
    .tenantId("tenant_001")
    .deptId("dept_001")
    .status("active")
    .pageSize(20)
    .build());
```

### 错误处理

```java
try {
    Organization org = client.organizations().create(request);
} catch (APIException e) {
    System.err.println("API Error: " + e.getMessage() + " (code: " + e.getCode() + ")");

    // 处理特定错误码
    switch (e.getCode()) {
        case 40901:
            System.out.println("组织编码重复");
            break;
        case 40301:
            System.out.println("权限不足");
            break;
        default:
            System.out.println("未知错误");
    }
} catch (IOException e) {
    // 网络错误
    System.err.println("Network Error: " + e.getMessage());
}
```

---

## 最佳实践

### 1. Token管理

```javascript
// 自动刷新Token
class TokenManager {
  constructor(client) {
    this.client = client;
    this.refreshTimer = null;
  }

  startAutoRefresh() {
    const expiresAt = this.getTokenExpiresAt();
    const refreshBefore = 5 * 60 * 1000; // 提前5分钟刷新

    this.refreshTimer = setTimeout(() => {
      this.refreshToken();
      this.startAutoRefresh();
    }, expiresAt - Date.now() - refreshBefore);
  }

  async refreshToken() {
    const newToken = await this.client.auth.refreshToken();
    localStorage.setItem('api_token', newToken);
    this.client.setToken(newToken);
  }

  stopAutoRefresh() {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
    }
  }
}
```

### 2. 请求缓存

```python
from functools import lru_cache

class CachedOrgClient:
    def __init__(self, client):
        self.client = client
        self.cache = {}

    def get_organization(self, org_id, use_cache=True):
        if use_cache and org_id in self.cache:
            return self.cache[org_id]

        org = self.client.organizations.get(org_id)
        self.cache[org_id] = org
        return org

    def clear_cache(self):
        self.cache.clear()
```

### 3. 批量操作

```go
// 批量创建员工
func (c *EmployeeClient) BatchCreate(ctx context.Context, employees []CreateEmployeeRequest) ([]Employee, error) {
    results := make([]Employee, 0, len(employees))
    errors := make([]error, 0)

    // 使用并发控制
    sem := make(chan struct{}, 10) // 最多10个并发
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, emp := range employees {
        wg.Add(1)
        sem <- struct{}{}

        go func(req CreateEmployeeRequest) {
            defer wg.Done()
            defer func() { <-sem }()

            emp, err := c.Create(ctx, req)
            mu.Lock()
            defer mu.Unlock()

            if err != nil {
                errors = append(errors, err)
            } else {
                results = append(results, *emp)
            }
        }(emp)
    }

    wg.Wait()

    if len(errors) > 0 {
        return results, fmt.Errorf("batch create failed with %d errors", len(errors))
    }

    return results, nil
}
```

### 4. 超时和重试

```typescript
import axios, { AxiosInstance } from 'axios';

export class ResilientClient {
  private client: AxiosInstance;

  constructor(config: ClientConfig) {
    this.client = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout || 30000,
    });

    this.setupRetryInterceptor();
  }

  private setupRetryInterceptor() {
    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        const config = error.config;

        if (!config || !config.__retryCount) {
          config.__retryCount = 0;
        }

        const maxRetries = 3;
        if (config.__retryCount >= maxRetries) {
          return Promise.reject(error);
        }

        config.__retryCount += 1;

        // 指数退避
        const delay = Math.pow(2, config.__retryCount) * 1000;
        await new Promise((resolve) => setTimeout(resolve, delay));

        return this.client(config);
      }
    );
  }
}
```

### 5. 日志和监控

```python
import logging
import time

class MonitoredClient:
    def __init__(self, client):
        self.client = client
        self.logger = logging.getLogger(__name__)

    def __getattr__(self, name):
        attr = getattr(self.client, name)
        if callable(attr):
            def wrapper(*args, **kwargs):
                start_time = time.time()
                self.logger.info(f"Calling {name} with args={args}")

                try:
                    result = attr(*args, **kwargs)
                    elapsed = time.time() - start_time
                    self.logger.info(f"{name} completed in {elapsed:.2f}s")
                    return result
                except Exception as e:
                    elapsed = time.time() - start_time
                    self.logger.error(f"{name} failed after {elapsed:.2f}s: {e}")
                    raise

            return wrapper
        return attr
```

---

## 技术支持

如有问题，请联系：

- 邮箱: sdk-support@coze.com
- GitHub Issues: https://github.com/coze-dev/coze-studio-sdk/issues
- 文档: https://docs.coze.com/sdk
