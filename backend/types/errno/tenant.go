// backend/types/errno/tenant.go
package errno

import (
	"net/http"
)

// 租户相关错误码常量
const (
	ErrTenantNotFoundCode        = 2001001
	ErrTenantAlreadyExistsCode   = 2001002
	ErrTenantSuspendedCode       = 2003001
	ErrTenantDeletedCode         = 2001003
	ErrTenantInvalidParamCode    = 2002001
	ErrInvalidTenantIDCode       = 2002002
	// 新增错误码常量
	ErrTenantInactiveCode                = 2003002
	ErrTenantQuotaExceededCode           = 2004001
	ErrTenantInvalidNameCode             = 2002003
	ErrTenantInvalidConfigCode           = 2002004
	ErrTenantMigrationFailedCode         = 2005001
	ErrTenantIsolationNotSupportedCode   = 2002005
	ErrTenantSubscriptionExpiredCode     = 2003003
	ErrMissingTenantIDCode               = 2002999 // 租户ID缺失（严格租户隔离模式）
	ErrTenantMigrationRecordNotFoundCode = 2005002
	ErrTenantMigrationAlreadyRolledBackCode = 2005003
	ErrTenantMigrationNotCompleteCode    = 2005004
	ErrCrossTenantAccessCode             = 2002041 // 跨租户访问
	ErrCrossTenantUpdateCode             = 2002042 // 跨租户更新
)

