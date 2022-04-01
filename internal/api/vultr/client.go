// Copyright (c) 2022 Beta Kuang <beta.kuang@gmail.com>
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

package vultr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewClient 返回 Vultr API 客户端.
func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type Client struct {
	apiKey string
}

// GetInstances 查询所有实例.
func (c *Client) GetInstances(ctx context.Context) ([]*Instance, error) {
	var response struct {
		Instances []*Instance `json:"instances"`
	}

	if err := c.get(ctx, "instances", &response); err != nil {
		return nil, err
	}

	return response.Instances, nil
}

type Instance struct {
	ID                  string `json:"id"`
	AllowedBandwidthGiB int    `json:"allowed_bandwidth"`
}

// GetInstanceBandwidth 查询实例过去一个月的带宽使用情况.
func (c *Client) GetInstanceBandwidth(ctx context.Context, instanceID string) (map[string]BandwidthUsage, error) {
	var response struct {
		Bandwidth map[string]BandwidthUsage `json:"bandwidth"`
	}

	if err := c.get(ctx, fmt.Sprintf("instances/%s/bandwidth", instanceID), &response); err != nil {
		return nil, err
	}

	return response.Bandwidth, nil
}

type BandwidthUsage struct {
	IncomingBytes int64 `json:"incoming_bytes"`
	OutgoingBytes int64 `json:"outgoing_bytes"`
}

func (c *Client) getURL(path string) string {
	const urlFmt = `https://api.vultr.com/v2/%s`
	return fmt.Sprintf(urlFmt, path)
}

func (c *Client) get(ctx context.Context, path string, dest interface{}) error {
	httpReq, err := http.NewRequest(http.MethodGet, c.getURL(path), nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %+v", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
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

	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, dest); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %+v", err)
	}
	return nil
}
