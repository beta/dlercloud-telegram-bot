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

package bot

import (
	"context"
	"fmt"
	"time"

	"dlercloud-telegarm-bot/internal/log"

	"gopkg.in/tucnak/telebot.v2"
)

// Info 查询信息.
func (bot *Bot) Info(m *telebot.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dlerInfo, err := bot.dler.GetUserInfo(ctx)
	if err != nil {
		log.Errorf("failed to get user info from Dler Cloud, error: %+v", err)
		bot.telebot.Send(m.Sender, "Opps，查询失败")
		return
	}

	if bot.vultrEnabled {
		vultrInfo, err := bot.queryVultrInfo(ctx)
		if err != nil {
			log.Errorf("failed to get bandwidth info from Vultr, error: %+v", err)
			bot.telebot.Send(m.Sender, "Opps，查询失败")
			return
		}

		msg := fmt.Sprintf(`*Dler Cloud*
已用流量: %s
可用流量: %s

`, dlerInfo.Used, dlerInfo.Unused)
		for _, inst := range vultrInfo {
			msg += fmt.Sprintf(`*%s*
已用流量: %s
可用流量: %s

`, inst.Name, inst.Used, inst.Unused)
		}

		bot.telebot.Send(m.Chat, msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
		return
	}

	bot.telebot.Send(m.Chat, fmt.Sprintf("已用流量: %s\n可用流量: %s", dlerInfo.Used, dlerInfo.Unused))
}

func (bot *Bot) queryVultrInfo(ctx context.Context) ([]*vultrInstanceInfo, error) {
	const (
		dateFmt = "2006-01-02"
	)

	tz, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %+v", err)
	}
	currentMonth := time.Now().In(tz).Month()

	// 查询所有实例的流量总额
	instances, err := bot.vultr.GetInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query all instances: %+v", err)
	}
	totalGiBs := make(map[string]float64, len(instances))
	for _, inst := range instances {
		totalGiBs[inst.ID] = float64(inst.AllowedBandwidthGiB)
	}

	ret := make([]*vultrInstanceInfo, 0, len(bot.vultrInstances))
	for _, inst := range bot.vultrInstances {
		if _, exist := totalGiBs[inst.InstanceID]; !exist {
			return nil, fmt.Errorf("vultr instance %s not found in your account", inst.InstanceID)
		}

		bandwidth, err := bot.vultr.GetInstanceBandwidth(ctx, inst.InstanceID)
		if err != nil {
			return nil, err
		}

		var usedBytes int64
		for date, usage := range bandwidth {
			d, err := time.Parse(dateFmt, date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse date %s, error: %+v", date, err)
			}
			if d.Month() != currentMonth {
				continue
			}

			usedBytes += usage.IncomingBytes + usage.OutgoingBytes
		}

		usedGiBs := float64(usedBytes) / (1024 * 1024 * 1024)
		unusedGiB := totalGiBs[inst.InstanceID] - usedGiBs
		ret = append(ret, &vultrInstanceInfo{
			Name:   inst.Name,
			Used:   fmt.Sprintf("%.2fGiB", usedGiBs),
			Unused: fmt.Sprintf("%.2fGiB", unusedGiB),
		})
	}

	return ret, nil
}

type vultrInstanceInfo struct {
	Name   string
	Used   string
	Unused string
}
