package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// Client API客户端
type Client struct {
	httpClient *resty.Client
	baseURL    string
	authToken  string
}

// NewClient 创建新的API客户端
func NewClient(baseURL, authToken string) *Client {
	client := resty.New().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "ZKER-CLI/1.0.0")

	if authToken != "" {
		client.SetAuthToken(authToken)
	}

	return &Client{
		httpClient: client,
		baseURL:    baseURL,
		authToken:  authToken,
	}
}

// SetDebug 设置调试模式
func (c *Client) SetDebug(debug bool) *Client {
	c.httpClient.SetDebug(debug)
	return c
}

// get 发送GET请求
func (c *Client) get(path string) (*resty.Response, error) {
	resp, err := c.httpClient.R().Get(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

// post 发送POST请求
func (c *Client) post(path string, body interface{}) (*resty.Response, error) {
	resp, err := c.httpClient.R().SetBody(body).Post(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

// put 发送PUT请求
func (c *Client) put(path string, body interface{}) (*resty.Response, error) {
	resp, err := c.httpClient.R().SetBody(body).Put(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

// delete 发送DELETE请求
func (c *Client) delete(path string) (*resty.Response, error) {
	resp, err := c.httpClient.R().Delete(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}
