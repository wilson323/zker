// backend/types/errno/org.go
package errno

import (
	"net/http"
)

// 组织中心相关错误码（ORG前缀）

// 组织中心错误码常量
const (
	ErrOrgNotFoundCode             = 20001
	ErrOrgAlreadyExistsCode        = 20002
	ErrParentOrgNotFoundCode       = 20003
	ErrTenantMismatchCode          = 20004
	ErrInvalidParentOrgCode        = 20005
	ErrRootOrgMustBeCompanyCode    = 20006
	ErrTenantAlreadyHasOrgCode     = 20007
	ErrOrgCodeAlreadyExistsCode    = 20008
	ErrCannotFreezeRootOrgCode     = 20009
	ErrInvalidParamCode            = 20010
	ErrParentDeptNotFoundCode      = 20011
	ErrDeptOrgMismatchCode         = 20012
	ErrDeptCodeAlreadyExistsCode   = 20013
	ErrDeptNotFoundCode            = 20014
	ErrEmployeeNotFoundCode        = 20015
	ErrPositionNotFoundCode        = 20016
	ErrEmpCodeAlreadyExistsCode    = 20017
	ErrInvalidEmailCode            = 20018
	ErrEmailAlreadyExistsCode      = 20019
	ErrInvalidPhoneCode            = 20020
	ErrInvalidIDCardCode           = 20021
	ErrPositionCodeAlreadyExistsCode = 20022
	ErrPositionHasEmployeesCode    = 20023
	ErrDeptOrgMismatchCodeOld      = 20024
	ErrOrgLeaderNotFoundCode       = 20025
	ErrEmployeeNotActiveCode       = 20026
	ErrHandoverEmployeeNotFoundCode = 20027
	ErrCannotHandoverToSelfCode    = 20028
	ErrVirtualOrgNotFoundCode      = 20029
)

