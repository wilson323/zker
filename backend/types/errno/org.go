// backend/types/errno/org.go
package errno

import (
	"net/http"
)

// 组织中心相关错误码（ORG前缀）
// 注意：这些错误码暂不使用，预留未来功能扩展

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
)
