package api

import (
	"encoding/json"
	"fmt"
)

// Bot Bot信息
type Bot struct {
	ID          string `json:"bot_id"`
	Name        string `json:"bot_name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Status      string `json:"status"`
	CreatedAt   string `json:"create_time"`
	UpdatedAt   string `json:"update_time"`
}

// CreateBotRequest 创建Bot请求
type CreateBotRequest struct {
	Name        string `json:"bot_name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Template    string `json:"template"`
}

// DeployBotRequest 部署Bot请求
type DeployBotRequest struct {
	BotID  string `json:"bot_id"`
	Env    string `json:"env"`
	Config string `json:"config"`
}

// Deployment 部署信息
type Deployment struct {
	ID     string `json:"deployment_id"`
	BotID  string `json:"bot_id"`
	Env    string `json:"env"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

// ListBots 列出所有Bot
func (c *Client) ListBots() ([]Bot, error) {
	resp, err := c.get("/api/bot/list")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list bots (status: %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Items []Bot `json:"items"`
			Total int   `json:"total"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return result.Data.Items, nil
}

// CreateBot 创建Bot
func (c *Client) CreateBot(req *CreateBotRequest) (*Bot, error) {
	resp, err := c.post("/api/bot/create", req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to create bot (status: %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    Bot   `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return &result.Data, nil
}

// GetBot 获取Bot详情
func (c *Client) GetBot(botID string) (*Bot, error) {
	resp, err := c.get(fmt.Sprintf("/api/bot/%s", botID))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get bot (status: %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    Bot   `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return &result.Data, nil
}

// DeployBot 部署Bot
func (c *Client) DeployBot(botID, env string) (*Deployment, error) {
	req := &DeployBotRequest{
		BotID: botID,
		Env:   env,
	}

	resp, err := c.post("/api/bot/deploy", req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to deploy bot (status: %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Code       int        `json:"code"`
		Message    string     `json:"message"`
		Data       Deployment `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return &result.Data, nil
}

// DeleteBot 删除Bot
func (c *Client) DeleteBot(botID string) error {
	resp, err := c.delete(fmt.Sprintf("/api/bot/%s", botID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete bot (status: %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %s", result.Message)
	}

	return nil
}
