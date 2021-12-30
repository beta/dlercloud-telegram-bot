# Dler Cloud Telegram Bot

A Telegram bot for managing your Dler Cloud account.

## Usage

1. Download [the latest release](https://github.com/beta/dlercloud-telegram-bot/releases/latest)
2. Fill the values into [config.toml](config.toml) and save it somewhere
3. `./bot -c <path-to-config.toml>`

# Configs

```toml
[telegram]
bot-token = "" # Bot token acquired from @BotFather

# Only user/group/channel with this recipient ID can use this bot.
# Omit this value so that any one can use it.
allowed-recipient = ""

[dler-cloud]
# Your Dler Cloud account.
email = ""
password = ""
```

## License

zlib
