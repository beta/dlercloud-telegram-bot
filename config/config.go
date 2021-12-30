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

package config

import (
	"github.com/BurntSushi/toml"
)

// Config stores app configurations.
type Config struct {
	Telegram struct {
		BotToken         string `toml:"bot-token"`
		AllowedRecipient string `toml:"allowed-recipient"`
	} `toml:"telegram"`

	DlerCloud struct {
		Email    string `toml:"email"`
		Password string `toml:"password"`
	} `toml:"dler-cloud"`
}

// FromFile parse configs from file.
func FromFile(path string) (*Config, error) {
	cfg := new(Config)
	_, err := toml.DecodeFile(path, cfg)
	return cfg, err
}
