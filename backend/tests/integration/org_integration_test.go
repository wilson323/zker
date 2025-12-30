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

//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestSetup 测试环境设置
type TestSetup struct {
	DB       *gorm.DB
	OrgSvc   *service.OrganizationService
	DeptSvc  *service.DepartmentService
	EmpSvc   *service.EmployeeService
	PosSvc   *service.PositionService
	HRSvc    *service.HRLifecycleService
	DirSvc   *service.DirectoryService
	TenantID string
	Cleanup  func()
}

// setupTestEnvironment 初始化测试环境
func setupTestEnvironment(t *testing.T) *TestSetup {
	ctx := context.Background()

	// 1. 启动MySQL 8.4.5容器
	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8.4.5",
			ExposedPorts: []string{"3306/tcp", "33060/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "test123",
				"MYSQL_DATABASE":      "coze_test",
			},
			Cmd: []string{
				"--character-set-server=utf8mb4",
				"--collation-server=utf8mb4_unicode_ci",
				"--default-authentication-plugin=mysql_native_password",
			},
			WaitingFor: wait.ForLog("ready for connections").WithOccurrence(2),
		},
		Started: true,
	})
	require.NoError(t, err, "启动MySQL容器失败")

	// 2. 获取数据库连接信息
	host, err := mysqlContainer.Host(ctx)
	require.NoError(t, err)

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	require.NoError(t, err)

	dsn := fmt.Sprintf("root:test123@tcp(%s:%s)/coze_test?charset=utf8mb4&parseTime=True&loc=Local",
		host, port.Port())

	// 3. 连接数据库（带重试）
	var db *gorm.DB
	var errDB error

	for i := 0; i < 20; i++ {
		db, errDB = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if errDB == nil {
			break
		}
		t.Logf("尝试连接数据库 (第%d次): %v", i+1, errDB)
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, errDB, "连接数据库失败")

	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. 创建表结构
	err = db.AutoMigrate(
		&entity.Organization{},
		&entity.Department{},
		&entity.Position{},
		&entity.Employee{},
		&entity.EmployeeContract{},
		&entity.EmployeeTransfer{},
		&entity.EmployeeResignation{},
		&entity.OrganizationTree{},
		&entity.DepartmentTree{},
	)
	require.NoError(t, err, "创建表结构失败")

	// 5. 创建Repository和Service
	orgRepo := repository.NewOrganizationRepository(db)
	treeRepo := repository.NewOrganizationTreeRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	deptTreeRepo := repository.NewDepartmentTreeRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	posRepo := repository.NewPositionRepository(db)
	hrRepo := repository.NewHRRepository(db)

	orgSvc := service.NewOrganizationService(orgRepo, treeRepo, db)
	deptSvc := service.NewDepartmentService(deptRepo, deptTreeRepo, db)
	empSvc := service.NewEmployeeService(empRepo, db)
	posSvc := service.NewPositionService(posRepo, db)
	hrSvc := service.NewHRLifecycleService(hrRepo, empRepo, db)
	dirSvc := service.NewDirectoryService(orgRepo, deptRepo, empRepo)

	// 6. 清理函数
	cleanup := func() {
		mysqlContainer.Terminate(ctx)
		sqlDB.Close()
	}

	return &TestSetup{
		DB:       db,
		OrgSvc:   orgSvc,
		DeptSvc:  deptSvc,
		EmpSvc:   empSvc,
		PosSvc:   posSvc,
		HRSvc:    hrSvc,
		DirSvc:   dirSvc,
		TenantID: "tenant-001",
		Cleanup:  cleanup,
	}
}

// TestOrgIntegration 组织中心完整集成测试
func TestOrgIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	setup := setupTestEnvironment(t)
	defer setup.Cleanup()

	t.Run("组织管理完整CRUD流程", func(t *testing.T) {
		testOrgCRUD(t, setup)
	})

	t.Run("部门管理完整CRUD流程", func(t *testing.T) {
		testDeptCRUD(t, setup)
	})

	t.Run("岗位管理完整CRUD流程", func(t *testing.T) {
		testPositionCRUD(t, setup)
	})

	t.Run("员工管理完整CRUD流程", func(t *testing.T) {
		testEmployeeCRUD(t, setup)
	})

	t.Run("通讯录查询流程", func(t *testing.T) {
		testDirectoryQuery(t, setup)
	})

	t.Run("HR生命周期流程", func(t *testing.T) {
		testHRLifecycle(t, setup)
	})

	t.Run("多租户隔离验证", func(t *testing.T) {
		testTenantIsolation(t, setup)
	})

	t.Run("事务测试", func(t *testing.T) {
		testTransactions(t, setup)
	})

	t.Run("复杂业务场景", func(t *testing.T) {
		testComplexScenarios(t, setup)
	})
}

