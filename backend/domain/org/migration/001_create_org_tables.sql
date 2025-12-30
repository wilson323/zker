-- backend/domain/org/migration/001_create_org_tables.sql
-- 组织中心表结构创建脚本
-- 创建时间：2025-01-01
-- 描述：创建组织、部门、岗位、员工及相关表的完整结构

-- ================================
-- 1. 组织表（organizations）
-- ================================
CREATE TABLE IF NOT EXISTS organizations (
    org_id VARCHAR(36) PRIMARY KEY COMMENT '组织ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    org_name VARCHAR(200) NOT NULL COMMENT '组织名称',
    org_type ENUM('company', 'division', 'department', 'project') NOT NULL COMMENT '组织类型：公司/分公司/部门（虚拟）/项目组',
    parent_id VARCHAR(36) DEFAULT NULL COMMENT '父组织ID',
    org_code VARCHAR(50) NOT NULL COMMENT '组织编码',
    level INT NOT NULL DEFAULT 1 COMMENT '层级：1-顶级，2-二级...',
    path VARCHAR(500) NOT NULL COMMENT '组织路径：/1/2/3',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序',
    status ENUM('active', 'inactive', 'frozen') NOT NULL DEFAULT 'active' COMMENT '状态：激活/停用/冻结',
    description TEXT COMMENT '描述',
    leader_id VARCHAR(36) DEFAULT NULL COMMENT '负责人ID',
    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    deleted_at BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_level (level),
    INDEX idx_status (status),
    INDEX idx_sort (sort_order),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX uk_tenant_code (tenant_id, org_code),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES organizations(org_id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (leader_id) REFERENCES employees(emp_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';

-- ================================
-- 2. 组织树闭包表（organization_trees）
-- ================================
CREATE TABLE IF NOT EXISTS organization_trees (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    ancestor_id VARCHAR(36) NOT NULL COMMENT '祖先节点ID',
    descendant_id VARCHAR(36) NOT NULL COMMENT '后代节点ID',
    depth INT NOT NULL COMMENT '层级深度：0表示自己，1表示子节点，2表示孙节点',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_ancestor (ancestor_id),
    INDEX idx_descendant (descendant_id),
    INDEX idx_ancestor_depth (ancestor_id, depth),

    -- 外键约束
    FOREIGN KEY (ancestor_id) REFERENCES organizations(org_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES organizations(org_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织树闭包表';

-- ================================
-- 3. 部门表（departments）
-- ================================
CREATE TABLE IF NOT EXISTS departments (
    dept_id VARCHAR(36) PRIMARY KEY COMMENT '部门ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    org_id VARCHAR(36) NOT NULL COMMENT '所属组织ID',
    parent_id VARCHAR(36) DEFAULT NULL COMMENT '父部门ID',
    dept_name VARCHAR(200) NOT NULL COMMENT '部门名称',
    dept_code VARCHAR(50) NOT NULL COMMENT '部门编码',
    level INT NOT NULL DEFAULT 1 COMMENT '层级：1-顶级，2-二级...',
    path VARCHAR(500) NOT NULL COMMENT '部门路径：/1/2/3',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序',
    status ENUM('active', 'inactive', 'frozen') NOT NULL DEFAULT 'active' COMMENT '状态',
    leader_id VARCHAR(36) DEFAULT NULL COMMENT '负责人ID',
    parent_leader VARCHAR(36) DEFAULT NULL COMMENT '上级部门负责人',
    description TEXT COMMENT '描述',
    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    deleted_at BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_org_id (org_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_level (level),
    INDEX idx_status (status),
    INDEX idx_sort (sort_order),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX uk_tenant_code (tenant_id, dept_code),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (org_id) REFERENCES organizations(org_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES departments(dept_id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (leader_id) REFERENCES employees(emp_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ================================
-- 4. 部门树闭包表（department_trees）
-- ================================
CREATE TABLE IF NOT EXISTS department_trees (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    ancestor_id VARCHAR(36) NOT NULL COMMENT '祖先节点ID',
    descendant_id VARCHAR(36) NOT NULL COMMENT '后代节点ID',
    depth INT NOT NULL COMMENT '层级深度：0表示自己，1表示子节点，2表示孙节点',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_ancestor (ancestor_id),
    INDEX idx_descendant (descendant_id),
    INDEX idx_ancestor_depth (ancestor_id, depth),

    -- 外键约束
    FOREIGN KEY (ancestor_id) REFERENCES departments(dept_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (descendant_id) REFERENCES departments(dept_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门树闭包表';

-- ================================
-- 5. 岗位表（positions）
-- ================================
CREATE TABLE IF NOT EXISTS positions (
    position_id VARCHAR(36) PRIMARY KEY COMMENT '岗位ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    dept_id VARCHAR(36) DEFAULT NULL COMMENT '所属部门ID',
    position_name VARCHAR(100) NOT NULL COMMENT '岗位名称',
    position_code VARCHAR(50) NOT NULL COMMENT '岗位编码',
    level INT NOT NULL DEFAULT 1 COMMENT '职级：1-初级，2-中级，3-高级',
    category VARCHAR(50) NOT NULL COMMENT '岗位类别：技术岗/管理岗/职能岗',
    responsibilities TEXT COMMENT '岗位职责',
    requirements TEXT COMMENT '任职要求',
    status ENUM('active', 'inactive', 'frozen') NOT NULL DEFAULT 'active' COMMENT '状态',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序',
    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    deleted_at BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_dept_id (dept_id),
    INDEX idx_level (level),
    INDEX idx_status (status),
    INDEX idx_sort (sort_order),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX uk_tenant_code (tenant_id, position_code),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (dept_id) REFERENCES departments(dept_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='岗位表';

-- ================================
-- 6. 员工表（employees）
-- ================================
CREATE TABLE IF NOT EXISTS employees (
    emp_id VARCHAR(36) PRIMARY KEY COMMENT '员工ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    user_id VARCHAR(36) DEFAULT NULL COMMENT '关联系统用户ID',
    org_id VARCHAR(36) NOT NULL COMMENT '组织ID',
    dept_id VARCHAR(36) DEFAULT NULL COMMENT '部门ID',
    position_id VARCHAR(36) DEFAULT NULL COMMENT '岗位ID',

    -- 基本信息
    emp_name VARCHAR(100) NOT NULL COMMENT '员工姓名',
    emp_code VARCHAR(50) NOT NULL COMMENT '工号',
    gender ENUM('male', 'female', 'other') DEFAULT NULL COMMENT '性别',
    phone VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    email VARCHAR(100) DEFAULT NULL COMMENT '邮箱',

    -- 职位信息
    employee_type ENUM('full_time', 'part_time', 'intern', 'outsourcing', 'contractor') NOT NULL DEFAULT 'full_time' COMMENT '员工类型',
    employee_status ENUM('active', 'trial', 'probation', 'resigned', 'retired', 'suspended') NOT NULL DEFAULT 'trial' COMMENT '员工状态',
    job_level INT NOT NULL DEFAULT 1 COMMENT '职级：P1-P10',
    job_title VARCHAR(100) DEFAULT NULL COMMENT '职衔',

    -- 入职信息
    hire_date BIGINT NOT NULL COMMENT '入职日期（毫秒时间戳）',
    regular_date BIGINT DEFAULT NULL COMMENT '转正日期（毫秒时间戳）',
    probation_days INT NOT NULL DEFAULT 90 COMMENT '试用天数',

    -- 工作信息
    work_location VARCHAR(200) DEFAULT NULL COMMENT '工作地点',
    direct_leader_id VARCHAR(36) DEFAULT NULL COMMENT '直接上级ID',

    -- 个人信息
    id_card VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
    birthday BIGINT DEFAULT NULL COMMENT '生日（毫秒时间戳）',
    address VARCHAR(500) DEFAULT NULL COMMENT '地址',
    education VARCHAR(50) DEFAULT NULL COMMENT '学历',
    graduate_school VARCHAR(200) DEFAULT NULL COMMENT '毕业院校',
    major VARCHAR(100) DEFAULT NULL COMMENT '专业',

    -- 紧急联系人
    emergency_contact VARCHAR(100) DEFAULT NULL COMMENT '紧急联系人',
    emergency_phone VARCHAR(20) DEFAULT NULL COMMENT '紧急联系电话',

    -- 其他
    status ENUM('active', 'trial', 'probation', 'resigned', 'retired', 'suspended') NOT NULL DEFAULT 'trial' COMMENT '状态',
    avatar_url VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
    description TEXT COMMENT '备注',

    -- 时间戳
    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    deleted_at BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_org_id (org_id),
    INDEX idx_dept_id (dept_id),
    INDEX idx_position_id (position_id),
    INDEX idx_emp_name (emp_name),
    INDEX idx_emp_code (emp_code),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_job_level (job_level),
    INDEX idx_hire_date (hire_date),
    INDEX idx_regular_date (regular_date),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX uk_tenant_code (tenant_id, emp_code),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (org_id) REFERENCES organizations(org_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (dept_id) REFERENCES departments(dept_id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (position_id) REFERENCES positions(position_id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (direct_leader_id) REFERENCES employees(emp_id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工表';

-- ================================
-- 7. 员工合同表（employee_contracts）
-- ================================
CREATE TABLE IF NOT EXISTS employee_contracts (
    contract_id VARCHAR(36) PRIMARY KEY COMMENT '合同ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    emp_id VARCHAR(36) NOT NULL COMMENT '员工ID',

    contract_type VARCHAR(50) NOT NULL COMMENT '合同类型：劳动合同/实习协议/外包协议',
    contract_no VARCHAR(100) NOT NULL COMMENT '合同编号',

    start_date BIGINT NOT NULL COMMENT '合同开始日期（毫秒时间戳）',
    end_date BIGINT DEFAULT NULL COMMENT '合同结束日期（无固定期限则为NULL）',

    salary VARCHAR(50) DEFAULT NULL COMMENT '薪资（密文）',
    salary_type VARCHAR(20) NOT NULL DEFAULT 'monthly' COMMENT '薪资类型：monthly/yearly',

    probation_days INT NOT NULL DEFAULT 0 COMMENT '试用天数',
    probation_salary VARCHAR(50) DEFAULT NULL COMMENT '试用期薪资',

    work_hours VARCHAR(50) NOT NULL DEFAULT '8:00-17:00' COMMENT '工作时间',
    work_place VARCHAR(200) DEFAULT NULL COMMENT '工作地点',

    contract_file_url VARCHAR(500) DEFAULT NULL COMMENT '合同文件URL',

    status ENUM('draft', 'active', 'expired', 'terminated') NOT NULL DEFAULT 'draft' COMMENT '状态：草稿/生效/过期/终止',

    signed_at BIGINT DEFAULT NULL COMMENT '签署时间（毫秒时间戳）',

    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    deleted_at BIGINT DEFAULT NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_emp_id (emp_id),
    INDEX idx_start_date (start_date),
    INDEX idx_end_date (end_date),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX uk_tenant_no (tenant_id, contract_no),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (emp_id) REFERENCES employees(emp_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工合同表';

-- ================================
-- 8. 员工调岗记录表（employee_transfers）
-- ================================
CREATE TABLE IF NOT EXISTS employee_transfers (
    transfer_id VARCHAR(36) PRIMARY KEY COMMENT '调岗记录ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    emp_id VARCHAR(36) NOT NULL COMMENT '员工ID',

    -- 变更前
    old_dept_id VARCHAR(36) DEFAULT NULL COMMENT '原部门ID',
    old_position_id VARCHAR(36) DEFAULT NULL COMMENT '原岗位ID',
    old_job_level INT NOT NULL DEFAULT 1 COMMENT '原职级',
    old_job_title VARCHAR(100) DEFAULT NULL COMMENT '原职衔',

    -- 变更后
    new_dept_id VARCHAR(36) DEFAULT NULL COMMENT '新部门ID',
    new_position_id VARCHAR(36) DEFAULT NULL COMMENT '新岗位ID',
    new_job_level INT NOT NULL DEFAULT 1 COMMENT '新职级',
    new_job_title VARCHAR(100) DEFAULT NULL COMMENT '新职衔',

    -- 调岗信息
    transfer_type VARCHAR(50) NOT NULL COMMENT '调岗类型：调岗/晋升/降职/平调',
    transfer_date BIGINT NOT NULL COMMENT '调岗日期（毫秒时间戳）',
    reason TEXT COMMENT '调岗原因',

    approver_id VARCHAR(36) DEFAULT NULL COMMENT '审批人ID',
    approved_at BIGINT DEFAULT NULL COMMENT '审批时间（毫秒时间戳）',

    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_emp_id (emp_id),
    INDEX idx_transfer_date (transfer_date),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (emp_id) REFERENCES employees(emp_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工调岗记录表';

-- ================================
-- 9. 员工离职记录表（employee_resignations）
-- ================================
CREATE TABLE IF NOT EXISTS employee_resignations (
    resignation_id VARCHAR(36) PRIMARY KEY COMMENT '离职记录ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    emp_id VARCHAR(36) NOT NULL COMMENT '员工ID',

    resignation_type VARCHAR(50) NOT NULL COMMENT '离职类型：主动离职/被动离职/合同到期/退休',
    resignation_reason TEXT COMMENT '离职原因',

    -- 时间信息
    apply_date BIGINT NOT NULL COMMENT '申请日期（毫秒时间戳）',
    last_work_date BIGINT NOT NULL COMMENT '最后工作日（毫秒时间戳）',
    resignation_date BIGINT NOT NULL COMMENT '离职生效日期（毫秒时间戳）',

    -- 流程信息
    handover_to_id VARCHAR(36) DEFAULT NULL COMMENT '工作交接人',
    handover_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '交接状态：pending/completed',

    approval_status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '审批状态：pending/approved/rejected',
    approver_id VARCHAR(36) DEFAULT NULL COMMENT '审批人ID',
    approved_at BIGINT DEFAULT NULL COMMENT '审批时间（毫秒时间戳）',
    approval_comment TEXT DEFAULT NULL COMMENT '审批意见',

    -- 离职后信息
    rehire_eligible BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否可再录用',
    notes TEXT DEFAULT NULL COMMENT '备注',

    created_at BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    updated_at BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_emp_id (emp_id),
    INDEX idx_apply_date (apply_date),
    INDEX idx_last_work_date (last_work_date),
    INDEX idx_resignation_date (resignation_date),
    INDEX idx_approval (approval_status),

    -- 外键约束
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (emp_id) REFERENCES employees(emp_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工离职记录表';

-- ================================
-- 创建完成注释
-- ================================
-- 执行完成后，系统将拥有完整的组织中心数据表结构
-- 包括：组织、部门、岗位、员工、合同、调岗、离职等7个主要实体
-- 以及2个闭包表用于高效的树形结构查询
