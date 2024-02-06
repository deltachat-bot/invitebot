package main

import (
	"testing"

	"github.com/deltachat-bot/deltabot-cli-go/botcli"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/stretchr/testify/require"
)

type TestCallback func(bot *deltachat.Bot, botAcc deltachat.AccountId, userRpc *deltachat.Rpc, userAcc deltachat.AccountId)

var acfactory *deltachat.AcFactory

func TestMain(m *testing.M) {
	acfactory = &deltachat.AcFactory{}
	acfactory.TearUp()
	defer acfactory.TearDown()
	m.Run()
}

func withBotAndUser(callback TestCallback) {
	acfactory.WithOnlineBot(func(bot *deltachat.Bot, botAcc deltachat.AccountId) {
		acfactory.WithOnlineAccount(func(userRpc *deltachat.Rpc, userAcc deltachat.AccountId) {
			cli := &botcli.BotCli{AppDir: acfactory.MkdirTemp()}
			onBotInit(cli, bot, nil, nil)
			go bot.Run() //nolint:errcheck
			callback(bot, botAcc, userRpc, userAcc)
		})
	})
}

func TestBot(t *testing.T) {
	withBotAndUser(func(bot *deltachat.Bot, botAcc deltachat.AccountId, userRpc *deltachat.Rpc, userAcc deltachat.AccountId) {
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
		require.NotEmpty(t, msg.File)
		require.Contains(t, msg.Text, "https://i.delta.chat")
	})
}