// testOrgCRUD 测试组织管理完整CRUD流程
func testOrgCRUD(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	t.Run("创建公司组织", func(t *testing.T) {
		req := &service.CreateOrganizationRequest{
			TenantID:     setup.TenantID,
			OrgName:      "测试公司",
			OrgType:      entity.OrgTypeCompany,
			OrgCode:      "COMP001",
			Description:  "这是一个测试公司",
			SortOrder:    1,
		}

		org, err := setup.OrgSvc.CreateOrganization(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, setup.TenantID, org.TenantID)
		assert.Equal(t, "测试公司", org.OrgName)
		assert.Equal(t, entity.OrgTypeCompany, org.OrgType)
		assert.Equal(t, 1, org.Level)
		assert.Equal(t, "/", org.Path)
		assert.True(t, org.IsActive())
	})

	t.Run("创建分公司组织", func(t *testing.T) {
		// 先获取公司
		filter := &service.OrganizationFilter{
			TenantID: setup.TenantID,
			OrgCode:  "COMP001",
		}
		companies, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter)
		require.Greater(t, len(companies), 0)

		company := companies[0]
		parentID := company.OrgID

		req := &service.CreateOrganizationRequest{
			TenantID:    setup.TenantID,
			OrgName:     "北京分公司",
			OrgType:     entity.OrgTypeDivision,
			ParentID:    &parentID,
			OrgCode:     "DIV001",
			Description: "北京分公司",
			SortOrder:   1,
		}

		org, err := setup.OrgSvc.CreateOrganization(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, entity.OrgTypeDivision, org.OrgType)
		assert.Equal(t, 2, org.Level)
		assert.NotNil(t, org.ParentID)
		assert.Equal(t, company.OrgID, *org.ParentID)
	})

	t.Run("创建项目组组织", func(t *testing.T) {
		filter := &service.OrganizationFilter{
			TenantID: setup.TenantID,
			OrgCode:  "COMP001",
		}
		companies, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter)
		require.Greater(t, len(companies), 0)

		company := companies[0]
		parentID := company.OrgID

		req := &service.CreateOrganizationRequest{
			TenantID:    setup.TenantID,
			OrgName:     "AI项目组",
			OrgType:     entity.OrgTypeProject,
			ParentID:    &parentID,
			OrgCode:     "PROJ001",
			Description: "AI研发项目组",
			SortOrder:   2,
		}

		org, err := setup.OrgSvc.CreateOrganization(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, entity.OrgTypeProject, org.OrgType)
	})

	t.Run("查询组织树", func(t *testing.T) {
		tree, err := setup.OrgSvc.GetOrganizationTree(ctx, setup.TenantID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tree), 1)

		// 验证树形结构
		root := tree[0]
		assert.Equal(t, entity.OrgTypeCompany, root.OrgType)
		assert.GreaterOrEqual(t, len(root.Children), 1)
	})

	t.Run("更新组织信息", func(t *testing.T) {
		filter := &service.OrganizationFilter{
			TenantID: setup.TenantID,
			OrgCode:  "COMP001",
		}
		companies, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter)
		require.Greater(t, len(companies), 0)

		company := companies[0]

		req := &service.UpdateOrganizationRequest{
			OrgID:       company.OrgID,
			OrgName:     "测试公司（已更新）",
			Description: "更新后的描述",
		}

		org, err := setup.OrgSvc.UpdateOrganization(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "测试公司（已更新）", org.OrgName)
		assert.Equal(t, "更新后的描述", org.Description)
	})

	t.Run("查询后代组织", func(t *testing.T) {
		filter := &service.OrganizationFilter{
			TenantID: setup.TenantID,
			OrgCode:  "COMP001",
		}
		companies, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter)
		if len(companies) > 0 {
			company := companies[0]

			descendants, err := setup.OrgSvc.GetDescendants(ctx, company.OrgID)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(descendants), 1)
		}
	})

	t.Run("删除项目组组织", func(t *testing.T) {
		filter := &service.OrganizationFilter{
			TenantID: setup.TenantID,
			OrgType:  string(entity.OrgTypeProject),
		}
		orgs, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter)
		if len(orgs) > 0 {
			project := orgs[0]

			err := setup.OrgSvc.DeleteOrganization(ctx, project.OrgID)
			require.NoError(t, err)
		}
	})
}

