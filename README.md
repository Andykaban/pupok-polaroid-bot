pupok-polaroid-bot
===========

Small and funny telegram bot which convert photos like a polaroid camera

## How to use

1 - Install Docker

2 - Clone this repo

3 - make
    or GOOS=linux GOARCH=amd64 make

4 - Setup config.json bot configuration file by example
```
{
    "bot_token" : "Telegram bot token",
    "bot_proxy_url" : "SOCK5 proxy url (optional if needed)",
    "bot_proxy_login" : "SOCK5 proxy login (optional if use proxy auth)",
    "bot_proxy_password" : "SOCK5 proxy login (optional if use proxy auth)",
    "bot_temp_dir" : "/tmp/picture_folder (or any other existing directory)",
    "font_path" : "./static/fonts/wqy-zenhei.ttf (You can use OS system ttf fonts)",
    "background_path" : "./static/images/background.png (Path to polaroid camera background)"
}
```
5 - Run the bot ./pupok-polaroid-bot %path_to_config% (without parameter if config placed in the same directory)

6 - Send /start command to the bot
