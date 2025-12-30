// backend/types/errno/bot.go
package errno

var (
	BOT201001 = &BaseErrorCode{"BOT201001", "Bot名称格式错误", "Bot名称格式错误", "Invalid bot name format", 400}
	BOT201002 = &BaseErrorCode{"BOT201002", "Bot描述过长", "Bot描述过长", "Bot description too long", 400}
	BOT401001 = &BaseErrorCode{"BOT401001", "Bot不存在", "Bot不存在", "Bot not found", 404}
	BOT403001 = &BaseErrorCode{"BOT403001", "无权访问Bot", "无权访问Bot", "No permission to access bot", 403}
	BOT403002 = &BaseErrorCode{"BOT403002", "Bot已发布", "Bot已发布", "Bot has been published", 403}
	BOT403003 = &BaseErrorCode{"BOT403003", "Bot已下线", "Bot已下线", "Bot has been offline", 403}
	BOT409001 = &BaseErrorCode{"BOT409001", "Bot名称已存在", "Bot名称已存在", "Bot name already exists", 409}
	BOT500001 = &BaseErrorCode{"BOT500001", "Bot创建失败", "Bot创建失败", "Failed to create bot", 500}
	BOT500002 = &BaseErrorCode{"BOT500002", "Bot更新失败", "Bot更新失败", "Failed to update bot", 500}
	BOT500003 = &BaseErrorCode{"BOT500003", "Bot删除失败", "Bot删除失败", "Failed to delete bot", 500}
)