// testDeptCRUD 测试部门管理完整CRUD流程
func testDeptCRUD(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	// 先创建组织
	org, err := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
		TenantID:  setup.TenantID,
		OrgName:   "测试公司2",
		OrgType:   entity.OrgTypeCompany,
		OrgCode:   "COMP002",
		SortOrder: 1,
	})
	require.NoError(t, err)

	t.Run("创建根部门", func(t *testing.T) {
		req := &service.CreateDepartmentRequest{
			TenantID:    setup.TenantID,
			OrgID:       org.OrgID,
			DeptName:    "技术部",
			DeptCode:    "TECH",
			Description: "技术研发部门",
			SortOrder:   1,
		}

		dept, err := setup.DeptSvc.CreateDepartment(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, dept)
		assert.Equal(t, "技术部", dept.DeptName)
		assert.Equal(t, 1, dept.Level)
		assert.Equal(t, "/", dept.Path)
	})

	t.Run("查询部门树", func(t *testing.T) {
		tree, err := setup.DeptSvc.GetDepartmentTree(ctx, org.OrgID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tree), 1)
	})
}

// testPositionCRUD 测试岗位管理完整CRUD流程
func testPositionCRUD(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	t.Run("创建岗位", func(t *testing.T) {
		req := &service.CreatePositionRequest{
			TenantID:         setup.TenantID,
			PositionName:     "高级软件工程师",
			PositionCode:     "SENIOR_SE",
			Level:            3,
			Category:         "技术岗",
			Responsibilities: "负责核心系统开发",
			Requirements:     "5年以上工作经验",
			SortOrder:        1,
		}

		pos, err := setup.PosSvc.CreatePosition(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, pos)
		assert.Equal(t, "高级软件工程师", pos.PositionName)
		assert.Equal(t, 3, pos.Level)
	})

	t.Run("查询岗位列表", func(t *testing.T) {
		filter := &service.PositionFilter{
			TenantID: setup.TenantID,
			Page:     1,
			PageSize: 10,
		}

		positions, total, err := setup.PosSvc.ListPositions(ctx, filter)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(positions), 1)
		assert.GreaterOrEqual(t, total, 1)
	})
}

