package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/deltachat-bot/deltabot-cli-go/botcli"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var cli = botcli.New("invitebot")

func onBotInit(cli *botcli.BotCli, bot *deltachat.Bot, cmd *cobra.Command, args []string) {
	bot.OnNewMsg(onNewMsg)

	accounts, err := bot.Rpc.GetAllAccountIds()
	if err != nil {
		cli.Logger.Error(err)
	}
	for _, accId := range accounts {
		name, err := bot.Rpc.GetConfig(accId, "displayname")
		if err != nil {
			cli.Logger.Error(err)
		}
		if name.UnwrapOr("") == "" {
			err = bot.Rpc.SetConfig(accId, "displayname", option.Some("InviteBot"))
			if err != nil {
				cli.Logger.Error(err)
			}
			status := "I am a bot that helps you invite friends to your private groups, send me /help for more info"
			err = bot.Rpc.SetConfig(accId, "selfstatus", option.Some(status))
			if err != nil {
				cli.Logger.Error(err)
			}
			err = bot.Rpc.SetConfig(accId, "delete_server_after", option.Some("1"))
			if err != nil {
				cli.Logger.Error(err)
			}
			err = bot.Rpc.SetConfig(accId, "delete_device_after", option.Some("1800"))
			if err != nil {
				cli.Logger.Error(err)
			}
		}
	}
}

func onNewMsg(bot *deltachat.Bot, accId deltachat.AccountId, msgId deltachat.MsgId) {
	logger := cli.GetLogger(accId).With("msg", msgId)
	msg, err := bot.Rpc.GetMessage(accId, msgId)
	if err != nil {
		logger.Error(err)
		return
	}

	if !msg.IsBot && !msg.IsInfo && msg.FromId > deltachat.ContactLastSpecial {
		chat, err := bot.Rpc.GetBasicChatInfo(accId, msg.ChatId)
		if err != nil {
			logger.Error(err)
			return
		}
		if chat.ChatType == deltachat.ChatSingle || strings.HasPrefix(msg.Text, "/") {
			err = bot.Rpc.MarkseenMsgs(accId, []deltachat.MsgId{msg.Id})
			if err != nil {
				logger.Error(err)
			}
		}

		args := strings.Split(msg.Text, " ")
		switch args[0] {
		case "/invite":
			if chat.ChatType == deltachat.ChatGroup {
				sendInviteQr(bot.Rpc, accId, msg.ChatId)
			} else {
				text := "The /invite command can only be used in groups, send /help for more info"
				_, err := bot.Rpc.SendMsg(accId, msg.ChatId, deltachat.MsgData{Text: text})
				if err != nil {
					logger.Error(err)
				}
			}
		case "/help":
			sendHelp(bot.Rpc, accId, msg.ChatId)
		default:
			if chat.ChatType == deltachat.ChatSingle {
				sendHelp(bot.Rpc, accId, msg.ChatId)
			}
		}
	}

	if msg.FromId > deltachat.ContactLastSpecial {
		err = bot.Rpc.DeleteMessages(accId, []deltachat.MsgId{msg.Id})
		if err != nil {
			logger.Error(err)
		}
	}
}

func sendHelp(rpc *deltachat.Rpc, accId deltachat.AccountId, chatId deltachat.ChatId) {
	text := "I am a bot that can help you invite friends to your private groups.\n\n"
	text += "You can also share your own invitation QR with them so why would you need me?\n"
	text += "If you share your QR, your friends will be able to join only when you are online, but since I am a bot I am always online!\n\n"
	text += "To get the invitation QR of a group, add me to the group and send in the group:\n\n/invite\n\n"
	text += "I will share the invitation QR, you can then send it to friends you want to invite.\n\n"
	text += "If you want to revoque te invitation QR just remove me from the group"
	_, err := rpc.SendMsg(accId, chatId, deltachat.MsgData{Text: text})
	if err != nil {
		cli.GetLogger(accId).With("chat", chatId).Error(err)
	}
}

func sendInviteQr(rpc *deltachat.Rpc, accId deltachat.AccountId, chatId deltachat.ChatId) {
	logger := cli.GetLogger(accId).With("chat", chatId)
	qrdata, _, err := rpc.GetChatSecurejoinQrCodeSvg(accId, option.Some(chatId))
	if err != nil {
		logger.Error(err)
		return
	}

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		logger.Error(err)
		return
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "qr.png")

	err = qrcode.WriteFile(qrdata, qrcode.Medium, 256, path)
	if err != nil {
		logger.Error(err)
		return
	}
	_, err = rpc.SendMsg(accId, chatId, deltachat.MsgData{Text: qrdata, File: path})
	if err != nil {
		logger.Error(err)
	}
}

func main() {
	cli.OnBotInit(onBotInit)
	if err := cli.Start(); err != nil {
		cli.Logger.Error(err)
	}
}
