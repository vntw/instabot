# InstaBot

Post images of an Instagram user periodically to Slack.

## Usage
Install dependencies via `dep` and run `make`.

Copy `config.json.dist` to `config.json` and configure your options. Then simply run the binary.

## Config
| Option | Description |
| --- | --- |
| `user` | The Instagram username |
| `wait-time` | The wait time between crawls in minutes, default is `10` |
| `slack-token` | Your bot's Slack token |
| `slack-channel` | The channel name where to post the images, e.g. `channel` |
| `proxy` | _(optional)_ The proxy URL if you have one, e.g. `http://x.x.x.x:x` |