var (
	// 租户存在性错误
	ErrTenantNotFound = &BaseErrorCode{
		code:       "TENANT_NOT_FOUND",
		message:    "Tenant not found",
		messageZH:  "租户不存在",
		messageEN:  "Tenant not found",
		httpStatus: http.StatusNotFound,
	}
	ErrTenantAlreadyExists = &BaseErrorCode{
		code:       "TENANT_ALREADY_EXISTS",
		message:    "Tenant already exists",
		messageZH:  "租户已存在",
		messageEN:  "Tenant already exists",
		httpStatus: http.StatusConflict,
	}
	ErrTenantSuspended = &BaseErrorCode{
		code:       "TENANT_SUSPENDED",
		message:    "Tenant is suspended",
		messageZH:  "租户已暂停",
		messageEN:  "Tenant is suspended",
		httpStatus: http.StatusForbidden,
	}
	ErrTenantDeleted = &BaseErrorCode{
		code:       "TENANT_DELETED",
		message:    "Tenant is deleted",
		messageZH:  "租户已删除",
		messageEN:  "Tenant is deleted",
		httpStatus: http.StatusGone,
	}
	ErrTenantExpired = &BaseErrorCode{
		code:       "TENANT_EXPIRED",
		message:    "Tenant is expired",
		messageZH:  "租户已过期",
		messageEN:  "Tenant is expired",
		httpStatus: http.StatusForbidden,
	}
	ErrTenantInactive = &BaseErrorCode{
		code:       "TENANT_INACTIVE",
		message:    "Tenant is inactive",
		messageZH:  "租户未激活",
		messageEN:  "Tenant is inactive",
		httpStatus: http.StatusForbidden,
	}

	// 租户注册错误
	ErrCompanyNameExists = &BaseErrorCode{
		code:       "COMPANY_NAME_EXISTS",
		message:    "Company name already exists",
		messageZH:  "企业名称已存在",
		messageEN:  "Company name already exists",
		httpStatus: http.StatusConflict,
	}
	ErrSubdomainExists = &BaseErrorCode{
		code:       "SUBDOMAIN_EXISTS",
		message:    "Subdomain already exists",
		messageZH:  "子域名已存在",
		messageEN:  "Subdomain already exists",
		httpStatus: http.StatusConflict,
	}
	ErrSubdomainInvalid = &BaseErrorCode{
		code:       "SUBDOMAIN_INVALID",
		message:    "Invalid subdomain format",
		messageZH:  "子域名格式错误",
		messageEN:  "Invalid subdomain format",
		httpStatus: http.StatusBadRequest,
	}
	ErrSubdomainReserved = &BaseErrorCode{
		code:       "SUBDOMAIN_RESERVED",
		message:    "Subdomain is reserved",
		messageZH:  "子域名被保留",
		messageEN:  "Subdomain is reserved",
		httpStatus: http.StatusForbidden,
	}
	ErrVerificationCodeInvalid = &BaseErrorCode{
		code:       "VERIFICATION_CODE_INVALID",
		message:    "Invalid verification code",
		messageZH:  "验证码错误",
		messageEN:  "Invalid verification code",
		httpStatus: http.StatusBadRequest,
	}
	ErrVerificationCodeExpired = &BaseErrorCode{
		code:       "VERIFICATION_CODE_EXPIRED",
		message:    "Verification code expired",
		messageZH:  "验证码已过期",
		messageEN:  "Verification code expired",
		httpStatus: http.StatusBadRequest,
	}
	ErrVerificationCodeRateLimit = &BaseErrorCode{
		code:       "VERIFICATION_CODE_RATE_LIMIT",
		message:    "Verification code send too frequently",
		messageZH:  "验证码发送过于频繁，请1分钟后重试",
		messageEN:  "Verification code send too frequently, please try again later",
		httpStatus: http.StatusTooManyRequests,
	}
	ErrVerificationCodeSendFailed = &BaseErrorCode{
		code:       "VERIFICATION_CODE_SEND_FAILED",
		message:    "Failed to send verification code",
		messageZH:  "验证码发送失败",
		messageEN:  "Failed to send verification code",
		httpStatus: http.StatusInternalServerError,
	}
	ErrVerificationCodeStoreFailed = &BaseErrorCode{
		code:       "VERIFICATION_CODE_STORE_FAILED",
		message:    "Failed to store verification code",
		messageZH:  "验证码存储失败",
		messageEN:  "Failed to store verification code",
		httpStatus: http.StatusInternalServerError,
	}
	ErrVerificationCodeCheckFailed = &BaseErrorCode{
		code:       "VERIFICATION_CODE_CHECK_FAILED",
		message:    "Failed to check rate limit",
		messageZH:  "验证码频率检查失败",
		messageEN:  "Failed to check rate limit",
		httpStatus: http.StatusInternalServerError,
	}
	ErrVerificationCodeGetFailed = &BaseErrorCode{
		code:       "VERIFICATION_CODE_GET_FAILED",
		message:    "Failed to get verification code",
		messageZH:  "获取验证码失败",
		messageEN:  "Failed to get verification code",
		httpStatus: http.StatusInternalServerError,
	}
	ErrEmailAlreadyRegistered = &BaseErrorCode{
		code:       "EMAIL_ALREADY_REGISTERED",
		message:    "Email already registered",
		messageZH:  "邮箱已注册",
		messageEN:  "Email already registered",
		httpStatus: http.StatusConflict,
	}

	// 租户操作错误
	ErrTenantUpdateFailed = &BaseErrorCode{
		code:       "TENANT_UPDATE_FAILED",
		message:    "Failed to update tenant",
		messageZH:  "租户更新失败",
		messageEN:  "Failed to update tenant",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantDeleteFailed = &BaseErrorCode{
		code:       "TENANT_DELETE_FAILED",
		message:    "Failed to delete tenant",
		messageZH:  "租户删除失败",
		messageEN:  "Failed to delete tenant",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantCheckFailed = &BaseErrorCode{
		code:       "TENANT_CHECK_FAILED",
		message:    "Failed to check tenant status",
		messageZH:  "租户检查失败",
		messageEN:  "Failed to check tenant status",
		httpStatus: http.StatusInternalServerError,
	}

	// 租户参数错误
	ErrTenantInvalidName = &BaseErrorCode{
		code:       "TENANT_INVALID_NAME",
		message:    "Invalid tenant name",
		messageZH:  "租户名称无效",
		messageEN:  "Invalid tenant name",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantInvalidConfig = &BaseErrorCode{
		code:       "TENANT_INVALID_CONFIG",
		message:    "Invalid tenant config",
		messageZH:  "租户配置无效",
		messageEN:  "Invalid tenant config",
		httpStatus: http.StatusBadRequest,
	}

	// 租户配额错误
	ErrTenantQuotaExceeded = &BaseErrorCode{
		code:       "TENANT_QUOTA_EXCEEDED",
		message:    "Tenant quota exceeded",
		messageZH:  "租户配额已用尽",
		messageEN:  "Tenant quota exceeded",
		httpStatus: http.StatusForbidden,
	}

	// 租户迁移错误
	ErrTenantMigrationFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_FAILED",
		message:    "Tenant migration failed",
		messageZH:  "租户迁移失败",
		messageEN:  "Tenant migration failed",
		httpStatus: http.StatusInternalServerError,
	}

	// 租户隔离错误
	ErrTenantIsolationNotSupported = &BaseErrorCode{
		code:       "TENANT_ISOLATION_NOT_SUPPORTED",
		message:    "Tenant isolation not supported",
		messageZH:  "不支持租户隔离",
		messageEN:  "Tenant isolation not supported",
		httpStatus: http.StatusBadRequest,
	}

	// 租户订阅错误
	ErrTenantSubscriptionExpired = &BaseErrorCode{
		code:       "TENANT_SUBSCRIPTION_EXPIRED",
		message:    "Tenant subscription expired",
		messageZH:  "租户订阅已过期",
		messageEN:  "Tenant subscription expired",
		httpStatus: http.StatusForbidden,
	}

	// 租户验证错误
	ErrTenantNameTooShort = &BaseErrorCode{
		code:       "TENANT_NAME_TOO_SHORT",
		message:    "Company name too short (min 2 characters)",
		messageZH:  "企业名称长度必须在2-200个字符之间",
		messageEN:  "Company name too short (min 2 characters)",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantNameTooLong = &BaseErrorCode{
		code:       "TENANT_NAME_TOO_LONG",
		message:    "Company name too long (max 200 characters)",
		messageZH:  "企业名称长度必须在2-200个字符之间",
		messageEN:  "Company name too long (max 200 characters)",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantSubdomainTooShort = &BaseErrorCode{
		code:       "TENANT_SUBDOMAIN_TOO_SHORT",
		message:    "Subdomain too short (min 3 characters)",
		messageZH:  "子域名长度必须在3-63个字符之间",
		messageEN:  "Subdomain too short (min 3 characters)",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantSubdomainTooLong = &BaseErrorCode{
		code:       "TENANT_SUBDOMAIN_TOO_LONG",
		message:    "Subdomain too long (max 63 characters)",
		messageZH:  "子域名长度必须在3-63个字符之间",
		messageEN:  "Subdomain too long (max 63 characters)",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantSubdomainInvalidFormat = &BaseErrorCode{
		code:       "TENANT_SUBDOMAIN_INVALID_FORMAT",
		message:    "Subdomain cannot start or end with hyphen",
		messageZH:  "子域名不能以连字符开头或结尾",
		messageEN:  "Subdomain cannot start or end with hyphen",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantRequestEmpty = &BaseErrorCode{
		code:       "TENANT_REQUEST_EMPTY",
		message:    "Request cannot be empty",
		messageZH:  "请求不能为空",
		messageEN:  "Request cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantNameEmpty = &BaseErrorCode{
		code:       "TENANT_NAME_EMPTY",
		message:    "Company name cannot be empty",
		messageZH:  "企业名称不能为空",
		messageEN:  "Company name cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantSubdomainEmpty = &BaseErrorCode{
		code:       "TENANT_SUBDOMAIN_EMPTY",
		message:    "Subdomain cannot be empty",
		messageZH:  "子域名不能为空",
		messageEN:  "Subdomain cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantEmailEmpty = &BaseErrorCode{
		code:       "TENANT_EMAIL_EMPTY",
		message:    "Email cannot be empty",
		messageZH:  "邮箱不能为空",
		messageEN:  "Email cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantCodeEmpty = &BaseErrorCode{
		code:       "TENANT_CODE_EMPTY",
		message:    "Verification code cannot be empty",
		messageZH:  "验证码不能为空",
		messageEN:  "Verification code cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantCodeInvalidFormat = &BaseErrorCode{
		code:       "TENANT_CODE_INVALID_FORMAT",
		message:    "Verification code format invalid",
		messageZH:  "验证码格式不正确",
		messageEN:  "Verification code format invalid",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantAdminNameEmpty = &BaseErrorCode{
		code:       "TENANT_ADMIN_NAME_EMPTY",
		message:    "Admin name cannot be empty",
		messageZH:  "管理员姓名不能为空",
		messageEN:  "Admin name cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantAdminPasswordEmpty = &BaseErrorCode{
		code:       "TENANT_ADMIN_PASSWORD_EMPTY",
		message:    "Admin password cannot be empty",
		messageZH:  "管理员密码不能为空",
		messageEN:  "Admin password cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantAdminPasswordTooShort = &BaseErrorCode{
		code:       "TENANT_ADMIN_PASSWORD_TOO_SHORT",
		message:    "Admin password too short (min 8 characters)",
		messageZH:  "管理员密码长度不能少于8位",
		messageEN:  "Admin password too short (min 8 characters)",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantTypeEmpty = &BaseErrorCode{
		code:       "TENANT_TYPE_EMPTY",
		message:    "Tenant type cannot be empty",
		messageZH:  "租户类型不能为空",
		messageEN:  "Tenant type cannot be empty",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantTypeInvalid = &BaseErrorCode{
		code:       "TENANT_TYPE_INVALID",
		message:    "Invalid tenant type",
		messageZH:  "无效的租户类型",
		messageEN:  "Invalid tenant type",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantVerifyFailed = &BaseErrorCode{
		code:       "TENANT_VERIFY_FAILED",
		message:    "Verification code validation failed",
		messageZH:  "验证码验证失败",
		messageEN:  "Verification code validation failed",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantCodeInvalidOrExpired = &BaseErrorCode{
		code:       "TENANT_CODE_INVALID_OR_EXPIRED",
		message:    "Verification code invalid or expired",
		messageZH:  "验证码错误或已过期",
		messageEN:  "Verification code invalid or expired",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantCompanyNameOccupied = &BaseErrorCode{
		code:       "TENANT_COMPANY_NAME_OCCUPIED",
		message:    "Company name already occupied",
		messageZH:  "企业名称已被占用",
		messageEN:  "Company name already occupied",
		httpStatus: http.StatusConflict,
	}
	ErrTenantSubdomainOccupied = &BaseErrorCode{
		code:       "TENANT_SUBDOMAIN_OCCUPIED",
		message:    "Subdomain already occupied",
		messageZH:  "子域名已被占用",
		messageEN:  "Subdomain already occupied",
		httpStatus: http.StatusConflict,
	}
	ErrTenantCheckNameFailed = &BaseErrorCode{
		code:       "TENANT_CHECK_NAME_FAILED",
		message:    "Failed to check company name",
		messageZH:  "企业名称检查失败",
		messageEN:  "Failed to check company name",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantCheckSubdomainFailed = &BaseErrorCode{
		code:       "TENANT_CHECK_SUBDOMAIN_FAILED",
		message:    "Failed to check subdomain",
		messageZH:  "子域名检查失败",
		messageEN:  "Failed to check subdomain",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantCreateFailed = &BaseErrorCode{
		code:       "TENANT_CREATE_FAILED",
		message:    "Failed to create tenant",
		messageZH:  "租户创建失败",
		messageEN:  "Failed to create tenant",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantSubscriptionCreateFailed = &BaseErrorCode{
		code:       "TENANT_SUBSCRIPTION_CREATE_FAILED",
		message:    "Failed to create subscription",
		messageZH:  "订阅创建失败",
		messageEN:  "Failed to create subscription",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantQuotaCreateFailed = &BaseErrorCode{
		code:       "TENANT_QUOTA_CREATE_FAILED",
		message:    "Failed to create quota",
		messageZH:  "配额创建失败",
		messageEN:  "Failed to create quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantQuotaNotFound = &BaseErrorCode{
		code:       "TENANT_QUOTA_NOT_FOUND",
		message:    "Quota not found for resource type",
		messageZH:  "配额不存在",
		messageEN:  "Quota not found for resource type",
		httpStatus: http.StatusNotFound,
	}
	ErrTenantSubscriptionNotFound = &BaseErrorCode{
		code:       "TENANT_SUBSCRIPTION_NOT_FOUND",
		message:    "Subscription not found for tenant",
		messageZH:  "订阅不存在",
		messageEN:  "Subscription not found for tenant",
		httpStatus: http.StatusNotFound,
	}
	ErrTenantInvalidUpgrade = &BaseErrorCode{
		code:       "TENANT_INVALID_UPGRADE",
		message:    "Invalid upgrade path",
		messageZH:  "无效的升级路径",
		messageEN:  "Invalid upgrade path",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantGetQuotaFailed = &BaseErrorCode{
		code:       "TENANT_GET_QUOTA_FAILED",
		message:    "Failed to get quota",
		messageZH:  "获取配额失败",
		messageEN:  "Failed to get quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantUpdateQuotaFailed = &BaseErrorCode{
		code:       "TENANT_UPDATE_QUOTA_FAILED",
		message:    "Failed to update quota",
		messageZH:  "更新配额失败",
		messageEN:  "Failed to update quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantUpdateSubscriptionTierFailed = &BaseErrorCode{
		code:       "TENANT_UPDATE_SUBSCRIPTION_TIER_FAILED",
		message:    "Failed to update subscription tier",
		messageZH:  "更新订阅层级失败",
		messageEN:  "Failed to update subscription tier",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantUpdateQuotaLimitsFailed = &BaseErrorCode{
		code:       "TENANT_UPDATE_QUOTA_LIMITS_FAILED",
		message:    "Failed to update quota limits",
		messageZH:  "更新配额限制失败",
		messageEN:  "Failed to update quota limits",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantGetSubscriptionFailed = &BaseErrorCode{
		code:       "TENANT_GET_SUBSCRIPTION_FAILED",
		message:    "Failed to get subscription",
		messageZH:  "获取订阅失败",
		messageEN:  "Failed to get subscription",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantCannotDowngrade = &BaseErrorCode{
		code:       "TENANT_CANNOT_DOWNGRADE",
		message:    "Cannot downgrade subscription",
		messageZH:  "不能降级订阅",
		messageEN:  "Cannot downgrade subscription",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantCancelSubscriptionFailed = &BaseErrorCode{
		code:       "TENANT_CANCEL_SUBSCRIPTION_FAILED",
		message:    "Failed to cancel subscription",
		messageZH:  "取消订阅失败",
		messageEN:  "Failed to cancel subscription",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantGetExistingSubscriptionFailed = &BaseErrorCode{
		code:       "TENANT_GET_EXISTING_SUBSCRIPTION_FAILED",
		message:    "Failed to get existing subscription",
		messageZH:  "获取现有订阅失败",
		messageEN:  "Failed to get existing subscription",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantUpdateSubscriptionFailed = &BaseErrorCode{
		code:       "TENANT_UPDATE_SUBSCRIPTION_FAILED",
		message:    "Failed to update subscription",
		messageZH:  "更新订阅失败",
		messageEN:  "Failed to update subscription",
		httpStatus: http.StatusInternalServerError,
	}

	// 租户迁移错误
	ErrTenantMigrationValidateFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_VALIDATE_FAILED",
		message:    "Migration validation failed",
		messageZH:  "验证迁移失败",
		messageEN:  "Migration validation failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationNullTenantID = &BaseErrorCode{
		code:       "TENANT_MIGRATION_NULL_TENANT_ID",
		message:    "Records with NULL tenant_id found",
		messageZH:  "存在tenant_id为NULL的记录",
		messageEN:  "Records with NULL tenant_id found",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationNotNullFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_NOT_NULL_FAILED",
		message:    "Failed to modify tenant_id to NOT NULL",
		messageZH:  "修改tenant_id为NOT NULL失败",
		messageEN:  "Failed to modify tenant_id to NOT NULL",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationValidateFieldTypeFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_VALIDATE_FIELD_TYPE_FAILED",
		message:    "Failed to validate field type",
		messageZH:  "验证字段类型失败",
		messageEN:  "Failed to validate field type",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationRollbackFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_ROLLBACK_FAILED",
		message:    "Migration rollback failed",
		messageZH:  "回滚切换失败",
		messageEN:  "Migration rollback failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationCutoverNotComplete = &BaseErrorCode{
		code:       "TENANT_MIGRATION_CUTOVER_NOT_COMPLETE",
		message:    "Cutover not complete",
		messageZH:  "切换未完成",
		messageEN:  "Cutover not complete",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationValidateNullFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_VALIDATE_NULL_FAILED",
		message:    "Failed to validate NULL",
		messageZH:  "验证NULL失败",
		messageEN:  "Failed to validate NULL",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationValidateForeignKeyFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_VALIDATE_FOREIGN_KEY_FAILED",
		message:    "Failed to validate foreign key",
		messageZH:  "验证外键失败",
		messageEN:  "Failed to validate foreign key",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationCutoverTablesFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_CUTOVER_TABLES_FAILED",
		message:    "Failed to cutover tables",
		messageZH:  "表切换失败",
		messageEN:  "Failed to cutover tables",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationDualWriteFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_DUAL_WRITE_FAILED",
		message:    "Failed to query data",
		messageZH:  "查询数据失败",
		messageEN:  "Failed to query data",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationConsistencyTooLow = &BaseErrorCode{
		code:       "TENANT_MIGRATION_CONSISTENCY_TOO_LOW",
		message:    "Data consistency below 99.9%",
		messageZH:  "数据一致性低于99.9%",
		messageEN:  "Data consistency below 99.9%",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationDDLNonRetryable = &BaseErrorCode{
		code:       "TENANT_MIGRATION_DDL_NON_RETRYABLE",
		message:    "Non-retryable DDL error",
		messageZH:  "不可重试错误",
		messageEN:  "Non-retryable DDL error",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationGetTotalFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_GET_TOTAL_FAILED",
		message:    "Failed to get total records",
		messageZH:  "获取总记录数失败",
		messageEN:  "Failed to get total records",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationSaveRecordFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_SAVE_RECORD_FAILED",
		message:    "Failed to save migration record",
		messageZH:  "保存迁移记录失败",
		messageEN:  "Failed to save migration record",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationAddFieldFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_ADD_FIELD_FAILED",
		message:    "Failed to add field",
		messageZH:  "添加字段失败",
		messageEN:  "Failed to add field",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationEnableDualWriteFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_ENABLE_DUAL_WRITE_FAILED",
		message:    "Failed to enable dual write",
		messageZH:  "启用双写失败",
		messageEN:  "Failed to enable dual write",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationBackfillFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_BACKFILL_FAILED",
		message:    "Failed to backfill data",
		messageZH:  "数据回填失败",
		messageEN:  "Failed to backfill data",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationModeNotSupported = &BaseErrorCode{
		code:       "TENANT_MIGRATION_MODE_NOT_SUPPORTED",
		message:    "Unsupported migration mode",
		messageZH:  "不支持的迁移模式",
		messageEN:  "Unsupported migration mode",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantMigrationExecuteSQLFailed = &BaseErrorCode{
		code:       "TENANT_MIGRATION_EXECUTE_SQL_FAILED",
		message:    "Failed to execute SQL",
		messageZH:  "执行SQL失败",
		messageEN:  "Failed to execute SQL",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantMigrationRecordNotFound = &BaseErrorCode{
		code:       "TENANT_MIGRATION_RECORD_NOT_FOUND",
		message:    "Migration record not found",
		messageZH:  "迁移记录不存在",
		messageEN:  "Migration record not found",
		httpStatus: http.StatusNotFound,
	}
	ErrTenantMigrationAlreadyRolledBack = &BaseErrorCode{
		code:       "TENANT_MIGRATION_ALREADY_ROLLED_BACK",
		message:    "Migration already rolled back",
		messageZH:  "迁移已经回滚",
		messageEN:  "Migration already rolled back",
		httpStatus: http.StatusConflict,
	}
	ErrTenantMigrationNotComplete = &BaseErrorCode{
		code:       "TENANT_MIGRATION_NOT_COMPLETE",
		message:    "Migration not complete",
		messageZH:  "迁移尚未完成",
		messageEN:  "Migration not complete",
		httpStatus: http.StatusConflict,
	}
	ErrTenantIsolationUpgradeSchemaFailed = &BaseErrorCode{
		code:       "TENANT_ISOLATION_UPGRADE_SCHEMA_FAILED",
		message:    "Failed to upgrade to schema level",
		messageZH:  "升级到schema级别失败",
		messageEN:  "Failed to upgrade to schema level",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantIsolationUpgradeDatabaseFailed = &BaseErrorCode{
		code:       "TENANT_ISOLATION_UPGRADE_DATABASE_FAILED",
		message:    "Failed to upgrade to database level",
		messageZH:  "升级到database级别失败",
		messageEN:  "Failed to upgrade to database level",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantIsolationUnknownStrategy = &BaseErrorCode{
		code:       "TENANT_ISOLATION_UNKNOWN_STRATEGY",
		message:    "Unknown target strategy",
		messageZH:  "未知的目标策略",
		messageEN:  "Unknown target strategy",
		httpStatus: http.StatusBadRequest,
	}
	ErrTenantIsolationUpdateStrategyFailed = &BaseErrorCode{
		code:       "TENANT_ISOLATION_UPDATE_STRATEGY_FAILED",
		message:    "Failed to update tenant isolation strategy",
		messageZH:  "更新租户隔离策略失败",
		messageEN:  "Failed to update tenant isolation strategy",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantIsolationCreateSchemaFailed = &BaseErrorCode{
		code:       "TENANT_ISOLATION_CREATE_SCHEMA_FAILED",
		message:    "Failed to create schema",
		messageZH:  "创建schema失败",
		messageEN:  "Failed to create schema",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantIsolationCreateTableFailed = &BaseErrorCode{
		code:       "TENANT_ISOLATION_CREATE_TABLE_FAILED",
		message:    "Failed to create table",
		messageZH:  "创建表失败",
		messageEN:  "Failed to create table",
		httpStatus: http.StatusInternalServerError,
	}
	ErrTenantIsolationMigrateDataFailed = &BaseErrorCode{
		code:       "TENANT_ISOLATION_MIGRATE_DATA_FAILED",
		message:    "Failed to migrate data",
		messageZH:  "迁移数据失败",
		messageEN:  "Failed to migrate data",
		httpStatus: http.StatusInternalServerError,
	}

	// 租户ID缺失错误（严格租户隔离模式）
	ErrMissingTenantID = &BaseErrorCode{
		code:       "TENANT_ID_MISSING",
		message:    "Tenant ID is required",
		messageZH:  "租户ID缺失",
		messageEN:  "Tenant ID is required",
		httpStatus: http.StatusUnauthorized,
	}

	// 跨租户访问错误
	ErrCrossTenantAccess = &BaseErrorCode{
		code:       "CROSS_TENANT_ACCESS",
		message:    "Cross-tenant access denied",
		messageZH:  "跨租户访问拒绝",
		messageEN:  "Cross-tenant access denied",
		httpStatus: http.StatusForbidden,
	}
	ErrCrossTenantUpdate = &BaseErrorCode{
		code:       "CROSS_TENANT_UPDATE",
		message:    "Cross-tenant update denied",
		messageZH:  "跨租户更新拒绝",
		messageEN:  "Cross-tenant update denied",
		httpStatus: http.StatusForbidden,
	}
)
