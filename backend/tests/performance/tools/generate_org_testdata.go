/*
 * 组织中心性能测试 - 测试数据生成工具
 *
 * 功能: 生成大量模拟组织数据用于性能测试
 *
 * 生成规模:
 *   - 租户: 100个
 *   - 每个租户: 1000个组织
 *   - 每个组织: 100个部门
 *   - 每个部门: 50个岗位
 *   - 每个部门: 1000个员工
 *   总计: 10,000租户, 100,000组织, 10,000,000部门, 500,000,000岗位, 10,000,000,000员工
 *
 * 使用方法:
 *   # 生成小规模测试数据 (1个租户, 10个组织, 100个部门, 1000个员工)
 *   go run generate_org_testdata.go -size=small
 *
 *   # 生成中等规模测试数据 (10个租户, 100个组织, 1000个部门, 10000个员工)
 *   go run generate_org_testdata.go -size=medium
 *
 *   # 生成大规模测试数据 (100个租户, 1000个组织, 10000个部门, 100000个员工)
 *   go run generate_org_testdata.go -size=large
 *
 *   # 自定义规模
 *   go run generate_org_testdata.go -tenants=10 -orgs=100 -depts=1000 -emps=10000
 *
 * @author 研发B (后端工程师)
 * @version 1.0
 * @date 2025-01-01
 */

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// ================================
// 配置参数
// ================================

type Config struct {
	Tenants         int    // 租户数量
	OrgsPerTenant   int    // 每个租户的组织数量
	DeptsPerOrg     int    // 每个组织的部门数量
	PositionsPerDept int   // 每个部门的岗位数量
	EmpsPerDept     int    // 每个部门的员工数量
	BatchSize       int    // 批量插入大小
	DSN             string // 数据库连接字符串
}

var (
	// 命令行参数
	size      = flag.String("size", "small", "测试数据规模: small, medium, large, custom")
	tenants   = flag.Int("tenants", 0, "租户数量 (自定义模式)")
	orgs      = flag.Int("orgs", 0, "每个租户的组织数量 (自定义模式)")
	depts     = flag.Int("depts", 0, "每个组织的部门数量 (自定义模式)")
	positions = flag.Int("positions", 0, "每个部门的岗位数量 (自定义模式)")
	emps      = flag.Int("emps", 0, "每个部门的员工数量 (自定义模式)")
	batch     = flag.Int("batch", 1000, "批量插入大小")
	dsn       = flag.String("dsn", "root:password@tcp(localhost:3306)/zker_prod?charset=utf8mb4&parseTime=True&loc=Local", "数据库连接字符串")
)

// ================================
// 预定义规模
// ================================

var predefinedSizes = map[string]Config{
	"small": {
		Tenants:         1,
		OrgsPerTenant:   10,
		DeptsPerOrg:     10,
		PositionsPerDept: 10,
		EmpsPerDept:     10,
		BatchSize:       100,
	},
	"medium": {
		Tenants:         10,
		OrgsPerTenant:   100,
		DeptsPerOrg:     100,
		PositionsPerDept: 50,
		EmpsPerDept:     100,
		BatchSize:       1000,
	},
	"large": {
		Tenants:         100,
		OrgsPerTenant:   1000,
		DeptsPerOrg:     100,
		PositionsPerDept: 50,
		EmpsPerDept:     1000,
		BatchSize:       10000,
	},
}

// ================================
// 辅助函数
// ================================

// 生成随机UUID
func generateUUID() string {
	return uuid.New().String()
}

// 当前时间戳（毫秒）
func currentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// 随机选择
func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

// 随机手机号
func randomPhone() string {
	return fmt.Sprintf("138%08d", rand.Intn(100000000))
}

// 随机邮箱
func randomEmail(name string) string {
	domains := []string{"qq.com", "163.com", "gmail.com", "outlook.com", "company.com"}
	return fmt.Sprintf("%s@%s", name, domains[rand.Intn(len(domains))])
}

// ================================
// 数据生成器
// ================================

// DataGenerator 数据生成器
type DataGenerator struct {
	config Config
	db     *sql.DB
	wg     sync.WaitGroup
}

