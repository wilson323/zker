// backend/types/errno/auth.go
package errno

import (
	"net/http"
)

// 认证模块错误码（AUTH）
// 格式: AUTH{HTTP状态}{编号}
// 编号范围: 201001-201099 (客户端错误), 401001-401099 (认证错误), 403001-403099 (权限错误), 500001-500099 (服务器错误)

var (
	// ===== 客户端错误 (2xx) =====

	AUTH201001 = &BaseErrorCode{"AUTH201001", "用户名或密码错误", "用户名或密码错误", "Invalid username or password", http.StatusBadRequest}
	AUTH201002 = &BaseErrorCode{"AUTH201003", "用户名已存在", "用户名已存在", "Username already exists", http.StatusBadRequest}
	AUTH201003 = &BaseErrorCode{"AUTH201004", "邮箱格式错误", "邮箱格式错误", "Invalid email format", http.StatusBadRequest}
	AUTH201004 = &BaseErrorCode{"AUTH201005", "手机号格式错误", "手机号格式错误", "Invalid phone number format", http.StatusBadRequest}
	AUTH201005 = &BaseErrorCode{"AUTH201006", "密码格式错误", "密码格式错误", "Invalid password format", http.StatusBadRequest}
	AUTH201006 = &BaseErrorCode{"AUTH201007", "验证码错误", "验证码错误", "Invalid verification code", http.StatusBadRequest}
	AUTH201008 = &BaseErrorCode{"AUTH201008", "验证码已过期", "验证码已过期", "Verification code has expired", http.StatusBadRequest}
	AUTH201009 = &BaseErrorCode{"AUTH201009", "注册参数缺失", "注册参数缺失", "Missing registration parameters", http.StatusBadRequest}
	AUTH201010 = &BaseErrorCode{"AUTH201010", "用户名长度错误", "用户名长度必须在4-20个字符之间", "Username length must be between 4-20 characters", http.StatusBadRequest}

	// ===== 认证错误 (401) =====

	AUTH401001 = &BaseErrorCode{"AUTH401001", "未登录", "未登录", "Not logged in", http.StatusUnauthorized}
	AUTH401002 = &BaseErrorCode{"AUTH401002", "Token已过期", "Token已过期", "Token has expired", http.StatusUnauthorized}
	AUTH401003 = &BaseErrorCode{"AUTH401003", "无效的Token", "无效的Token", "Invalid token", http.StatusUnauthorized}
	AUTH401004 = &BaseErrorCode{"AUTH401004", "Token缺失", "Token缺失", "Missing token", http.StatusUnauthorized}
	AUTH401005 = &BaseErrorCode{"AUTH401005", "Refresh Token无效", "Refresh Token无效", "Invalid refresh token", http.StatusUnauthorized}
	AUTH401006 = &BaseErrorCode{"AUTH401006", "Refresh Token已过期", "Refresh Token已过期", "Refresh token has expired", http.StatusUnauthorized}
	AUTH401007 = &BaseErrorCode{"AUTH401007", "登录已过期", "登录已过期", "Login has expired", http.StatusUnauthorized}
	AUTH401008 = &BaseErrorCode{"AUTH401008", "会话已失效", "会话已失效", "Session has been invalidated", http.StatusUnauthorized}
	AUTH401009 = &BaseErrorCode{"AUTH401009", "重复登录", "检测到重复登录，请重新登录", "Duplicate login detected", http.StatusUnauthorized}
	AUTH401010 = &BaseErrorCode{"AUTH401010", "登录失败次数过多", "登录失败次数过多，请稍后重试", "Too many failed login attempts", http.StatusUnauthorized}

	// ===== 权限错误 (403) =====

	AUTH403001 = &BaseErrorCode{"AUTH403001", "账号已禁用", "账号已禁用", "Account has been disabled", http.StatusForbidden}
	AUTH403002 = &BaseErrorCode{"AUTH403002", "账号已锁定", "账号已锁定", "Account has been locked", http.StatusForbidden}
	AUTH403003 = &BaseErrorCode{"AUTH403003", "账号已封禁", "账号已封禁", "Account has been banned", http.StatusForbidden}
	AUTH403004 = &BaseErrorCode{"AUTH403004", "权限不足", "权限不足", "Permission denied", http.StatusForbidden}
	AUTH403005 = &BaseErrorCode{"AUTH403005", "需要管理员权限", "需要管理员权限", "Admin permission required", http.StatusForbidden}

	// ===== 资源错误 (404) =====

	AUTH404001 = &BaseErrorCode{"AUTH404001", "用户不存在", "用户不存在", "User not found", http.StatusNotFound}
	AUTH404002 = &BaseErrorCode{"AUTH404002", "角色不存在", "角色不存在", "Role not found", http.StatusNotFound}

	// ===== 服务器错误 (500) =====

	AUTH500001 = &BaseErrorCode{"AUTH500001", "认证服务异常", "认证服务异常", "Authentication service error", http.StatusInternalServerError}
	AUTH500002 = &BaseErrorCode{"AUTH500002", "Token生成失败", "Token生成失败", "Failed to generate token", http.StatusInternalServerError}
	AUTH500003 = &BaseErrorCode{"AUTH500003", "密码加密失败", "密码加密失败", "Failed to encrypt password", http.StatusInternalServerError}
	AUTH500004 = &BaseErrorCode{"AUTH500004", "会话创建失败", "会话创建失败", "Failed to create session", http.StatusInternalServerError}
	AUTH500005 = &BaseErrorCode{"AUTH500005", "验证码发送失败", "验证码发送失败", "Failed to send verification code", http.StatusInternalServerError}
	AUTH500006 = &BaseErrorCode{"AUTH500006", "用户创建失败", "用户创建失败", "Failed to create user", http.StatusInternalServerError}
	AUTH500007 = &BaseErrorCode{"AUTH500007", "用户更新失败", "用户更新失败", "Failed to update user", http.StatusInternalServerError}
	AUTH500008 = &BaseErrorCode{"AUTH500008", "用户删除失败", "用户删除失败", "Failed to delete user", http.StatusInternalServerError}

	// ===== 第三方登录错误 (42x) =====

	AUTH422001 = &BaseErrorCode{"AUTH422001", "第三方登录失败", "第三方登录失败", "Third-party login failed", http.StatusUnprocessableEntity}
	AUTH422002 = &BaseErrorCode{"AUTH422002", "第三方授权失败", "第三方授权失败", "Third-party authorization failed", http.StatusUnprocessableEntity}
	AUTH422003 = &BaseErrorCode{"AUTH422003", "未绑定账号", "未绑定账号", "Account not linked", http.StatusUnprocessableEntity}
	AUTH422004 = &BaseErrorCode{"AUTH422004", "已绑定其他账号", "已绑定其他账号", "Already linked to another account", http.StatusUnprocessableEntity}

	// ===== OAuth相关错误 =====

	AUTH426001 = &BaseErrorCode{"AUTH426001", "无效的OAuth客户端", "无效的OAuth客户端", "Invalid OAuth client", http.StatusUpgradeRequired}
	AUTH426002 = &BaseErrorCode{"AUTH426002", "OAuth授权码无效", "OAuth授权码无效", "Invalid OAuth authorization code", http.StatusUpgradeRequired}
	AUTH426003 = &BaseErrorCode{"AUTH426003", "OAuth scope无效", "OAuth scope无效", "Invalid OAuth scope", http.StatusUpgradeRequired}

	// ===== 账号安全相关错误 =====

	AUTH423001 = &BaseErrorCode{"AUTH423001", "密码过于简单", "密码过于简单", "Password is too simple", http.StatusUnprocessableEntity}
	AUTH423002 = &BaseErrorCode{"AUTH423002", "密码与旧密码相同", "新密码不能与旧密码相同", "New password must be different from old password", http.StatusUnprocessableEntity}
	AUTH423003 = &BaseErrorCode{"AUTH423003", "旧密码错误", "旧密码错误", "Old password is incorrect", http.StatusUnprocessableEntity}
	AUTH423004 = &BaseErrorCode{"AUTH423004", "两次密码不一致", "两次输入的密码不一致", "Passwords do not match", http.StatusUnprocessableEntity}
	AUTH423005 = &BaseErrorCode{"AUTH423005", "需要验证旧密码", "需要验证旧密码", "Old password verification required", http.StatusUnprocessableEntity}
	AUTH423006 = &BaseErrorCode{"AUTH423006", "密码重置次数过多", "密码重置次数过多，请稍后重试", "Too many password reset attempts", http.StatusUnprocessableEntity}
	AUTH423007 = &BaseErrorCode{"AUTH423007", "验证码错误次数过多", "验证码错误次数过多", "Too many incorrect verification code attempts", http.StatusUnprocessableEntity}

	// ===== 邮箱验证相关错误 =====

	AUTH424001 = &BaseErrorCode{"AUTH424001", "邮箱未验证", "邮箱未验证", "Email has not been verified", http.StatusFailedDependency}
	AUTH424002 = &BaseErrorCode{"AUTH424002", "邮箱已验证", "邮箱已验证", "Email has already been verified", http.StatusFailedDependency}
	AUTH424003 = &BaseErrorCode{"AUTH424003", "验证链接已失效", "验证链接已失效", "Verification link has expired", http.StatusFailedDependency}
	AUTH424004 = &BaseErrorCode{"AUTH424004", "验证链接无效", "验证链接无效", "Invalid verification link", http.StatusFailedDependency}
)
