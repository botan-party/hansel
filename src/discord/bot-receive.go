package discord

import (
	"hansel/aws"
	"hansel/config"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TargetChannel Botがメッセージを投稿するDiscordチャンネル
type TargetChannel struct {
	s     *discordgo.Session
	event *discordgo.MessageCreate
}

func (tc *TargetChannel) messageSend(message string) error {
	// コマンドが投稿されたチャンネル
	targetChannel, err := tc.s.State.Channel(tc.event.ChannelID)
	if err != nil {
		log.Println("チャンネルの取得に失敗 :", err)
		return err
	}

	// Botからメッセージ投稿
	if _, err := tc.s.ChannelMessageSend(targetChannel.ID, message); err != nil {
		log.Println("チャンネルメッセージの送信に失敗 :", err)
		return err
	}
	return nil
}

func (b *Bot) receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	targetChannel := TargetChannel{
		s:     s,
		event: event,
	}

	messages, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if event.Content == messages.StartTriggerMessage {
		// 起動時
		log.Println("開始 : インスタンス起動...")
		targetChannel.messageSend("インスタンスの起動コマンドを検知")

		statusErr := aws.StartInstance()
		if !statusErr.IsEmpty() {
			errLogMsg := ""
			errDiscordMsg := ""
			switch statusErr.Code {
			case aws.ERR_FAILED_START_INSTANCE:
				errLogMsg = "起動に失敗した :"
				errDiscordMsg = "インスタンスの起動に失敗"
			case aws.ERR_INVALID_RESPONSE_START_INSTANCE:
				errLogMsg = "起動時のレスポンスに異常 :"
				errDiscordMsg = "インスタンスの起動に失敗"
			case aws.ERR_INSTANCE_ALREADY_STARTED:
				errLogMsg = "既に起動している"
				errDiscordMsg = "インスタンスは起動済み"
			case aws.ERR_STARTING_INSTANCE:
				errLogMsg = "起動処理実行中"
				errDiscordMsg = "インスタンスは既に起動準備中"
			case aws.ERR_FAILED_WAIT_START_INSTANCE:
				errLogMsg = "起動待ちに失敗した"
				errDiscordMsg = "インスタンスの起動状態不明　再度のコマンド入力を要求"
			}
			log.Println(errLogMsg, statusErr.Err)
			targetChannel.messageSend(errDiscordMsg)
			return
		}

		log.Println("正常終了 : インスタンス起動")
		targetChannel.messageSend("インスタンスの起動に成功")

		// IPアドレス通知
		log.Println("IPアドレス取得待機中...")
		time.Sleep(time.Second)

		ipaddress, statusErr := aws.GetIPAddress()
		if !statusErr.IsEmpty() {
			errLogMsg := ""
			errDiscordMsg := "IPアドレスの取得に失敗"
			switch statusErr.Code {
			case aws.ERR_FAILED_GET_IP_ADDRESS:
				errLogMsg = "IPアドレス取得時、コマンド実行に失敗 : "
			case aws.ERR_INVALID_RESPONSE_GET_IP_ADDRESS:
				errLogMsg = "IPアドレス取得時のレスポンスに異常 :"
			}
			log.Println(errLogMsg, err)
			targetChannel.messageSend(errDiscordMsg)
			return
		}

		targetChannel.messageSend("今回のIPアドレス : " + ipaddress)

	} else if event.Content == messages.HibernateTriggerMessage {
		// 停止時
		log.Println("開始 : インスタンス停止...")
		targetChannel.messageSend("インスタンスの停止コマンドを検知")

		statusErr := aws.StopInstance()
		if !statusErr.IsEmpty() {
			errLogMsg := ""
			errDiscordMsg := ""
			switch statusErr.Code {
			case aws.ERR_FAILED_STOP_INSTANCE:
				errLogMsg = "停止に失敗した :"
				errDiscordMsg = "インスタンスの停止に失敗"
			case aws.ERR_INVALID_RESPONSE_STOP_INSTANCE:
				errLogMsg = "停止時のレスポンスに異常 :"
				errDiscordMsg = "インスタンスの停止に失敗"
			case aws.ERR_INSTANCE_ALREADY_STOPPED:
				errLogMsg = "既に停止している"
				errDiscordMsg = "インスタンスは停止済み"
			case aws.ERR_STOPPING_INSTANCE:
				errLogMsg = "停止処理実行中"
				errDiscordMsg = "インスタンスは既に停止準備中"
			case aws.ERR_FAILED_WAIT_STOP_INSTANCE:
				errLogMsg = "停止待ちに失敗した :"
				errDiscordMsg = "インスタンスの停止状態不明　再度のコマンド入力を要求"
			}
			log.Println(errLogMsg, err)
			targetChannel.messageSend(errDiscordMsg)
			return
		}

		log.Println("正常終了 : インスタンス停止")
		targetChannel.messageSend("インスタンスの停止に成功")

	} else if event.Content == messages.GetStatusTriggerMessage {
		// 起動状態の確認(IPアドレスの取得)
		log.Println("開始 : インスタンスステータス確認")
		targetChannel.messageSend("インスタンスの確認コマンドを検知")

		ipaddress, statusErr := aws.GetIPAddress()
		if !statusErr.IsEmpty() {
			errLogMsg := ""
			errDiscordMsg := "インスタンスの確認に失敗"
			switch statusErr.Code {
			case aws.ERR_FAILED_GET_IP_ADDRESS:
				errLogMsg = "IPアドレス取得時、コマンド実行に失敗 : "
			case aws.ERR_INVALID_RESPONSE_GET_IP_ADDRESS:
				errLogMsg = "IPアドレス取得時のレスポンスに異常 :"
			}
			log.Println(errLogMsg, err)
			targetChannel.messageSend(errDiscordMsg)
			return
		}

		if ipaddress != "" {
			targetChannel.messageSend("インスタンスは起動済み :" + ipaddress)
		} else {
			targetChannel.messageSend("インスタンスは未起動")
		}

	}
}
