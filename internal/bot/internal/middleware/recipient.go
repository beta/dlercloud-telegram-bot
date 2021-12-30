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
	"dlercloud-telegarm-bot/internal/log"

	"gopkg.in/tucnak/telebot.v2"
)

// FilterRecipient 限制消息来源中间件.
func FilterRecipient(allowedRecipient string) func(*telebot.Update) bool {
	return func(update *telebot.Update) bool {
		if update == nil {
			return false
		}
		if update.Message == nil {
			log.Errorf("[Update updateID=%d] non-message update, ignore", update.ID)
			return false
		}

		m := update.Message
		if m.Chat == nil {
			log.Errorf("[Message updateID=%d] chat is nil", update.ID)
			return false
		}
		if m.Chat.Recipient() != allowedRecipient {
			log.Errorf("[Message updateID=%d] message from unallowed recipient %s, chatTitle=%s, sender=%s, ignore", update.ID, m.Chat.Recipient(), m.Chat.Title, getSenderName(m.Sender))
			return false
		}

		return true
	}
}
