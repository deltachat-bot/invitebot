package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deltachat-bot/deltabot-cli-go/botcli"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var cli = botcli.New("invitebot")

func main() {
	cli.OnBotInit(func(bot *deltachat.Bot, cmd *cobra.Command, args []string) {
		name, _ := bot.GetConfig("displayname")
		if name == "" {
			bot.SetConfig("displayname", "Invite Bot")
			bot.SetConfig("selfstatus", "I am a bot that helps you invite friends to your private groups, send me /help for more info")
		}
		bot.OnNewMsg(func(message *deltachat.Message) { onNewMsg(bot, message) })
	})
	cli.OnBotStart(func(bot *deltachat.Bot, cmd *cobra.Command, args []string) {
		addr, _ := bot.GetConfig("configured_addr")
		cli.Logger.Infof("Listening at: %v", addr)
	})
	cli.Start()
}

func onNewMsg(bot *deltachat.Bot, message *deltachat.Message) {
	msg, err := message.Snapshot()
	if err != nil || msg.IsInfo || msg.IsBot {
		return
	}
	chat := &deltachat.Chat{bot.Account, msg.ChatId}
	args := strings.Split(msg.Text, " ")
	switch args[0] {
	case "/invite":
		chatInfo, err := chat.BasicSnapshot()
		if err != nil {
			cli.Logger.Error(err)
			return
		}
		if chatInfo.ChatType == deltachat.ChatGroup {
			sendInviteQr(chat)
		} else {
			chat.SendText("The /invite command can only be used in groups, send /help for more info")
		}
	case "/help":
		sendHelp(chat)
	default:
		chatInfo, err := chat.BasicSnapshot()
		if err != nil {
			cli.Logger.Error(err)
			return
		}
		if chatInfo.ChatType != deltachat.ChatSingle {
			return
		}
		sendHelp(chat)
	}
}

func sendHelp(chat *deltachat.Chat) {
	text := "I am a bot that can help you invite friends to your private groups using a QR.\n\n"
	text += "You can also share your own invitation QR with them so why would you need me?\n"
	text += "Well, if you share your QR, your friends will be able to join only when you are online, but since I am a bot I am always online!\n\n"
	text += "To get the invitation QR of a group, add me to the group and send in the group:\n\n/invite\n\n"
	text += "I will share the invitation QR, you can then send it to friends you want to invite.\n\n"
	text += "If you want to revoque te invitation QR just remove me from the group"
	chat.SendText(text)
}

func sendInviteQr(chat *deltachat.Chat) {
	qrdata, _, err := chat.QrCode()
	if err != nil {
		cli.Logger.Error(err)
		return
	}
	chatInfo, err := chat.BasicSnapshot()
	if err != nil {
		cli.Logger.Error(err)
		return
	}

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		cli.Logger.Error(err)
		return
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "qr.png")

	err = qrcode.WriteFile(qrdata, qrcode.Medium, 256, path)
	if err != nil {
		cli.Logger.Error(err)
		return
	}
	text := fmt.Sprintf("Scan to join group %s", chatInfo.Name)
	chat.SendMsg(deltachat.MsgData{Text: text, File: path})
}
