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

	"dlercloud-telegarm-bot/internal/api/dler"
	"dlercloud-telegarm-bot/internal/api/vultr"
	"dlercloud-telegarm-bot/internal/bot/internal/middleware"
	"dlercloud-telegarm-bot/internal/config"

	"gopkg.in/tucnak/telebot.v2"
)

// NewBot 返回新的 bot 实例.
func NewBot(cfg *config.Config) *Bot {
	bot := &Bot{
		dler:             dler.NewClient(cfg.DlerCloud.Email, cfg.DlerCloud.Password),
		vultrEnabled:     cfg.Vultr.Enabled,
		telebotSettings:  telebot.Settings{Token: cfg.Telegram.BotToken},
		allowedRecipient: cfg.Telegram.AllowedRecipient,
	}

	if cfg.Vultr.Enabled {
		bot.vultrInstances = make([]*vultrInstance, 0, len(cfg.Vultr.Instances))
		for name, inst := range cfg.Vultr.Instances {
			bot.vultrInstances = append(bot.vultrInstances, &vultrInstance{
				Name:       name,
				InstanceID: inst.ID,
			})
		}
		bot.vultr = vultr.NewClient(cfg.Vultr.APIKey)
	}

	return bot
}

// Bot.
type Bot struct {
	dler *dler.Client

	vultrEnabled   bool
	vultrInstances []*vultrInstance
	vultr          *vultr.Client

	telebotSettings  telebot.Settings
	allowedRecipient string

	telebot *telebot.Bot
}

type vultrInstance struct {
	Name       string
	InstanceID string
}

// Start 启动 bot.
func (bot *Bot) Start() error {
	if err := bot.loginToDler(); err != nil {
		return fmt.Errorf("failed to log in to Dler Cloud, error: %+v", err)
	}

	if err := bot.createTelebot(); err != nil {
		return fmt.Errorf("failed to create Telegram bot, error: %+v", err)
	}

	bot.registerRoutes()
	bot.telebot.Start()
	return nil
}

func (bot *Bot) loginToDler() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return bot.dler.Login(ctx)
}

func (bot *Bot) createTelebot() error {
	b, err := telebot.NewBot(bot.telebotSettings)
	if err != nil {
		return err
	}

	// 添加 logger 中间件
	b.Poller = telebot.NewMiddlewarePoller(&telebot.LongPoller{Timeout: 10 * time.Second}, bot.middleware())

	bot.telebot = b
	return nil
}

func (bot *Bot) middleware() func(u *telebot.Update) bool {
	return func(u *telebot.Update) bool {
		if u == nil {
			return false
		}
		if !middleware.Logger(u) {
			return false
		}
		if !middleware.FilterRecipient(bot.allowedRecipient)(u) {
			return false
		}
		return true
	}

}

func (bot *Bot) registerRoutes() {
	bot.telebot.Handle("/info", bot.Info)
}