var (
	// ========== 组织管理错误（2xxx） ==========

	// 组织存在性错误 (20xxx)
	ErrOrgNotFound      = &BaseErrorCode{"ORG20001", "组织不存在", "Organization not found", "Organization not found", http.StatusNotFound}
	ErrOrgAlreadyExists = &BaseErrorCode{"ORG20002", "组织已存在", "Organization already exists", "Organization already exists", http.StatusConflict}
	ErrOrgCodeExists    = &BaseErrorCode{"ORG20003", "组织编码已存在", "Organization code already exists", "Organization code already exists", http.StatusConflict}
	ErrOrgNameExists    = &BaseErrorCode{"ORG20004", "组织名称已存在", "Organization name already exists", "Organization name already exists", http.StatusConflict}

	// 组织操作错误 (21xxx)
	ErrOrgCreateFailed  = &BaseErrorCode{"ORG21001", "创建组织失败", "Failed to create organization", "Failed to create organization", http.StatusInternalServerError}
	ErrOrgUpdateFailed  = &BaseErrorCode{"ORG21002", "更新组织失败", "Failed to update organization", "Failed to update organization", http.StatusInternalServerError}
	ErrOrgDeleteFailed  = &BaseErrorCode{"ORG21003", "删除组织失败", "Failed to delete organization", "Failed to delete organization", http.StatusInternalServerError}
	ErrOrgHasChildren   = &BaseErrorCode{"ORG21004", "组织下有子组织，无法删除", "Organization has children, cannot delete", "Organization has children, cannot delete", http.StatusForbidden}
	ErrOrgHasEmployees  = &BaseErrorCode{"ORG21005", "组织下有员工，无法删除", "Organization has employees, cannot delete", "Organization has employees, cannot delete", http.StatusForbidden}
	ErrOrgMoveFailed    = &BaseErrorCode{"ORG21006", "移动组织失败", "Failed to move organization", "Failed to move organization", http.StatusInternalServerError}
	ErrOrgInvalidParent = &BaseErrorCode{"ORG21007", "无效的父组织", "Invalid parent organization", "Invalid parent organization", http.StatusBadRequest}
	ErrOrgCycleDetected = &BaseErrorCode{"ORG21008", "检测到循环依赖", "Cycle detected in organization hierarchy", "Cycle detected in organization hierarchy", http.StatusBadRequest}
	ErrOrgLevelExceeded = &BaseErrorCode{"ORG21009", "组织层级超出限制", "Organization level exceeded limit", "Organization level exceeded limit", http.StatusBadRequest}

	// ========== 业务规则错误 ==========

	// 租户相关
	ErrTenantMismatch      = &BaseErrorCode{"ORG20004", "租户不匹配", "Tenant mismatch", "Tenant mismatch", http.StatusBadRequest}
	ErrTenantAlreadyHasOrg = &BaseErrorCode{"ORG20007", "租户已有组织", "Tenant already has organization", "Tenant already has organization", http.StatusConflict}

	// 父组织相关
	ErrParentOrgNotFound = &BaseErrorCode{"ORG20003", "父组织不存在", "Parent organization not found", "Parent organization not found", http.StatusNotFound}
	ErrInvalidParentOrg  = &BaseErrorCode{"ORG20005", "无效的父组织", "Invalid parent organization", "Invalid parent organization", http.StatusBadRequest}
	ErrRootOrgMustBeCompany = &BaseErrorCode{"ORG20006", "根组织必须是公司类型", "Root organization must be company type", "Root organization must be company type", http.StatusBadRequest}

	// 组织编码相关
	ErrOrgCodeAlreadyExists = &BaseErrorCode{"ORG20008", "组织编码已存在", "Organization code already exists", "Organization code already exists", http.StatusConflict}

	// 组织状态相关
	ErrCannotFreezeRootOrg = &BaseErrorCode{"ORG20009", "不能冻结根组织", "Cannot freeze root organization", "Cannot freeze root organization", http.StatusForbidden}

	// 通用参数错误
	ErrInvalidParam = &BaseErrorCode{"ORG20010", "无效参数", "Invalid parameter", "Invalid parameter", http.StatusBadRequest}

	// ========== 部门管理错误 ==========

	// 部门存在性错误
	ErrParentDeptNotFound = &BaseErrorCode{"ORG20011", "父部门不存在", "Parent department not found", "Parent department not found", http.StatusNotFound}
	ErrDeptNotFound       = &BaseErrorCode{"ORG20014", "部门不存在", "Department not found", "Department not found", http.StatusNotFound}

	// 部门业务规则错误
	ErrDeptOrgMismatch       = &BaseErrorCode{"ORG20012", "部门组织不匹配", "Department organization mismatch", "Department organization mismatch", http.StatusBadRequest}
	ErrDeptCodeAlreadyExists = &BaseErrorCode{"ORG20013", "部门编码已存在", "Department code already exists", "Department code already exists", http.StatusConflict}

	// ========== 员工管理错误 ==========

	// 员工存在性错误
	ErrEmployeeNotFound = &BaseErrorCode{"ORG20015", "员工不存在", "Employee not found", "Employee not found", http.StatusNotFound}

	// 员工业务规则错误
	ErrEmpCodeAlreadyExists = &BaseErrorCode{"ORG20017", "员工编码已存在", "Employee code already exists", "Employee code already exists", http.StatusConflict}
	ErrInvalidEmail         = &BaseErrorCode{"ORG20018", "无效邮箱", "Invalid email", "Invalid email", http.StatusBadRequest}
	ErrEmailAlreadyExists   = &BaseErrorCode{"ORG20019", "邮箱已存在", "Email already exists", "Email already exists", http.StatusConflict}
	ErrInvalidPhone         = &BaseErrorCode{"ORG20020", "无效手机号", "Invalid phone number", "Invalid phone number", http.StatusBadRequest}
	ErrInvalidIDCard        = &BaseErrorCode{"ORG20021", "无效身份证号", "Invalid ID card number", "Invalid ID card number", http.StatusBadRequest}
	ErrEmployeeNotActive    = &BaseErrorCode{"ORG20026", "员工未激活", "Employee not active", "Employee not active", http.StatusBadRequest}
	ErrHandoverEmployeeNotFound = &BaseErrorCode{"ORG20027", "交接员工不存在", "Handover employee not found", "Handover employee not found", http.StatusNotFound}
	ErrCannotHandoverToSelf = &BaseErrorCode{"ORG20028", "不能交接给自己", "Cannot handover to self", "Cannot handover to self", http.StatusBadRequest}

	// ========== 职位管理错误 ==========

	// 职位存在性错误
	ErrPositionNotFound = &BaseErrorCode{"ORG20016", "职位不存在", "Position not found", "Position not found", http.StatusNotFound}

	// 职位业务规则错误
	ErrPositionCodeAlreadyExists = &BaseErrorCode{"ORG20022", "职位编码已存在", "Position code already exists", "Position code already exists", http.StatusConflict}
	ErrPositionHasEmployees       = &BaseErrorCode{"ORG20023", "职位下有员工，无法删除", "Position has employees, cannot delete", "Position has employees, cannot delete", http.StatusForbidden}

	// ========== 矩阵组织错误 ==========

	ErrOrgLeaderNotFound = &BaseErrorCode{"ORG20025", "组织负责人不存在", "Organization leader not found", "Organization leader not found", http.StatusNotFound}

	// ========== 虚拟组织错误 ==========

	ErrVirtualOrgNotFound = &BaseErrorCode{"ORG20029", "虚拟组织不存在", "Virtual organization not found", "Virtual organization not found", http.StatusNotFound}
)