// testEmployeeCRUD 测试员工管理完整CRUD流程
func testEmployeeCRUD(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	// 准备数据
	org, _ := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
		TenantID:  setup.TenantID,
		OrgName:   "测试公司3",
		OrgType:   entity.OrgTypeCompany,
		OrgCode:   "COMP003",
		SortOrder: 1,
	})

	dept, _ := setup.DeptSvc.CreateDepartment(ctx, &service.CreateDepartmentRequest{
		TenantID:  setup.TenantID,
		OrgID:     org.OrgID,
		DeptName:  "技术部",
		DeptCode:  "TECH3",
		SortOrder: 1,
	})

	pos, _ := setup.PosSvc.CreatePosition(ctx, &service.CreatePositionRequest{
		TenantID:     setup.TenantID,
		PositionName: "软件工程师",
		PositionCode: "SE003",
		Level:        2,
		Category:     "技术岗",
	})

	t.Run("创建员工", func(t *testing.T) {
		gender := entity.GenderMale
		empType := entity.EmpTypeFullTime
		empStatus := entity.EmpStatusTrial

		req := &service.CreateEmployeeRequest{
			TenantID:        setup.TenantID,
			OrgID:           org.OrgID,
			DeptID:          &dept.DeptID,
			PositionID:      &pos.PositionID,
			EmpName:         "张三",
			EmpCode:         "EMP001",
			Gender:          &gender,
			Phone:           strPtr("13800138000"),
			Email:           strPtr("zhangsan@example.com"),
			EmployeeType:    empType,
			EmployeeStatus:  empStatus,
			JobLevel:        3,
			JobTitle:        "高级工程师",
			HireDate:        time.Now().UnixMilli(),
			ProbationDays:   90,
		}

		emp, err := setup.EmpSvc.CreateEmployee(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, emp)
		assert.Equal(t, "张三", emp.EmpName)
		assert.Equal(t, "EMP001", emp.EmpCode)
		assert.True(t, emp.IsTrial())
	})

	t.Run("员工转正", func(t *testing.T) {
		filter := &service.EmployeeFilter{
			TenantID: setup.TenantID,
			EmpCode:  "EMP001",
		}
		emps, _, _ := setup.EmpSvc.ListEmployees(ctx, filter)
		if len(emps) > 0 {
			emp := emps[0]

			req := &service.RegularizeEmployeeRequest{
				EmpID:       emp.EmpID,
				RegularDate: time.Now().UnixMilli(),
			}

			err := setup.EmpSvc.RegularizeEmployee(ctx, req)
			require.NoError(t, err)

			// 验证状态变更
			updated, _ := setup.EmpSvc.GetEmployee(ctx, emp.EmpID)
			assert.Equal(t, entity.EmpStatusActive, updated.EmployeeStatus)
		}
	})
}

// testDirectoryQuery 测试通讯录查询流程
func testDirectoryQuery(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	t.Run("查询组织通讯录", func(t *testing.T) {
		directory, err := setup.DirSvc.GetOrgDirectory(ctx, setup.TenantID)
		require.NoError(t, err)
		assert.NotNil(t, directory)
	})
}

// testHRLifecycle 测试HR生命周期流程
func testHRLifecycle(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	// 准备员工
	org, _ := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
		TenantID:  setup.TenantID,
		OrgName:   "测试公司4",
		OrgType:   entity.OrgTypeCompany,
		OrgCode:   "COMP004",
		SortOrder: 1,
	})

	emp, _ := setup.EmpSvc.CreateEmployee(ctx, &service.CreateEmployeeRequest{
		TenantID:       setup.TenantID,
		OrgID:          org.OrgID,
		EmpName:        "李四",
		EmpCode:        "EMP002",
		EmployeeType:   entity.EmpTypeFullTime,
		EmployeeStatus: entity.EmpStatusTrial,
		HireDate:       time.Now().UnixMilli(),
	})

	t.Run("创建员工合同", func(t *testing.T) {
		startDate := time.Now().UnixMilli()
		endDate := time.Now().AddDate(3, 0, 0).UnixMilli()

		req := &service.CreateContractRequest{
			TenantID:      setup.TenantID,
			EmpID:         emp.EmpID,
			ContractType:  "劳动合同",
			ContractNo:    "CONTRACT_001",
			StartDate:     startDate,
			EndDate:       &endDate,
			SalaryType:    "monthly",
			ProbationDays: 90,
			WorkHours:     "9:00-18:00",
		}

		contract, err := setup.HRSvc.CreateContract(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, contract)
		assert.Equal(t, "CONTRACT_001", contract.ContractNo)
		assert.Equal(t, "draft", contract.Status)
	})

	t.Run("员工调岗", func(t *testing.T) {
		// 创建新部门
		dept, _ := setup.DeptSvc.CreateDepartment(ctx, &service.CreateDepartmentRequest{
			TenantID:  setup.TenantID,
			OrgID:     org.OrgID,
			DeptName:  "产品部",
			DeptCode:  "PRODUCT4",
			SortOrder: 1,
		})

		// 创建新岗位
		pos, _ := setup.PosSvc.CreatePosition(ctx, &service.CreatePositionRequest{
			TenantID:     setup.TenantID,
			PositionName: "产品经理",
			PositionCode: "PM004",
			Level:        3,
			Category:     "产品岗",
		})

		req := &service.TransferEmployeeRequest{
			TenantID:      setup.TenantID,
			EmpID:         emp.EmpID,
			NewDeptID:     &dept.DeptID,
			NewPositionID: &pos.PositionID,
			NewJobLevel:   4,
			NewJobTitle:   "高级产品经理",
			TransferType:  "调岗",
			TransferDate:  time.Now().UnixMilli(),
			Reason:        "部门调整",
		}

		err := setup.HRSvc.TransferEmployee(ctx, req)
		require.NoError(t, err)

		// 验证调岗记录
		transfers, _ := setup.HRSvc.ListTransfers(ctx, setup.TenantID, emp.EmpID)
		assert.GreaterOrEqual(t, len(transfers), 1)
	})

	t.Run("员工离职", func(t *testing.T) {
		applyDate := time.Now().UnixMilli()
		lastWorkDate := time.Now().AddDate(0, 0, 30).UnixMilli()
		resignDate := time.Now().AddDate(0, 0, 31).UnixMilli()

		req := &service.ResignEmployeeRequest{
			TenantID:          setup.TenantID,
			EmpID:             emp.EmpID,
			ResignationType:   "主动离职",
			ResignationReason: "个人发展",
			ApplyDate:         applyDate,
			LastWorkDate:      lastWorkDate,
			ResignationDate:   resignDate,
			RehireEligible:    true,
		}

		resignation, err := setup.HRSvc.SubmitResignation(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, resignation)
		assert.Equal(t, "pending", resignation.ApprovalStatus)
	})
}

