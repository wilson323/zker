// backend/types/errno/chat.go
package errno

var (
	CHAT201001 = &BaseErrorCode{"CHAT201001", "消息内容为空", "消息内容为空", "Message content is empty", 400}
	CHAT201002 = &BaseErrorCode{"CHAT201002", "消息过长", "消息过长", "Message too long", 400}
	CHAT404001 = &BaseErrorCode{"CHAT404001", "对话不存在", "对话不存在", "Conversation not found", 404}
	CHAT403001 = &BaseErrorCode{"CHAT403001", "无权访问对话", "无权访问对话", "No permission to access conversation", 403}
	CHAT500001 = &BaseErrorCode{"CHAT500001", "消息发送失败", "消息发送失败", "Failed to send message", 500}
)
