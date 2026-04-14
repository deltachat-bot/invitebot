package main

import (
	"testing"

	"github.com/chatmail/rpc-client-go/v2/deltachat"
	"github.com/deltachat-bot/deltabot-cli-go/v2/botcli"
	"github.com/stretchr/testify/require"
)

type TestCallback func(bot *deltachat.Bot, botAcc uint32, userRpc *deltachat.Rpc, userAcc uint32)

var acfactory *deltachat.AcFactory

func TestMain(m *testing.M) {
	acfactory = &deltachat.AcFactory{}
	acfactory.TearUp()
	defer acfactory.TearDown()
	m.Run()
}

func withBotAndUser(callback TestCallback) {
	acfactory.WithOnlineBot(func(bot *deltachat.Bot, botAcc uint32) {
		acfactory.WithOnlineAccount(func(userRpc *deltachat.Rpc, userAcc uint32) {
			cli := &botcli.BotCli{AppDir: acfactory.MkdirTemp()}
			onBotInit(cli, bot, nil, nil)
			go bot.Run() //nolint:errcheck
			callback(bot, botAcc, userRpc, userAcc)
		})
	})
}

func TestBot(t *testing.T) {
	withBotAndUser(func(bot *deltachat.Bot, botAcc uint32, userRpc *deltachat.Rpc, userAcc uint32) {
		chatWithBot := acfactory.CreateChat(userRpc, userAcc, bot.Rpc, botAcc)
		_, err := userRpc.MiscSendTextMessage(userAcc, chatWithBot, "hi")
		require.Nil(t, err)
		msg := acfactory.NextMsg(userRpc, userAcc)
		require.Contains(t, msg.Text, "I am a bot")

		groupWithBot, err := userRpc.CreateGroupChat(userAcc, "test group", false)
		require.Nil(t, err)
		require.Nil(t, userRpc.AddContactToChat(userAcc, groupWithBot, msg.FromId))
		_, err = userRpc.MiscSendTextMessage(userAcc, groupWithBot, "hi")
		require.Nil(t, err)
		_, err = userRpc.MiscSendTextMessage(userAcc, groupWithBot, "/invite")
		require.Nil(t, err)
		msg = acfactory.NextMsg(userRpc, userAcc)
		require.Contains(t, msg.Text, "https://i.delta.chat")
	})
}