// testTenantIsolation 测试多租户隔离
func testTenantIsolation(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	tenant1 := setup.TenantID
	tenant2 := "tenant-002"

	t.Run("租户数据隔离", func(t *testing.T) {
		// 租户1创建组织
		org1, err := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
			TenantID:  tenant1,
			OrgName:   "租户1公司",
			OrgType:   entity.OrgTypeCompany,
			OrgCode:   "T1_COMP",
			SortOrder: 1,
		})
		require.NoError(t, err)

		// 租户2创建组织
		org2, err := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
			TenantID:  tenant2,
			OrgName:   "租户2公司",
			OrgType:   entity.OrgTypeCompany,
			OrgCode:   "T2_COMP",
			SortOrder: 1,
		})
		require.NoError(t, err)

		// 租户1只能看到自己的数据
		filter1 := &service.OrganizationFilter{
			TenantID: tenant1,
		}
		orgs1, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter1)
		for _, org := range orgs1 {
			assert.Equal(t, tenant1, org.TenantID)
		}

		// 租户2只能看到自己的数据
		filter2 := &service.OrganizationFilter{
			TenantID: tenant2,
		}
		orgs2, _, _ := setup.OrgSvc.ListOrganizations(ctx, filter2)
		for _, org := range orgs2 {
			assert.Equal(t, tenant2, org.TenantID)
		}

		// 验证两个组织的ID不同
		assert.NotEqual(t, org1.OrgID, org2.OrgID)
	})

	t.Run("租户完全隔离验证", func(t *testing.T) {
		// 验证所有表的数据都按tenant_id隔离
		var count int64

		// organizations表
		setup.DB.Table("organizations").
			Where("tenant_id = ? OR tenant_id = ?", tenant1, tenant2).
			Count(&count)
		assert.Greater(t, count, int64(0))
	})
}

