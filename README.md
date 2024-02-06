#  InviteBot

![Latest release](https://img.shields.io/github/v/tag/deltachat-bot/invitebot?label=release)
[![CI](https://github.com/deltachat-bot/invitebot/actions/workflows/ci.yml/badge.svg)](https://github.com/deltachat-bot/invitebot/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-65.4%25-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/deltachat-bot/invitebot)](https://goreportcard.com/report/github.com/deltachat-bot/invitebot)

Small bot that allows to generate invitation QRs for your private Delta Chat groups. The bot is always online
and can add people to groups in "real time" while if you use your own invitation QRs, others will not be able
to join until you are online.

## Install

Binary releases can be found at: https://github.com/deltachat-bot/invitebot/releases

To install from source:

```sh
go install github.com/deltachat-bot/invitebot@latest
```

### Installing deltachat-rpc-server

This program depends on a standalone Delta Chat RPC server `deltachat-rpc-server` program that must be
available in your `PATH`. For installation instructions check:
https://github.com/deltachat/deltachat-core-rust/tree/master/deltachat-rpc-server

## Running the bot

Configure the bot:

```sh
invitebot init bot@example.com PASSWORD
```

Start the bot:

```sh
invitebot serve
```

Run `invitebot --help` to see all available options.


## Usage in Delta Chat

Once the bot is running:

1. Add the bot address to some group in Delta Chat.
2. Send `/invite` in the group.
3. The bot will reply with an invitation QR.
4. Share the invitation QR with friends.
5. To revoque invitations, simply remove the bot from the group.
