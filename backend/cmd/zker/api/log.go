package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry 日志条目
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	BotID     string `json:"bot_id"`
	RunID     string `json:"run_id"`
}

// LogQuery 日志查询参数
type LogQuery struct {
	Tail   int       `json:"tail"`
	Level  string    `json:"level"`
	Since  time.Time `json:"since"`
	Follow bool      `json:"follow"`
}

// GetLogs 获取日志
func (c *Client) GetLogs(botID string, query *LogQuery) ([]LogEntry, error) {
	resp, err := c.get(fmt.Sprintf("/api/bot/%s/logs?tail=%d&level=%s", botID, query.Tail, query.Level))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get logs (status: %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    []LogEntry `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return result.Data, nil
}

// StreamLogs 实时流式日志
func (c *Client) StreamLogs(botID string, level string, callback func(LogEntry)) error {
	// TODO: 实现WebSocket或SSE流式日志
	return fmt.Errorf("streaming logs not implemented yet")
}
