package saga

import (
	"github.com/google/uuid"
)

// GenerateID 生成唯一ID
//
// 使用UUID v4算法生成全局唯一标识符，适用于Saga实例、执行记录等场景。
// 返回36字符的UUID字符串（格式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx）。
//
// 示例:
//
//	id := GenerateID()
//	fmt.Println(id) // "550e8400-e29b-41d4-a716-446655440000"
func GenerateID() string {
	return uuid.New().String()
}

// generateID 内部使用的ID生成函数（保持向后兼容）
// Deprecated: 使用 GenerateID 代替
func generateID() string {
	return GenerateID()
}
