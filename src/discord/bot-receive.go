package discord

import (
	"hansel/aws"
	"hansel/config"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	messages, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if event.Content == messages.StartTriggerMessage {
		// 起動時
		log.Println("開始 : インスタンス起動...")
		_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスの起動コマンドを検知")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

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
			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		log.Println("正常終了 : インスタンス起動")
		_, err = s.ChannelMessageSend(event.ChannelID, "インスタンスの起動に成功")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

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
			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		_, err = s.ChannelMessageSend(event.ChannelID, "今回のIPアドレス : "+ipaddress)
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

	} else if event.Content == messages.HibernateTriggerMessage {
		// 停止時
		log.Println("開始 : インスタンス停止...")
		_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスの停止コマンドを検知")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

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
			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		log.Println("正常終了 : インスタンス停止")
		_, err = s.ChannelMessageSend(event.ChannelID, "インスタンスの停止に成功")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

	} else if event.Content == messages.GetStatusTriggerMessage {
		// 起動状態の確認(IPアドレスの取得)
		log.Println("開始 : インスタンスステータス確認")
		_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスの確認コマンドを検知")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

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
			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		if ipaddress != "" {
			_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスは起動済み :"+ipaddress)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
		} else {
			_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスは未起動")
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
		}
	}
}
