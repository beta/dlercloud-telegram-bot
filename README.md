# Dler Cloud Telegram Bot

A Telegram bot for querying the bandwidth usage of your Dler Cloud & Vultr account.

## Usage

1. Download [the latest release](https://github.com/beta/dlercloud-telegram-bot/releases/latest)
2. Fill the values into [config.toml](config.toml) and save it somewhere
3. `./bot -c <path-to-config.toml>`

## Configs

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

[vultr]
# Change to true to use Vultr
enabled = false
# Your Vultr API key. Enable in https://my.vultr.com/settings/#settingsapi
api-key = ""

# Uncomment the following options to add Vultr instances.
# Change INSTANCE_NAME_* to an recognizable instance name.
# INSTANCE_ID_* can be found in the URL of Vultr's product page.

#   [vultr.instances.INSTANCE_NAME_1]
#   id = "INSTANCE_ID_1"

#   [vultr.instances.INSTANCE_NAME_2]
#   id = "INSTANCE_ID_2"

```

## License

zlib
