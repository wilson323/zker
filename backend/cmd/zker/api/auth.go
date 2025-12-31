package api

import (
	"fmt"
)

// ValidateAPIKey 验证API Key
func (c *Client) ValidateAPIKey(apiKey string) error {
	client := NewClient(c.baseURL, apiKey)
	resp, err := client.get("/api/developer/auth/validate")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("API Key验证失败 (status: %d): %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// ValidateToken 验证JWT Token
func (c *Client) ValidateToken(token string) error {
	client := NewClient(c.baseURL, token)
	resp, err := client.get("/api/developer/user/me")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("Token验证失败 (status: %d): %s", resp.StatusCode(), resp.String())
	}

	return nil
}