// NewDataGenerator 创建数据生成器
func NewDataGenerator(config Config) (*DataGenerator, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库ping失败: %w", err)
	}

	return &DataGenerator{
		config: config,
		db:     db,
	}, nil
}

// Close 关闭连接
func (dg *DataGenerator) Close() error {
	return dg.db.Close()
}

// GenerateAll 生成所有测试数据
func (dg *DataGenerator) GenerateAll() error {
	log.Println("🚀 开始生成测试数据...")
	log.Printf("📊 配置: 租户=%d, 组织/租户=%d, 部门/组织=%d, 岗位/部门=%d, 员工/部门=%d",
		dg.config.Tenants,
		dg.config.OrgsPerTenant,
		dg.config.DeptsPerOrg,
		dg.config.PositionsPerDept,
		dg.config.EmpsPerDept,
	)

	startTime := time.Now()

	// 1. 生成租户
	log.Println("📝 步骤1: 生成租户...")
	tenantIDs, err := dg.generateTenants()
	if err != nil {
		return fmt.Errorf("生成租户失败: %w", err)
	}
	log.Printf("✅ 生成 %d 个租户", len(tenantIDs))

	// 2. 生成组织
	log.Println("📝 步骤2: 生成组织...")
	orgIDs, err := dg.generateOrganizations(tenantIDs)
	if err != nil {
		return fmt.Errorf("生成组织失败: %w", err)
	}
	log.Printf("✅ 生成 %d 个组织", len(orgIDs))

	// 3. 生成部门
	log.Println("📝 步骤3: 生成部门...")
	deptIDs, err := dg.generateDepartments(orgIDs)
	if err != nil {
		return fmt.Errorf("生成部门失败: %w", err)
	}
	log.Printf("✅ 生成 %d 个部门", len(deptIDs))

	// 4. 生成岗位
	log.Println("📝 步骤4: 生成岗位...")
	positionCount, err := dg.generatePositions(deptIDs)
	if err != nil {
		return fmt.Errorf("生成岗位失败: %w", err)
	}
	log.Printf("✅ 生成 %d 个岗位", positionCount)

	// 5. 生成员工
	log.Println("📝 步骤5: 生成员工...")
	empCount, err := dg.generateEmployees(deptIDs)
	if err != nil {
		return fmt.Errorf("生成员工失败: %w", err)
	}
	log.Printf("✅ 生成 %d 个员工", empCount)

	elapsed := time.Since(startTime)
	log.Printf("✅ 测试数据生成完成! 总耗时: %v", elapsed)

	return nil
}

// generateTenants 生成租户
func (dg *DataGenerator) generateTenants() ([]string, error) {
	var tenantIDs []string
	var mu sync.Mutex

	// 并发生成
	sem := make(chan struct{}, 10) // 最多10个并发
	var wg sync.WaitGroup

	for i := 0; i < dg.config.Tenants; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(index int) {
			defer wg.Done()
			defer func() { <-sem }()

			tenantID := generateUUID()
			tenantName := fmt.Sprintf("测试租户_%04d", index)
			adminEmail := fmt.Sprintf("admin%04d@test.com", index)

			query := `
				INSERT INTO tenants (tenant_id, tenant_name, tenant_type, admin_email, status, created_at, updated_at)
				VALUES (?, ?, 'enterprise', ?, 'active', ?, ?)
			`

			_, err := dg.db.Exec(query, tenantID, tenantName, adminEmail, currentTimestamp(), currentTimestamp())
			if err != nil {
				log.Printf("❌ 插入租户失败: %v", err)
				return
			}

			mu.Lock()
			tenantIDs = append(tenantIDs, tenantID)
			mu.Unlock()

			if (index+1)%100 == 0 {
				log.Printf("  已生成 %d/%d 租户", index+1, dg.config.Tenants)
			}
		}(i)
	}

	wg.Wait()
	return tenantIDs, nil
}

