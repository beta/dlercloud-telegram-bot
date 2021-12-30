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

	info, err := bot.dler.GetUserInfo(ctx)
	if err != nil {
		log.Errorf("failed to get user info from Dler Cloud, error: %+v", err)
		bot.telebot.Send(m.Sender, "Opps，查询失败")
		return
	}

	bot.telebot.Send(m.Chat, fmt.Sprintf("已用流量: %s\n可用流量: %s", info.Used, info.Unused))
}
