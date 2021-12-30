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

package middleware

import (
	"strings"

	"dlercloud-telegarm-bot/internal/log"

	"gopkg.in/tucnak/telebot.v2"
)

// Logger 输出消息日志中间件.
func Logger(update *telebot.Update) bool {
	if update == nil {
		return false
	}

	switch {
	case update.Message != nil:
		m := update.Message
		if m.Sender == nil {
			log.Errorf("[Message updateID=%d] sender is nil", update.ID)
			return false
		}
		log.Infof("[Message updateID=%d] sender=%s, fromGroup=%v, recipient=%s, content=%s", update.ID, getSenderName(m.Sender), m.FromGroup(), m.Chat.Recipient(), m.Text)

	default:
		log.Infof("[Update updateID=%d] non-message update", update.ID)
	}

	return true
}

func getSenderName(sender *telebot.User) string {
	if sender == nil {
		return ""
	}

	parts := make([]string, 0, 3)
	if sender.FirstName != "" {
		parts = append(parts, sender.FirstName)
	}
	if sender.LastName != "" {
		parts = append(parts, sender.LastName)
	}
	if sender.Username != "" {
		parts = append(parts, "(@"+sender.Username+")")
	}
	return strings.Join(parts, " ")
}