// generateOrganizations 生成组织
func (dg *DataGenerator) generateOrganizations(tenantIDs []string) (map[string][]string, error) {
	orgMap := make(map[string][]string) // tenant_id -> []org_id
	var mu sync.Mutex

	// 并发生成
	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup

	totalOrgs := 0
	for _, tenantID := range tenantIDs {
		for i := 0; i < dg.config.OrgsPerTenant; i++ {
			totalOrgs++
			wg.Add(1)
			sem <- struct{}{}

			go func(tenantID string, index int) {
				defer wg.Done()
				defer func() { <-sem }()

				orgID := generateUUID()
				orgName := fmt.Sprintf("测试组织_%s_%04d", tenantID[:8], index)
				orgCode := fmt.Sprintf("ORG_%s_%04d", tenantID[:8], index)

				// 随机组织类型
				orgTypes := []string{"company", "division", "department", "project"}
				orgType := orgTypes[rand.Intn(len(orgTypes))]

				query := `
					INSERT INTO organizations (org_id, tenant_id, org_name, org_type, org_code, level, path, status, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, 1, ?, 'active', ?, ?)
				`

				_, err := dg.db.Exec(query, orgID, tenantID, orgName, orgType, orgCode, "/"+orgID, currentTimestamp(), currentTimestamp())
				if err != nil {
					log.Printf("❌ 插入组织失败: %v", err)
					return
				}

				// 插入闭包表记录（自身）
				treeQuery := `
					INSERT INTO organization_trees (tenant_id, ancestor_id, descendant_id, depth)
					VALUES (?, ?, ?, 0)
				`
				dg.db.Exec(treeQuery, tenantID, orgID, orgID)

				mu.Lock()
				orgMap[tenantID] = append(orgMap[tenantID], orgID)
				mu.Unlock()

			}(tenantID, i)
		}
	}

	wg.Wait()
	log.Printf("  总计生成 %d 个组织", totalOrgs)
	return orgMap, nil
}

// generateDepartments 生成部门
func (dg *DataGenerator) generateDepartments(orgMap map[string][]string) (map[string][]string, error) {
	deptMap := make(map[string][]string) // org_id -> []dept_id
	var mu sync.Mutex

	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup

	totalDepts := 0
	for tenantID, orgIDs := range orgMap {
		for _, orgID := range orgIDs {
			for i := 0; i < dg.config.DeptsPerOrg; i++ {
				totalDepts++
				wg.Add(1)
				sem <- struct{}{}

				go func(tenantID, orgID string, index int) {
					defer wg.Done()
					defer func() { <-sem }()

					deptID := generateUUID()
					deptName := fmt.Sprintf("测试部门_%s_%04d", orgID[:8], index)
					deptCode := fmt.Sprintf("DEPT_%s_%04d", orgID[:8], index)

					query := `
						INSERT INTO departments (dept_id, tenant_id, org_id, dept_name, dept_code, level, path, status, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, 1, ?, 'active', ?, ?)
					`

					_, err := dg.db.Exec(query, deptID, tenantID, orgID, deptName, deptCode, "/"+deptID, currentTimestamp(), currentTimestamp())
					if err != nil {
						log.Printf("❌ 插入部门失败: %v", err)
						return
					}

					// 插入闭包表记录（自身）
					treeQuery := `
						INSERT INTO department_trees (tenant_id, ancestor_id, descendant_id, depth)
						VALUES (?, ?, ?, 0)
					`
					dg.db.Exec(treeQuery, tenantID, deptID, deptID)

					mu.Lock()
					deptMap[orgID] = append(deptMap[orgID], deptID)
					mu.Unlock()

				}(tenantID, orgID, i)
			}
		}
	}

	wg.Wait()
	log.Printf("  总计生成 %d 个部门", totalDepts)
	return deptMap, nil
}

