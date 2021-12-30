// Copyright (c) 2021 Beta Kuang <beta.kuang@gmail.com>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would be
//    appreciated but is not required.
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
// 3. This notice may not be removed or altered from any source distribution.

package dler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// NewClient 返回 Dler Cloud API 客户端.
func NewClient(email string, password string) *Client {
	return &Client{email: email, password: password}
}

// Client Dler Cloud API 客户端.
type Client struct {
	email    string
	password string

	token string
}

// HasLoggedIn 返回是否已登录.
func (c *Client) HasLoggedIn() bool {
	return len(c.token) > 0
}

// Login 登录.
func (c *Client) Login(ctx context.Context) error {
	var response = new(struct {
		Token string `json:"token"`
	})
	err := c.post(ctx, "login", map[string]interface{}{
		"email":  c.email,
		"passwd": c.password,
	}, response)
	if err != nil {
		return err
	}

	c.token = response.Token
	return nil
}

// GetUserInfo 获取用户信息.
func (c *Client) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	if !c.HasLoggedIn() {
		return nil, fmt.Errorf("not logged in")
	}

	var response = new(UserInfo)
	err := c.post(ctx, "information", map[string]interface{}{
		"access_token": c.token,
	}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UserInfo 用户信息.
type UserInfo struct {
	Plan      string `json:"plan"`
	PlanTime  string `json:"plan_time"`
	Money     string `json:"money"`
	AffMoney  string `json:"aff_money"`
	TodayUsed string `json:"today_used"`
	Used      string `json:"used"`
	Unused    string `json:"unused"`
	Traffic   string `json:"traffic"`
	Integral  string `json:"integral"`
}

func (c *Client) getURL(path string) string {
	const urlFmt = `https://dler.cloud/api/v1/%s`
	return fmt.Sprintf(urlFmt, path)
}

func (c *Client) post(ctx context.Context, path string, body map[string]interface{}, dest interface{}) error {
	var (
		httpReq *http.Request
		err     error
	)

	if body != nil {
		form := make(url.Values, len(body))
		for k, v := range body {
			form.Set(k, fmt.Sprint(v))
		}
		httpReq, err = http.NewRequest(http.MethodPost, c.getURL(path), strings.NewReader(form.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %+v", err)
		}
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		httpReq, err = http.NewRequest(http.MethodPost, c.getURL(path), nil)
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %+v", err)
		}
	}

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to do request: %+v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status code: %d", httpResp.StatusCode)
	}
	if httpResp.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %+v", err)
	}
	resp := new(response)
	if err := json.Unmarshal(respBody, resp); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %+v", err)
	}
	if resp.Code != http.StatusOK {
		return fmt.Errorf("invalid result code %d with message: %s", resp.Code, resp.Message)
	}

	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(resp.Data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal response data: %+v", err)
	}
	return nil
}

type response struct {
	Code    int             `json:"ret"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}