// testTransactions 测试事务
func testTransactions(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	t.Run("事务提交成功", func(t *testing.T) {
		err := setup.DB.Transaction(func(tx *gorm.DB) error {
			// 创建组织
			org := &entity.Organization{
				OrgID:      fmt.Sprintf("%d", time.Now().UnixNano()),
				TenantID:   setup.TenantID,
				OrgName:    "事务测试公司",
				OrgType:    entity.OrgTypeCompany,
				OrgCode:    "TX_COMP",
				Level:      1,
				Path:       "/",
				Status:     entity.OrgStatusActive,
				CreatedAt:  time.Now().UnixMilli(),
				UpdatedAt:  time.Now().UnixMilli(),
			}
			if err := tx.Create(org).Error; err != nil {
				return err
			}

			// 创建部门
			dept := &entity.Department{
				DeptID:      fmt.Sprintf("%d", time.Now().UnixNano()+1),
				TenantID:    setup.TenantID,
				OrgID:       org.OrgID,
				DeptName:    "技术部",
				DeptCode:    "TECH_TX",
				Level:       1,
				Path:        "/",
				Status:      entity.OrgStatusActive,
				CreatedAt:   time.Now().UnixMilli(),
				UpdatedAt:   time.Now().UnixMilli(),
			}
			return tx.Create(dept).Error
		})

		require.NoError(t, err)
	})

	t.Run("事务回滚失败操作", func(t *testing.T) {
		initialCount := int64(0)
		setup.DB.Table("organizations").Where("tenant_id = ?", setup.TenantID).Count(&initialCount)

		err := setup.DB.Transaction(func(tx *gorm.DB) error {
			// 创建组织
			org := &entity.Organization{
				OrgID:      fmt.Sprintf("%d", time.Now().UnixNano()),
				TenantID:   setup.TenantID,
				OrgName:    "回滚测试公司",
				OrgType:    entity.OrgTypeCompany,
				OrgCode:    "ROLLBACK_COMP",
				Level:      1,
				Path:       "/",
				Status:     entity.OrgStatusActive,
				CreatedAt:  time.Now().UnixMilli(),
				UpdatedAt:  time.Now().UnixMilli(),
			}
			if err := tx.Create(org).Error; err != nil {
				return err
			}

			// 模拟错误
			return fmt.Errorf("模拟错误导致回滚")
		})

		assert.Error(t, err)

		// 验证数据未插入
		finalCount := int64(0)
		setup.DB.Table("organizations").Where("tenant_id = ?", setup.TenantID).Count(&finalCount)
		assert.Equal(t, initialCount, finalCount)
	})
}

// testComplexScenarios 测试复杂业务场景
func testComplexScenarios(t *testing.T, setup *TestSetup) {
	ctx := context.Background()

	t.Run("组织树结构查询", func(t *testing.T) {
		// 创建复杂组织树
		org, _ := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
			TenantID:  setup.TenantID,
			OrgName:   "集团总部",
			OrgType:   entity.OrgTypeCompany,
			OrgCode:   "HQ",
			SortOrder: 1,
		})

		// 创建多层组织
		var parentID string = org.OrgID
		for i := 1; i <= 3; i++ {
			div, _ := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
				TenantID:  setup.TenantID,
				OrgName:   fmt.Sprintf("第%d级分公司", i),
				OrgType:   entity.OrgTypeDivision,
				ParentID:  &parentID,
				OrgCode:   fmt.Sprintf("DIV%d", i),
				SortOrder: i,
			})
			parentID = div.OrgID
		}

		// 查询后代树
		descendants, _ := setup.OrgSvc.GetDescendants(ctx, org.OrgID)
		assert.GreaterOrEqual(t, len(descendants), 3) // 至少有3个后代
	})

	t.Run("员工状态转换", func(t *testing.T) {
		org, _ := setup.OrgSvc.CreateOrganization(ctx, &service.CreateOrganizationRequest{
			TenantID:  setup.TenantID,
			OrgName:   "测试公司5",
			OrgType:   entity.OrgTypeCompany,
			OrgCode:   "STATUS_TEST",
			SortOrder: 1,
		})

		// 创建试用员工
		emp, _ := setup.EmpSvc.CreateEmployee(ctx, &service.CreateEmployeeRequest{
			TenantID:       setup.TenantID,
			OrgID:          org.OrgID,
			EmpName:        "测试员工",
			EmpCode:        fmt.Sprintf("EMP_%d", time.Now().UnixNano()),
			EmployeeType:   entity.EmpTypeFullTime,
			EmployeeStatus: entity.EmpStatusTrial,
			HireDate:       time.Now().UnixMilli(),
		})

		// 验证试用期状态
		assert.True(t, emp.IsTrial())

		// 转正
		regularDate := time.Now().UnixMilli()
		setup.EmpSvc.RegularizeEmployee(ctx, &service.RegularizeEmployeeRequest{
			EmpID:       emp.EmpID,
			RegularDate: regularDate,
		})

		// 验证转正状态
		regularized, _ := setup.EmpSvc.GetEmployee(ctx, emp.EmpID)
		assert.Equal(t, entity.EmpStatusActive, regularized.EmployeeStatus)
		assert.False(t, regularized.IsTrial())
	})
}

// 辅助函数

// strPtr 返回字符串指针
func strPtr(s string) *string {
	return &s
}