// generatePositions 生成岗位
func (dg *DataGenerator) generatePositions(deptMap map[string][]string) (int, error) {
	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup
	totalPositions := 0

	for orgID, deptIDs := range deptMap {
		tenantID := orgID[:36] // 从orgID提取tenantID（简化处理）

		for _, deptID := range deptIDs {
			for i := 0; i < dg.config.PositionsPerDept; i++ {
				totalPositions++
				wg.Add(1)
				sem <- struct{}{}

				go func(tenantID, deptID string, index int) {
					defer wg.Done()
					defer func() { <-sem }()

					positionID := generateUUID()
					positionName := fmt.Sprintf("岗位_%04d", index)
					positionCode := fmt.Sprintf("POS_%s_%04d", deptID[:8], index)

					// 随机职级和类别
					level := rand.Intn(10) + 1
					categories := []string{"技术岗", "管理岗", "职能岗", "销售岗", "运营岗"}
					category := categories[rand.Intn(len(categories))]

					query := `
						INSERT INTO positions (position_id, tenant_id, dept_id, position_name, position_code, level, category, status, sort_order, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, 'active', ?, ?, ?)
					`

					_, err := dg.db.Exec(query, positionID, tenantID, deptID, positionName, positionCode, level, category, index, currentTimestamp(), currentTimestamp())
					if err != nil {
						log.Printf("❌ 插入岗位失败: %v", err)
						return
					}

				}(tenantID, deptID, i)
			}
		}
	}

	wg.Wait()
	return totalPositions, nil
}

// generateEmployees 生成员工
func (dg *DataGenerator) generateEmployees(deptMap map[string][]string) (int, error) {
	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup
	totalEmps := 0

	// 姓氏和名字库
	surnames := []string{"张", "李", "王", "刘", "陈", "杨", "赵", "黄", "周", "吴"}
	names := []string{"伟", "芳", "娜", "敏", "静", "丽", "强", "磊", "军", "洋", "勇", "艳", "杰", "涛", "明"}

	for orgID, deptIDs := range deptMap {
		tenantID := orgID[:36]

		for _, deptID := range deptIDs {
			for i := 0; i < dg.config.EmpsPerDept; i++ {
				totalEmps++
				wg.Add(1)
				sem <- struct{}{}

				go func(tenantID, orgID, deptID string, index int) {
					defer wg.Done()
					defer func() { <-sem }()

					empID := generateUUID()
					surname := surnames[rand.Intn(len(surnames))]
					name := names[rand.Intn(len(names))]
					empName := surname + name
					empCode := fmt.Sprintf("EMP%s%04d", deptID[:8], index)

					// 随机员工类型和状态
					empTypes := []string{"full_time", "part_time", "intern", "outsourcing"}
					empStatuses := []string{"active", "trial", "probation"}
					empType := empTypes[rand.Intn(len(empTypes))]
					empStatus := empStatuses[rand.Intn(len(empStatuses))]

					// 随机入职日期（过去3年内）
					hireDate := time.Now().Add(-time.Duration(rand.Intn(1095)) * 24 * time.Hour).UnixNano() / int64(time.Millisecond)

					query := `
						INSERT INTO employees (emp_id, tenant_id, org_id, dept_id, emp_name, emp_code, employee_type, status, hire_date, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
					`

					_, err := dg.db.Exec(query, empID, tenantID, orgID, deptID, empName, empCode, empType, empStatus, hireDate, currentTimestamp(), currentTimestamp())
					if err != nil {
						log.Printf("❌ 插入员工失败: %v", err)
						return
					}

				}(tenantID, orgID, deptID, i)
			}
		}
	}

	wg.Wait()
	return totalEmps, nil
}

// ================================
// Main 函数
// ================================

func main() {
	flag.Parse()

	// 确定配置
	var config Config
	switch *size {
	case "small", "medium", "large":
		config = predefinedSizes[*size]
	case "custom":
		if *tenants == 0 || *orgs == 0 || *depts == 0 || *emps == 0 {
			log.Fatal("❌ 自定义模式需要指定 -tenants, -orgs, -depts, -emps 参数")
		}
		config = Config{
			Tenants:         *tenants,
			OrgsPerTenant:   *orgs,
			DeptsPerOrg:     *depts,
			PositionsPerDept: *positions,
			EmpsPerDept:     *emps,
			BatchSize:       *batch,
		}
	default:
		log.Fatalf("❌ 未知的规模: %s (支持: small, medium, large, custom)", *size)
	}

	config.DSN = *dsn

	// 创建生成器
	gen, err := NewDataGenerator(config)
	if err != nil {
		log.Fatalf("❌ 创建数据生成器失败: %v", err)
	}
	defer gen.Close()

	// 生成数据
	if err := gen.GenerateAll(); err != nil {
		log.Fatalf("❌ 生成测试数据失败: %v", err)
	}

	log.Println("✅ 所有测试数据生成完成!")
}
